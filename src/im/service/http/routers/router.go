/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   :
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/

package routers

import (
	"im/service/http/controllers"
	"im/service/http/controllers/push"
	"im/service/http/controllers/users"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/service/v1/test", &controllers.TestController{}, "get,post:GetTest")

	//	// push
	beego.Router("/service/v1/push", &push.PushController{}, "post:Push")
	//	//beego.Router("/service/v1/push/settoken", &push.PushController{}, "get,post:SetToken")
	//	//beego.Router("/service/v1/push/deltoken", &push.PushController{}, "get,post:DelToken")

	//	//用户注册
	beego.Router("/service/v1/user/signup", &users.UserController{}, "post:SignUp")
	beego.Router("/service/v1/user/login", &users.UserController{}, "*:Login")
	beego.Router("/service/v1/user/loginbyid", &users.UserController{}, "*:LoginByUid")
	beego.Router("/service/v1/user/logout", &users.UserController{}, "get,post:Logout")
	beego.Router("/service/v1/user/checkusername", &users.UserController{}, "get,post:CheckUserName")
	beego.Router("/service/v1/user/userinfo", &users.UserController{}, "get:GetUserInfo")
	beego.Router("/service/v1/user/rm", &users.UserController{}, "get,post:RemoveUser")
	//	beego.Router("/service/v1/user/getusersinfo", &users.UserController{}, "get:GetUsersInfo")
	beego.Router("/service/v1/user/updatepassword", &users.UserController{}, "post:UpdatePassword")
	//	beego.Router("/service/v1/user/updateinfo", &users.UserController{}, "get,post:UpdateUserInfo")
	//	beego.Router("/service/v1/user/updateusername", &users.UserController{}, "get,post:UpdateUsername")
	//	beego.Router("/service/v1/user/getuserinfofornickname", &users.UserController{}, "get:GetUserInfosForNickname")
	//	beego.Router("/service/v1/user/getuserinfoforusername", &users.UserController{}, "get:GetUserInfosForUsername")

	//	beego.Router("/service/v1/user/checksessionid", &users.UserController{}, "get,post:CheckSessionId")
	//	beego.Router("/service/v1/user/statisticslogin", &users.UserController{}, "post:StatisticsLogin")

	//	//好友
	beego.Router("/service/v1/user/friends/add", &users.UserController{}, "*:AddFriend")
	beego.Router("/service/v1/user/friends/rm", &users.UserController{}, "*:RemoveFriend")
	beego.Router("/service/v1/user/friends/all", &users.UserController{}, "*:GetFriendsList")
	beego.Router("/service/v1/user/friends/version", &users.UserController{}, "*:GetFriendVersion")

	//	//用户群组
	beego.Router("/service/v1/user/group/create", &users.GroupController{}, "post:CreateGroup")
	beego.Router("/service/v1/user/group/getinfo", &users.GroupController{}, "get:GetGroupInfo")
	beego.Router("/service/v1/user/group/getjoins", &users.GroupController{}, "get:GetJoinGroups")
	beego.Router("/service/v1/user/group/getcreates", &users.GroupController{}, "get:GetCreateGroups")
	beego.Router("/service/v1/user/group/delete", &users.GroupController{}, "get:DeleteGroup")
	beego.Router("/service/v1/user/group/quit", &users.GroupController{}, "get:QuitGroup")
	beego.Router("/service/v1/user/group/addmember", &users.GroupController{}, "get,post:AddGroupMember")
	beego.Router("/service/v1/user/group/addmembers", &users.GroupController{}, "get,post:AddGroupMembers")
	beego.Router("/service/v1/user/group/rmmembers", &users.GroupController{}, "get,post:RemoveGroupMembers")
	beego.Router("/service/v1/user/group/getallmembers", &users.GroupController{}, "get,post:GetGroupMemberList")
	beego.Router("/service/v1/user/group/version", &users.GroupController{}, "get,post:GetGroupMemberList")
	//beego.Router("/service/v1/user/group/updateinfo", &users.GroupController{}, "get,post:UpdateGroupInfo")

	//	//历史消息
	//	beego.Router("/service/v1/user/historymsg/broadcast", &users.HistoryMsgController{}, "get,post:GetBroadcastHistoryMsg")
	//	beego.Router("/service/v1/user/historymsg/personal", &users.HistoryMsgController{}, "get,post:GetPersonalHistoryMsg")
	//	//离线消息
	//	beego.Router("/service/v1/user/offlinemsg/broadcast", &users.HistoryMsgController{}, "get,post:GetBroadcastHistoryMsg")
	//	beego.Router("/service/v1/user/offlinemsg/personal", &users.HistoryMsgController{}, "get,post:GetPersonalOffLineMsg")
	//	//离线消息数量
	//	beego.Router("/service/v1/user/offlinemsgcount/broadcast", &users.HistoryMsgController{}, "get,post:GetBroadcastOffLineMsgCount")
	//	beego.Router("/service/v1/user/offlinemsgcount/personal", &users.HistoryMsgController{}, "get,post:GetPersonalOffLineMsgCount")

	//	//App
	//	beego.Router("/service/v1/appshop/addapp", &appshop.AppshopController{}, "post:AddApp")
	//	beego.Router("/service/v1/appshop/removeapp", &appshop.AppshopController{}, "get,post:RemoveApp")
	//	beego.Router("/service/v1/appshop/updateapp", &appshop.AppshopController{}, "post:UpdateApp")
}
