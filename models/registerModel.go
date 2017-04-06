package models

import (
	"fmt"
	. "nuobelcloud/conf"
	sql "nuobelcloud/mysql"
	. "nuobelcloud/tools"
	"time"
)

/*

注册操作相关数据处理

*/

//注册操作
func (c *UserInfo) RegisterHandle() error {
	//查询是否已经被注册
	row := sql.DB.QueryRow("select username from user where username=?", c.UserName)

	var username string
	row.Scan(&username)
	if len(username) >= 1 {
		return &MyError{"RegisterHandle", "此用户已被注册", nil}
	}

	//插入数据库
	date := time.Now().Format("20060102")
	md5Pass := Md5Hash(c.Pass + Md5Salt)
	result, err := sql.DB.Exec("insert into user(username,password,email,dir1,createtime) values(?,?,?,?,?)", c.UserName, md5Pass, c.Email, date, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		return &MyError{"RegisterHandle", "服务器出故障啦，请稍后再试", err}
	}
	//创建文件夹
	id, _ := result.LastInsertId()
	filepath := fmt.Sprintf("%s/%s/%d", FileRootPath, date, id)
	err = makeDir(filepath)
	if err != nil {
		return &MyError{"RegisterHandle", "服务器出故障啦，请稍后再试", err}
	}
	return nil
}
