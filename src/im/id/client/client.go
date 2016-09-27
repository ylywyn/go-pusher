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

type Error string

func (err Error) Error() string { return string(err) }

type Conn interface {
	Close() error

	Err() error

	Do(commandName string) (reply *Reply, err error)

	Send(commandName string) error

	Flush() error

	Receive() (reply *Reply, err error)
}
