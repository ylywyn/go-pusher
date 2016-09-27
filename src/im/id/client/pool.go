/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   :
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/

package id

import (
	"container/list"
	"crypto/rand"
	"crypto/sha1"
	"errors"
	"io"
	"strconv"
	"sync"
	"time"
)

var nowFunc = time.Now // for testing

// ErrPoolExhausted is returned from a pool connection method (Do, Send,
// Receive, Flush, Err) when the maximum number of database connections in the
// pool has been reached.
var ErrPoolExhausted = errors.New("client: connection pool exhausted")

var (
	errPoolClosed = errors.New("client: connection pool closed")
	errConnClosed = errors.New("client: connection closed")
)

//.
//
//  func newPool(server, password string) *redis.Pool {
//      return &redis.Pool{
//          MaxIdle: 3,
//          IdleTimeout: 240 * time.Second,
//          Dial: func () (redis.Conn, error) {
//              c, err := redis.Dial("tcp", server)
//              if err != nil {
//                  return nil, err
//              }
//              if _, err := c.Do("AUTH", password); err != nil {
//                  c.Close()
//                  return nil, err
//              }
//              return c, err
//          },
//          TestOnBorrow: func(c redis.Conn, t time.Time) error {
//              _, err := c.Do("PING")
//              return err
//          },
//      }
//  }
//
//  var (
//      pool *redis.Pool
//      redisServer = flag.String("redisServer", ":6379", "")
//      redisPassword = flag.String("redisPassword", "", "")
//  )
//
//  func main() {
//      flag.Parse()
//      pool = newPool(*redisServer, *redisPassword)
//      ...
//  }
//
// A request handler gets a connection from the pool and closes the connection
// when the handler is done:
//
//  func serveHome(w http.ResponseWriter, r *http.Request) {
//      conn := pool.Get()
//      defer conn.Close()
//      ....
//  }
//
type Pool struct {
	Dial func() (Conn, error)

	TestOnBorrow func(c Conn, t time.Time) error

	MaxIdle int

	MaxActive int

	IdleTimeout time.Duration

	Wait bool

	// mu protects fields defined below.
	mu     sync.Mutex
	cond   *sync.Cond
	closed bool
	active int

	// Stack of idleConn with most recently used at the front.
	idle list.List
}

type idleConn struct {
	c Conn
	t time.Time
}

// NewPool creates a new pool. This function is deprecated. Applications should
// initialize the Pool fields directly as shown in example.
func NewPool(newFn func() (Conn, error), maxIdle int) *Pool {
	return &Pool{Dial: newFn, MaxIdle: maxIdle}
}

func (p *Pool) Get() Conn {
	c, err := p.get()
	if err != nil {
		return errorConnection{err}
	}
	if c == nil {
		return errorConnection{err: errors.New("conn err")}
	}
	return &pooledConnection{p: p, c: c}
}

// ActiveCount returns the number of active connections in the pool.
func (p *Pool) ActiveCount() int {
	p.mu.Lock()
	active := p.active
	p.mu.Unlock()
	return active
}

// Close releases the resources used by the pool.
func (p *Pool) Close() error {
	p.mu.Lock()
	idle := p.idle
	p.idle.Init()
	p.closed = true
	p.active -= idle.Len()
	if p.cond != nil {
		p.cond.Broadcast()
	}
	p.mu.Unlock()
	for e := idle.Front(); e != nil; e = e.Next() {
		e.Value.(idleConn).c.Close()
	}
	return nil
}

// release decrements the active count and signals waiters. The caller must
// hold p.mu during the call.
func (p *Pool) release() {
	p.active -= 1
	if p.cond != nil {
		p.cond.Signal()
	}
}

