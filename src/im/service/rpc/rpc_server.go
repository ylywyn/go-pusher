package rpc

import (
	log "im/common/log4go"
	"im/service/rpc/gen"
	"im/service/rpc/service/push"
	"net"

	grpc "google.golang.org/grpc"
)

//运行服务
func RunRpcService(s string) error {
	addr, err := net.ResolveTCPAddr("tcp4", s)
	if err != nil {
		log.Error("[RPCService|ResolveTCPAddr] %s error: (%v)", addr, err)
		return err
	}

	l, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		log.Error("[RPCService|ListenTCP] error: (%v)", err)
		return err
	}

	log.Debug("[rpc service|run] %s", s)

	server := grpc.NewServer()

	// Register Service
	gen.RegisterIPusherServer(server, &push.PushService{})

	server.Serve(l)
	return nil
}
