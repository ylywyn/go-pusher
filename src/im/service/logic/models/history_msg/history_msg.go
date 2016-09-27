/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   :
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/

package history_msg

import (
	"os"
	"pie-push-service/db"
	. "pie-push-service/db/history_msg"

	"github.com/astaxie/beego"

	"pie-utils/g"
)

func init() {
	//获取配置信息
	s, err := beego.AppConfig.GetSection("history_msg_service")
	if err != nil {
		g.Log.Error(err.Error())
		return
	}

	//是否启动
	if v, ok := s["service"]; (!ok) || (v == "off") {
		g.Log.Debug("history_msg_service is off")
		return
	}

	//获取 mogodb配置项
	var ok bool
	var dbAddr, dbName string
	if dbAddr, ok = s["mongodb_addr"]; !ok {
		g.Log.Error("mongodb_addr can't find")
		return
	}

	if dbName, ok = s["mongodb_name"]; !ok {
		g.Log.Error("mongodb_addr can't find")
		return
	}

	//初始化连接池
	err = db.InitMsgDB(dbAddr, dbName)
	if err != nil {
		g.Log.Error(err.Error())
		os.Exit(0)
		return
	}
}

func GetPersonalHistoryMsg(uid, friend_uid uint32, pageSize, pageIndex int) []HistoryMsg {
	hm := &HistoryMsgs{Uid: uid}

	hm.GetHistroyMsgs(friend_uid, pageSize, pageIndex)

	return hm.Msgs
}
func GetBroadcastHistoryMsg(appid uint16, pageSize, pageIndex int) []Broadcast_HistoryMsg {
	hm := &Broadcast_HistoryMsgs{Appid: appid}

	hm.GetHistroyMsgs(pageSize, pageIndex)

	return hm.Msgs
}
func GetPersonalOffLineMsg(uid uint32, lastmagid uint64) []HistoryMsg {
	hm := &HistoryMsgs{Uid: uid}

	hm.GetOffLineMsgs(lastmagid)

	return hm.Msgs
}
func GetBroadcastOffLineMsg(appid uint16, lastmagid uint64) []Broadcast_HistoryMsg {
	hm := &Broadcast_HistoryMsgs{Appid: appid}

	hm.GetOffLineMsgs(lastmagid)

	return hm.Msgs
}
