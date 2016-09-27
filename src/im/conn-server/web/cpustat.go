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
	"runtime"
)

func CpuStat(m map[string]interface{}) {

	res := make(map[string]interface{})
	res["numcpu"] = runtime.NumCPU()
	res["numgoroutine"] = runtime.NumGoroutine()

	m[DATA] = res
}
