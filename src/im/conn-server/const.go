/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   :
 *  Date   :
 *  Author :
 *  Description:
 ******************************************************************/

package main

import (
	"time"
)

const ConnTimeOut = 5 * 60 * time.Second

const (
	SignalShutDown = -1
	SignalClose    = 0
	SignalMsg      = 1
)
const (
	StatusClosed  = -1
	StatusClosing = 0
	StatusRunning = 1
	StatusLogin   = 2
)
