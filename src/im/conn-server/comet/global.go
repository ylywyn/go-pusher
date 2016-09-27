/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   :
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/

package comet

import (
	"time"
)

const (
	TcpAcceptRouts     = 2
	MaxIndex           = 65536
	TCP_PACKHEADER_LEN = 4
	WRITE_WAIT         = 8 * time.Second
)

const (
	StatusClosed  = -1
	StatusClosing = 0
	StatusRunning = 1
)
