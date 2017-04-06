package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	. "nuobelcloud/tools"
)

/*

首页控制器，负责首页html模板和除文件列表操作外的其它信息获取

*/
//index
func OnIndex(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		sessionInfo := CheckSession(w, r)
		if sessionInfo == nil {
			return
		}

		//认证通过
		temp, err := template.ParseFiles("views/nuobel_index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		temp.Execute(w, nil)
	}
}

//getuserinfo
func OnGetUserInfo(w http.ResponseWriter, r *http.Request) {

	sessionInfo := CheckSession(w, r)
	Log(INFO, "请求index3", sessionInfo.UserName)
	if sessionInfo == nil {
		return
	}
	var reply *ReplyMess
	reply = &ReplyMess{
		Code: 0,
		Mess: "",
		Data: fmt.Sprintf("{\"username\":\"%s\"}", sessionInfo.UserName),
	}
	w.Write(reply.MessToJson())

}
