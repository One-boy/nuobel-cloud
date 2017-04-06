package controllers

import (
	"html/template"
	"net/http"
	. "nuobelcloud/conf"
	"nuobelcloud/models"
	. "nuobelcloud/tools"
)

/*

登陆控制器，负责登陆功能

*/
//登录
func OnLogin(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {

		temp, err := template.ParseFiles("views/nuobel_login.html")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		temp.Execute(w, nil)
	} else if r.Method == "POST" {

		username := r.FormValue("username")
		pass := r.FormValue("pass")
		autologin := r.FormValue("autologin")
		cookie, err := r.Cookie(CookieName)
		updateSession := false

		if err != nil {
			Log(ERROR, "老cookie不存在", err)
		} else {
			Log(INFO, "老cookie=", cookie.Value, "err=", err)
			updateSession = true
		}

		if len(username) == 0 || len(pass) == 0 {
			reply := &ReplyMess{1, "用户名或密码不能为空", ""}
			w.Write(reply.MessToJson())
			return
		}

		userInfo := new(models.UserInfo)
		userInfo.UserName = username
		userInfo.Pass = pass

		sessioninfo, err := userInfo.LoginHandle()
		if err != nil {
			reply := &ReplyMess{1, err.(*MyError).GetDetail(), ""} //这里获取错误详细信息，需要强转成我们自己定义的错误类型
			w.Write(reply.MessToJson())
			return
		}
		//如果老的cookie在，更新
		if updateSession {
			err = sessioninfo.Update(cookie.Value)
			if err != nil {
				reply := &ReplyMess{1, "服务器错误", ""}
				w.Write(reply.MessToJson())
				return
			}
		} else {
			err = sessioninfo.Set()
			if err != nil {
				reply := &ReplyMess{1, "服务器错误", ""}
				w.Write(reply.MessToJson())
				return
			}
		}
		reply := &ReplyMess{0, "登录成功", ""}
		//设置cookie
		Log(INFO, "autologin=", autologin)
		var maxage int
		//maxage=0表示cookie的周期为session，>0表示到指定时间过期
		if autologin == "false" {
			maxage = 0
		} else {
			maxage = sessioninfo.ExpireTime
		}
		setcookie := &http.Cookie{
			Name:   CookieName,
			Value:  sessioninfo.Sid,
			Path:   "/",
			MaxAge: maxage,
		}
		http.SetCookie(w, setcookie)
		w.Write(reply.MessToJson())
	}
}

//退出
func OnExit(w http.ResponseWriter, r *http.Request) {
	sessionInfo := CheckSession(w, r)
	if sessionInfo == nil {
		Log(ERROR, "错误session")
		return
	}
	sessionInfo.Del()
	reply := &ReplyMess{0, "退出成功", ""}
	w.Write(reply.MessToJson())
}
