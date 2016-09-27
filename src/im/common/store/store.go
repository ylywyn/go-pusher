/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   :
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/

package store

import (
	"container/list"
)

//文件存储，暂未实现
type FileStore struct {
}

func NewFileStore() *FileStore {
	return nil
}

func (this *FileStore) Write([]byte) {

}

func (this *FileStore) WriteList(l *list.List) {

}
