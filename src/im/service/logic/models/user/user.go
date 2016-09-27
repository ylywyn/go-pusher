/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   : user.go
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/

package user

import (
	"errors"
	"fmt"
	. "im/service/logic/models/utils"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	MAX_NAME_LEN = 24
)

//MogoDB 中的用户表模型
type Users struct {
	Uid           uint64
	Username      string
	Password      string
	Sex           string
	Email         string
	Appid         uint32
	Depid         int64
	Status        string
	Registertime  int64
	Nikename      string
	Tel           string
	Lastlogintime int64
	Role          string
	Address       string
	Admin         bool //true代表管理员
	Option        string
}

func (m *Users) TableName() string {
	return fmt.Sprintf("t_user_%d", m.Appid)
}

func (this *Users) Read() error {
	return UserDBPool.M(this.TableName(), func(c *mgo.Collection) error {
		return c.Find(bson.M{"uid": this.Uid}).One(this)
	})
}

func (this *Users) ReadByUserName() error {
	return UserDBPool.M(this.TableName(), func(c *mgo.Collection) error {
		return c.Find(bson.M{"username": this.Username}).One(this)
	})
}

func (this *Users) Delete() error {
	return UserDBPool.M(this.TableName(), func(c *mgo.Collection) error {
		return c.Remove(bson.M{"uid": this.Uid})
	})
}

func (this *Users) Update() error {
	return UserDBPool.M(this.TableName(), func(c *mgo.Collection) error {
		return c.Update(bson.M{"uid": this.Uid}, this)
	})
}

func (this *Users) Insert() error {
	return UserDBPool.M(this.TableName(), func(c *mgo.Collection) error {
		return c.Insert(this)
	})
}

//以后增加分页查询
func GetUsersInfo(appid, depid uint32) ([]Users, error) {
	var users []Users
	err := UserDBPool.M(fmt.Sprintf("t_user_%d", appid), func(c *mgo.Collection) error {
		return c.Find(nil).All(&users)
	})
	return users, err
}

func GetUids(appid, depid uint16, pageSize, pageIndex int) ([]Users, error) {
	var users []Users
	err := UserDBPool.M(fmt.Sprintf("t_user_%d", appid), func(c *mgo.Collection) error {
		return c.Find(nil).Select(bson.M{"uid": 1, "_id": 0}).All(&users)
	})
	return users, err
}

//分页查询(此方法假定uid按升序排列)
func GetUidsLimit(appid, depid uint16, lastid uint32, limit int) ([]Users, error) {
	var users []Users
	err := UserDBPool.M(fmt.Sprintf("t_user_%d", appid), func(c *mgo.Collection) error {
		return c.Find(bson.M{"uid": bson.M{"$gt": lastid}}).Select(bson.M{"uid": 1, "_id": 0}).Limit(limit).All(&users)
	})
	return users, err
}

//如果可用返回true
func CheckUserName(username string, appid uint32) (rst bool, err error) {

	l := len(username)
	if l > 0 && l < MAX_NAME_LEN {
		user := &Users{Username: username, Appid: appid}
		err := user.ReadByUserName()
		if err != nil {
			if err == mgo.ErrNotFound {
				return true, nil
			} else {
				return false, err
			}
		} else {
			return false, errors.New("用户名已被占用")
		}
	} else {
		return false, errors.New("参数错误")
	}
}
