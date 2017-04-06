package models

import (
	. "nuobelcloud/conf"
	sql "nuobelcloud/mysql"
	"nuobelcloud/session"
	. "nuobelcloud/tools"
	"time"
)

/*

登陆相关数据处理

*/

//用户信息结构
type UserInfo struct {
	Uid      int
	UserName string
	Pass     string
	Email    string
}

//登陆操作
func (c *UserInfo) LoginHandle() (*session.SessionInfo, error) {

	row := sql.DB.QueryRow("select uid,dir1,password from user where username=?", c.UserName)

	var uid int
	var dir1 string
	var oldPass string
	row.Scan(&uid, &dir1, &oldPass)

	if len(oldPass) < 1 {
		return nil, &MyError{"LoginHandle", "没有此用户", nil}
	}
	newPass := Md5Hash(c.Pass + Md5Salt)
	if newPass != oldPass {
		return nil, &MyError{"LoginHandle", "密码验证不通过", nil}
	}

	//加入session
	sessionInfo := &session.SessionInfo{
		Sid:        MakeSession(c.UserName),
		Uid:        uid,
		UserName:   c.UserName,
		SetDate:    time.Now().Unix(),
		ExpireTime: CookieExpireTime,
		Dir1:       dir1,
	}
	return sessionInfo, nil
}