// get prunes stale connections and returns a connection from the idle list or
// creates a new connection.
func (p *Pool) get() (Conn, error) {
	p.mu.Lock()

	// Prune stale connections.
	if timeout := p.IdleTimeout; timeout > 0 {
		for i, n := 0, p.idle.Len(); i < n; i++ {
			e := p.idle.Back()
			if e == nil {
				break
			}
			ic := e.Value.(idleConn)
			if ic.t.Add(timeout).After(nowFunc()) {
				break
			}
			p.idle.Remove(e)
			p.release()
			p.mu.Unlock()
			ic.c.Close()
			p.mu.Lock()
		}
	}

	for {

		// Get idle connection.

		for i, n := 0, p.idle.Len(); i < n; i++ {
			e := p.idle.Front()
			if e == nil {
				break
			}
			ic := e.Value.(idleConn)
			p.idle.Remove(e)
			test := p.TestOnBorrow
			p.mu.Unlock()
			if test == nil || test(ic.c, ic.t) == nil {
				return ic.c, nil
			}
			ic.c.Close()
			p.mu.Lock()
			p.release()
		}

		// Check for pool closed before dialing a new connection.

		if p.closed {
			p.mu.Unlock()
			return nil, errors.New("client: get on closed pool")
		}

		// Dial new connection if under limit.

		if p.MaxActive == 0 || p.active < p.MaxActive {
			dial := p.Dial
			p.active += 1
			p.mu.Unlock()
			c, err := dial()
			if err != nil {
				p.mu.Lock()
				p.release()
				p.mu.Unlock()
				c = nil
			}
			return c, err
		}

		if !p.Wait {
			p.mu.Unlock()
			return nil, ErrPoolExhausted
		}

		if p.cond == nil {
			p.cond = sync.NewCond(&p.mu)
		}
		p.cond.Wait()
	}
}

func (p *Pool) put(c Conn, forceClose bool) error {
	err := c.Err()
	p.mu.Lock()
	if !p.closed && err == nil && !forceClose {
		p.idle.PushFront(idleConn{t: nowFunc(), c: c})
		if p.idle.Len() > p.MaxIdle {
			c = p.idle.Remove(p.idle.Back()).(idleConn).c
		} else {
			c = nil
		}
	}

	if c == nil {
		if p.cond != nil {
			p.cond.Signal()
		}
		p.mu.Unlock()
		return nil
	}

	p.release()
	p.mu.Unlock()
	return c.Close()
}

type pooledConnection struct {
	p *Pool
	c Conn
}

var (
	sentinel     []byte
	sentinelOnce sync.Once
)

func initSentinel() {
	p := make([]byte, 64)
	if _, err := rand.Read(p); err == nil {
		sentinel = p
	} else {
		h := sha1.New()
		io.WriteString(h, "Oops, rand failed. Use time instead.")
		io.WriteString(h, strconv.FormatInt(time.Now().UnixNano(), 10))
		sentinel = h.Sum(nil)
	}
}

func (pc *pooledConnection) Close() error {
	c := pc.c
	if _, ok := c.(errorConnection); ok {
		return nil
	}

	//c.Do("")
	pc.p.put(c, false)
	return nil
}

func (pc *pooledConnection) Err() error {
	return pc.c.Err()
}

func (pc *pooledConnection) Do(commandName string) (reply *Reply, err error) {
	return pc.c.Do(commandName)
}

func (pc *pooledConnection) Send(commandName string) error {
	return pc.c.Send(commandName)
}

func (pc *pooledConnection) Flush() error {
	return pc.c.Flush()
}

func (pc *pooledConnection) Receive() (reply *Reply, err error) {
	return pc.c.Receive()
}

type errorConnection struct{ err error }

func (ec errorConnection) Do(string) (*Reply, error) { return nil, ec.err }
func (ec errorConnection) Send(string) error         { return ec.err }
func (ec errorConnection) Err() error                { return ec.err }
func (ec errorConnection) Close() error              { return ec.err }
func (ec errorConnection) Flush() error              { return ec.err }
func (ec errorConnection) Receive() (*Reply, error)  { return nil, ec.err }
