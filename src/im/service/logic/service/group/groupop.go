/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   :
 *  Date   :
 *  Author :
 *  Description:
 ******************************************************************/

package group

import (
	"errors"
	"fmt"
	log "im/common/log4go"
	"im/service/logic/models/group"
	"im/service/logic/models/user"
	"im/service/logic/models/utils"
	"im/service/logic/service/id"

	"gopkg.in/mgo.v2"
)

const MAX_CREATE_GROUPS = 16
const MAX_GROUP_MUMBERS = 1000

func GroupCreate(uid uint64, g *group.Group) (uint64, error) {

	log.Debug("1.检查这个ownerid 创建Grup的个数，进行限制")
	uext := user.UserExt{Uid: uid, Appid: g.Appid}
	err := uext.GetGroupInfo()
	if err != nil {
		return 0, err
	}

	if uext.CreateGroups != nil {
		count := len(uext.CreateGroups)
		if count > MAX_CREATE_GROUPS {
			return 0, errors.New(fmt.Sprintf("每个用户最多创建%d个群组", MAX_CREATE_GROUPS))
		}
	}

	log.Debug("2.插入群信息")
	gid := id.GetId("group")
	if gid == 0 {
		return 0, errors.New("get group id error")
	}

	g.Gid = gid
	g.Members = []uint64{uid}
	g.Version = 1
	g.Membercount = 1
	err = g.Insert()
	if err != nil {
		return 0, err
	}

	log.Debug("3.更新群组元数据")
	err = uext.CreateGroupUpdate(gid)
	if err != nil {
		g.Delete()
		return 0, err
	}

	return gid, nil
}

func GroupDelete(uid uint64, g *group.Group, notice bool) error {
	//1. 判断权限
	err := g.ReadWithMembers()
	if err != nil || g.Name == "" {
		return errors.New("未找到群组")
	} else {
		if g.Ownerid != uid {
			return errors.New("没有权限删除该群组")
		}
	}

	//2. 删除群信息
	err = g.Delete()
	if err != nil {
		return err
	}

	//3. 删除用户扩展表群信息
	uext := user.UserExt{Uid: uid, Appid: g.Appid}
	uext.DelGroupUpdate(g.Gid)

	//4. 发送群解散通知
	if notice {
		//push service 负责发送消息，发完消息后进行5
		//cmd := msg.Group_DELGROUP_NOTICE{uint32(g.Gid)}
		//notice_group(uint32(g.Gid), 0, uint32(uid), fbs.GroupMsgTypeMT_GROUP_DISGROUP, cmd.Serialize())
	}

	return nil
}

////申请加入群， 将向群主发送通知
//func AskForAdd_Group(uid, gid uint32) error {
//	return nil
//}

//添加组成员（这个函数应该由群的创建人调用）
func GroupAddMember(uid, adduid uint64, g *group.Group, notice bool) error {

	//1. 验证权限
	err := g.Read()
	if err != nil || g.Name == "" {
		return errors.New("未找到群组")
	} else {
		if g.Ownerid != uid {
			return errors.New("没有权限添加成员")
		}
		if uid == adduid {
			return errors.New("不能添加自己")
		}

		if g.Membercount > MAX_GROUP_MUMBERS {
			return errors.New("群成员达到上限")
		}
	}

	//2.将adduid添加到成员表
	err = g.AddMember(adduid)
	if err != nil && err != utils.ErrExist {
		return err
	}

	//3.更新adduid userext表
	uext := user.UserExt{Uid: adduid, Appid: g.Appid}
	err = uext.AddJoinGroup(g.Gid)
	if err != nil {
		e := g.DeleteMember(adduid)
		log.Error("[group|GroupAddMember|recove mem] gid:%d, uid:%d, %s", g.Gid, adduid, e.Error())
		return err
	}

	//发送通知
	if notice {
		//cmd := msg.Group_ADDMEMBER_NOTICE{uint32(gid), uint32(adduid)}
		//notice_group(uint32(gid), uint32(adduid), uint32(uid), fbs.GroupMsgTypeMT_GROUP_ADDMEMBER, cmd.Serialize())
	}
	return nil
}

// 删除成员
func GroupDelMember(uid, deluid uint64, g *group.Group, notice bool) error {

	//1. 验证权限
	err := g.Read()
	if err != nil || g.Name == "" {
		return errors.New("未找到群组")
	} else {
		if (g.Ownerid == uid) && (uid == deluid) {
			return errors.New("群主不能删除自己，需要解散群")
		}
		if (g.Ownerid != uid) && (uid != deluid) {
			return errors.New("没有权限删除其他人成员")
		}
	}

	//2.删除成员表成员
	err = g.DeleteMember(deluid)
	if err != nil && err != mgo.ErrNotFound {
		return err
	}

	//3.更新 userext表
	uext := user.UserExt{Uid: deluid, Appid: g.Appid}
	err = uext.DelJoinGroup(g.Gid)
	if err != nil {
		e := g.AddMember(deluid)
		log.Error("[group|GroupDelMember|recove mem] gid:%d, uid:%d, %s", g.Gid, deluid, e.Error())
		return err
	}

	//发送通知
	if notice {
		//cmd := msg.Group_DELMEMBER_NOTICE{uint32(gid), uint32(deluid)}
		//notice_group(uint32(gid), uint32(deluid), uint32(uid), fbs.GroupMsgTypeMT_GROUP_DELMEMBER, cmd.Serialize())
	}
	return nil
}

// 退群
func GroupQuit(uid uint64, g *group.Group, notice bool) error {

	//1. 验证权限
	err := g.Read()
	if err != nil || g.Name == "" {
		return errors.New("未找到群组")
	} else {
		if g.Ownerid == uid {
			return errors.New("群主不能删除自己，需要解散群")
		}
	}

	//2.删除成员表成员
	err = g.DeleteMember(uid)
	if err != nil {
		return err
	}

	//3.更新 userext表
	uext := user.UserExt{Uid: uid, Appid: g.Appid}
	err = uext.DelJoinGroup(g.Gid)
	if err != nil {
		e := g.AddMember(uid)
		log.Error("[group|GroupDelMember|recove mem] gid:%d, uid:%d, %s", g.Gid, uid, e.Error())
		return err
	}

	//发送通知
	if notice {
		//cmd := msg.Group_DELMEMBER_NOTICE{uint32(gid), uint32(uid)}
		//notice_group(uint32(gid), uint32(g.Ownerid), uint32(uid), fbs.GroupMsgTypeMT_GROUP_QUITGROUP, cmd.Serialize())
	}
	return nil
}
