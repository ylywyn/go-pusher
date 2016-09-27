/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   :
 *  Date   :
 *  Author :
 *  Description:
 ******************************************************************/

package user

import (
	//"time"
	"im/service/logic/models/user"
	"im/service/logic/models/utils"

	"gopkg.in/mgo.v2"
)

//申请加好友
func AskForAdd_Friend(uid, uidaim uint64, appid uint32, notice bool) error {
	return nil
}

//添加好友， 双向添加
func AddFriend(uid, uidaim uint64, appid uint32, notice bool) error {
	uex := user.UserExt{Uid: uid, Appid: appid}
	err := uex.AddFriends(uidaim) //已经更新了版本
	if err != nil && err != utils.ErrExist {
		if err == mgo.ErrNotFound {
			//朋友关系表里没有记录
			uex.FriendsVersion = 1
			uex.Friends = []uint64{uidaim}
			err = uex.Insert()
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	ufaim := user.UserExt{Uid: uidaim, Appid: appid}
	err = ufaim.AddFriends(uid)
	if err != nil && err != utils.ErrExist {
		if err == mgo.ErrNotFound {
			ufaim.FriendsVersion = 1
			ufaim.Friends = []uint64{uid}
			err = ufaim.Insert()
			if err != nil {
				return err
			}
		} else {
			uex.RemoveFriends(uidaim)
			return err
		}
	}

	if notice {
		//cmd := msg.Friend_ADD_NOTICE{uint32(uid), uint32(uidaim)}
		//Notice_Friend(uint32(uid), uint32(uidaim), uint16(appid), uint16(fbs.FriendMsgTypeMT_FRIEND_ADD), cmd.Serialize())
	}
	return nil
}

//删除好友， 双向删除
func DelFriend(uid, uidaim uint64, appid uint32, notice bool) error {
	uex := user.UserExt{Uid: uid, Appid: appid}
	err := uex.RemoveFriends(uidaim)
	if err != nil && err != mgo.ErrNotFound {
		return err
	}

	ufaim := user.UserExt{Uid: uidaim, Appid: appid}
	err = ufaim.RemoveFriends(uid)
	if err != nil {
		if err != mgo.ErrNotFound {
			//出现错误，恢复
			uex.AddFriends(uidaim)
			return err
		}
	}

	if notice && err != mgo.ErrNotFound {
		//cmd := msg.Friend_DEL_NOTICE{uint32(uid), uint32(uidaim)}
		//Notice_Friend(uint32(uid), uint32(uidaim), uint16(appid), uint16(fbs.FriendMsgTypeMT_FRIEND_DEL), cmd.Serialize())
	}
	return nil
}

//func Notice_Friend(uid, frienduid uint32, appid, msgtype uint16, cmdbytes []byte) (uint64, error) {

//	Options_ := &msg.Options{
//		Time_live:       100000,
//		Start_time:      100000,
//		Apns_production: false,
//		Command:         cmdbytes,
//	}
//	msgid, err := push.GetMsgId("Friend_Relation")
//	if err != nil {
//		return 0, err
//	}

//	msg_ := &msg.Message{Type: msgtype,
//		Appid:    appid,
//		From:     uid,
//		To:       []uint32{frienduid},
//		Text:     "",
//		Time:     uint64(time.Now().Unix()),
//		Msgid:    msgid,
//		Platform: 1, //向所有平台
//		Payload:  nil,
//		Options:  Options_,
//	}

//	data := msg_.Serialize()

//	err = push.PushFriendMsg(data)
//	if err != nil {
//		return 0, err
//	} else {
//		return msg_.Msgid, nil
//	}
//}
