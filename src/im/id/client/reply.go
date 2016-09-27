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

type Reply struct {
	value uint64
	str   string
}

func (this *Reply) Uint64() uint64 {
	return this.value
}

func (this *Reply) String() string {
	return this.str
}
