/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   :
 *  Date   :
 *  Author :
 *  Description:
 ******************************************************************/

package web

import (
	"im/conn-server/comet"
)

func NetWorkStat(m map[string]interface{}) {
	stat := comet.Server.Stat()
	m[DATA] = stat
}
