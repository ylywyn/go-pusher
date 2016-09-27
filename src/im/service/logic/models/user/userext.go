/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   : userext.go
 *  Date   :
 *  Author : yangl
 *  Description:  Friends and Groups Info
 ******************************************************************/

package user

import (
	"fmt"
	. "im/service/logic/models/utils"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//用户扩展表 实体
type UserExt struct {
	Uid            uint64   `bson:"uid"`
	Appid          uint32   `bson:"appid"`
	Groups         []uint64 `bson:"groups"`
	CreateGroups   []uint64 `bson:"creategroups"`
	Friends        []uint64 `bson:"friends"`
	FriendsVersion int32    `bson:"friendsversion"`
}

//获取MongoDB的用户扩展表名
func (this *UserExt) TableName() string {
	return fmt.Sprintf("t_userext_%d", this.Appid)
}

func (this *UserExt) Get() error {
	return UserDBPool.M(this.TableName(), func(c *mgo.Collection) error {
		return c.Find(bson.M{"uid": this.Uid}).Select(bson.M{"friends": 0}).One(this)
	})
}

func (this *UserExt) Delete() error {
	return UserDBPool.M(this.TableName(), func(c *mgo.Collection) error {
		err := c.Remove(bson.M{"uid": this.Uid})
		return err
	})
}

//func (this *UserExt) Update() error {
//	return UserDBPool.M(this.TableName(), func(c *mgo.Collection) error {
//		return c.Update(bson.M{"uid": this.Uid}, this)
//	})
//}

func (this *UserExt) Insert() error {
	return UserDBPool.M(this.TableName(), func(c *mgo.Collection) error {
		return c.Insert(this)
	})
}

//记录不存在则为插入，存在则为更新
func (this *UserExt) AddJoinGroup(g uint64) error {

	return UserDBPool.M(this.TableName(), func(c *mgo.Collection) error {
		return c.Update(bson.M{"uid": this.Uid}, bson.M{"$addToSet": bson.M{
			"groups": g,
		}})
	})
}

//删除用户扩展表中的 其在的group
func (this *UserExt) DelJoinGroup(g uint64) error {

	return UserDBPool.M(this.TableName(), func(c *mgo.Collection) error {
		return c.Update(bson.M{"uid": this.Uid}, bson.M{"$pull": bson.M{
			"groups": g,
		}})
	})
}

//记录不存在则为插入，存在则为更新
func (this *UserExt) AddCreateGroup(g uint64) error {
	return UserDBPool.M(this.TableName(), func(c *mgo.Collection) error {
		return c.Update(bson.M{"uid": this.Uid}, bson.M{"$addToSet": bson.M{
			"creategroups": g,
		}})
	})
}

//删除用户扩展表中的 其在的group
func (this *UserExt) DelCreateGroup(g uint64) error {
	return UserDBPool.M(this.TableName(), func(c *mgo.Collection) error {
		return c.Update(bson.M{"uid": this.Uid}, bson.M{"$pull": bson.M{
			"creategroups": g,
		}})
	})
}

//增加用户扩展表中的 其在的group, AddCreateGroup, AddJoinGroup
func (this *UserExt) CreateGroupUpdate(g uint64) error {
	return UserDBPool.M(this.TableName(), func(c *mgo.Collection) error {

		return c.Update(bson.M{"uid": this.Uid}, bson.M{"$addToSet": bson.M{
			"creategroups": g, "groups": g,
		}})

	})
}

//删除用户扩展表中的 其在的group, 合并DelCreateGroup, DelJoinGroup两个删除操作
func (this *UserExt) DelGroupUpdate(g uint64) error {
	return UserDBPool.M(this.TableName(), func(c *mgo.Collection) error {
		return c.Update(bson.M{"uid": this.Uid}, bson.M{"$pull": bson.M{
			"creategroups": g, "groups": g,
		}})
	})
}

//获取创建的组
func (this *UserExt) GetGroupInfo() error {
	return UserDBPool.M(this.TableName(), func(c *mgo.Collection) error {
		n, err := c.Find(bson.M{"uid": this.Uid}).Count()
		if err != nil {
			return err
		}
		if n == 0 {
			this.Insert()
		}
		return c.Find(bson.M{"uid": this.Uid}).Select(bson.M{"friends": 0}).One(this)
	})
}

//添加好友, 修改了版本
func (this *UserExt) AddFriends(uid uint64) error {
	return UserDBPool.M(this.TableName(), func(c *mgo.Collection) error {
		n, err := c.Find(bson.M{"uid": this.Uid, "friends": uid}).Count()
		if err != nil {
			return err
		}
		if n != 0 {
			return ErrExist
		}
		return c.Update(bson.M{"uid": this.Uid},
			bson.M{
				"$push": bson.M{"friends": uid},
				"$inc":  bson.M{"friendsversion": 1},
			})
	})
}

//删除好友，, 修改了版本
func (this *UserExt) RemoveFriends(uid uint64) error {
	return UserDBPool.M(this.TableName(), func(c *mgo.Collection) error {
		n, err := c.Find(bson.M{"uid": this.Uid, "friends": uid}).Count()
		if err != nil {
			return err
		}
		if n == 0 {
			return mgo.ErrNotFound
		}
		return c.Update(bson.M{"uid": this.Uid},
			bson.M{
				"$pull": bson.M{"friends": uid},
				"$inc":  bson.M{"friendsversion": 1},
			})
	})
}

//获取好友关系版本号
func (this *UserExt) GetFriendsVersion() error {
	return UserDBPool.M(this.TableName(), func(c *mgo.Collection) error {
		err := c.Find(bson.M{"uid": this.Uid}).Select(bson.M{"friendsversion": 1, "_id": 0}).One(this)
		return err
	})
}

func (this *UserExt) GetFriends() error {
	return UserDBPool.M(this.TableName(), func(c *mgo.Collection) error {
		err := c.Find(bson.M{"uid": this.Uid}).Select(bson.M{"friends": 1, "_id": 0}).One(this)
		return err
	})
}
