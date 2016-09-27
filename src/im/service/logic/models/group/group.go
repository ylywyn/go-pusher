/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   :
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/

package group

import (
	"fmt"

	. "im/service/logic/models/utils"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//MongoDB中的群表模型
type Group struct {
	Gid         uint64
	Appid       uint32
	Creattime   int64
	Ownerid     uint64
	Name        string
	Desc        string
	Membercount int64
	Version     int64
	Members     []uint64
}

func (this *Group) TableName() string {
	return fmt.Sprintf("t_group_%d", this.Appid)
}

//不读取用户群成员
func (this *Group) Read() error {
	return UserDBPool.M(this.TableName(), func(c *mgo.Collection) error {
		return c.Find(bson.M{"gid": this.Gid}).Select(bson.M{"members": 0}).One(this)
	})
}

func (this *Group) ReadWithMembers() error {
	return UserDBPool.M(this.TableName(), func(c *mgo.Collection) error {
		return c.Find(bson.M{"gid": this.Gid}).One(this)
	})
}

//删除群组
func (this *Group) Delete() error {
	return UserDBPool.M(this.TableName(), func(c *mgo.Collection) error {
		return c.Remove(bson.M{"gid": this.Gid})
	})
}

//func (this *Group) Update() error {
//	return UserDBPool.M(this.TableName(), func(c *mgo.Collection) error {
//		return c.Update(bson.M{"gid": this.Gid}, this)
//	})
//}

//插入一条
func (this *Group) Insert() error {
	return UserDBPool.M(this.TableName(), func(c *mgo.Collection) error {
		return c.Insert(this)
	})
}

//群组成员添加记录
func (this *Group) AddMember(member uint64) error {

	ret := UserDBPool.M(this.TableName(), func(c *mgo.Collection) error {
		count, err := c.Find(bson.M{"gid": this.Gid, "members": member}).Count()
		if err != nil {
			return err
		}
		if count != 0 {
			return ErrExist
		}

		return c.Update(bson.M{"gid": this.Gid},
			bson.M{
				"$push": bson.M{"members": member},
				"$inc":  bson.M{"version": 1, "membercount": 1},
			})
	})

	return ret
}

// 删除成员
func (this *Group) DeleteMember(member uint64) error {

	ret := UserDBPool.M(this.TableName(), func(c *mgo.Collection) error {
		count, err := c.Find(bson.M{"gid": this.Gid, "members": member}).Count()
		if err != nil {
			return err
		}
		if count == 0 {
			return mgo.ErrNotFound
		}

		return c.Update(bson.M{"gid": this.Gid},
			bson.M{
				"$pull": bson.M{"members": member},
				"$inc":  bson.M{"version": 1, "membercount": -1},
			})
	})

	return ret
}

//获取成员列表
func (this *Group) GetMembers() error {
	return UserDBPool.M(this.TableName(), func(c *mgo.Collection) error {
		return c.Find(bson.M{"gid": this.Gid}).Select(bson.M{"members": 1, "_id": 0}).One(this)
	})
}
