package controllers

import (
	"net/http"
	"nuobelcloud/models"
	. "nuobelcloud/tools"
	"strings"
)

/*

注册控制器，负责注册相关功能

*/
//register

func OnRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		pass := r.FormValue("pass")
		repass := r.FormValue("repass")
		email := r.FormValue("email")

		if len(username) == 0 || len(pass) == 0 {
			reply := &ReplyMess{1, "用户名或密码不能为空", ""}
			w.Write(reply.MessToJson())
			return
		}
		if pass != repass {
			reply := &ReplyMess{1, "两次密码不相同", ""}
			w.Write(reply.MessToJson())
			return
		}
		if strings.ContainsAny(username, "/ || \\ || , || @ || ~") {
			reply := &ReplyMess{1, "用户名不能包含/\\,@~等特殊字符", ""}
			w.Write(reply.MessToJson())
			return
		}
		if strings.ContainsAny(pass, "/ || \\ || , || @ || ~") {
			reply := &ReplyMess{1, "密码不能包含/\\,@~等特殊字符", ""}
			w.Write(reply.MessToJson())
			return

		}
		//注册信息
		userInfo := &models.UserInfo{
			UserName: username,
			Pass:     pass,
			Email:    email,
		}

		err := userInfo.RegisterHandle()
		var reply = new(ReplyMess)
		if err != nil {
			reply = &ReplyMess{1, err.(*MyError).GetDetail(), ""}
		} else {
			reply = &ReplyMess{0, "注册成功", ""}
		}
		w.Write(reply.MessToJson())
	}
}
