package tools

import (
	"crypto/md5"
	"fmt"
	"time"
)

/*

一些工具，比如打印，md5摘要等。

*/

const (
	INFO  = 0
	WARN  = 1
	ERROR = 2
)

//自定义错误
type MyError struct {
	Op     string
	Detail string
	Err    error
}

//实现Error方法，MyError才实现了error接口
func (c *MyError) Error() string {
	if c.Err != nil {
		Log(ERROR, "Handle:"+c.Op+" Detail:"+c.Detail+" Error:"+c.Err.Error())
		return "Handle:" + c.Op + " Detail:" + c.Detail + " Error:" + c.Err.Error()
	}
	Log(ERROR, "Handle:"+c.Op+" Detail:"+c.Detail)
	return "Handle:" + c.Op + " Detail:" + c.Detail
}
func (c *MyError) GetDetail() string {
	return c.Detail
}
func Log(level int, args ...interface{}) {
	logTime := time.Now().Format("2006-01-02 15:04:05")
	switch level {
	case INFO:
		fmt.Println("[INFO]", logTime, ":", args)
	case WARN:
		fmt.Println("[WARN]", logTime, ":", args)
	case ERROR:
		fmt.Println("[ERROR]", logTime, ":", args)
	default:
		fmt.Println(args)
	}
}

func Md5Hash(src string) string {

	hash := md5.New()
	hash.Write([]byte(src))

	return fmt.Sprintf("%x", hash.Sum(nil))
}
