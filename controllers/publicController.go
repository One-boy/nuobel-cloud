package controllers

//控制器公用
import (
	"fmt"
	"net/http"
	. "nuobelcloud/conf"
	"nuobelcloud/session"
	. "nuobelcloud/tools"
)

/*

控制器通用的一些操作在这里

*/

//统一的消息回复格式

type ReplyMess struct {
	Code int
	Mess string
	Data string
}

func (c *ReplyMess) MessToJson() []byte {
	if len(c.Data) > 0 {

		return []byte(fmt.Sprintf("{\"code\":%d,\"mess\":\"%s\",\"data\":%s}", c.Code, c.Mess, c.Data))
	}
	return []byte(fmt.Sprintf("{\"code\":%d,\"mess\":\"%s\",\"data\":\"%s\"}", c.Code, c.Mess, c.Data))
}

//session认证
func CheckSession(w http.ResponseWriter, r *http.Request) *session.SessionInfo {

	//判断cookie
	cookie, err := r.Cookie(CookieName)
	//没有cookie时
	if err != nil {
		Log(ERROR, "无cookie信息", err)
		http.Redirect(w, r, LoginUrl, 302)
		return nil
	}
	cookieVal := cookie.Value
	sessionInfo := &session.SessionInfo{Sid: cookieVal}
	err = sessionInfo.Get()
	//没有此session
	if err != nil {
		Log(ERROR, "未找到session", err)
		http.Redirect(w, r, LoginUrl, 302)
		return nil
	}
	if len(sessionInfo.UserName) <= 0 {
		Log(ERROR, "未获取到session用户名")
		http.Redirect(w, r, LoginUrl, 302)
		return nil
	}

	return sessionInfo
}
