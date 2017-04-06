package models

import (
	"fmt"
	. "nuobelcloud/tools"
	"os"
	"time"
)

/*

模型通用的一些功能

*/

//根据用户名和当前时间，生产sessionid
func MakeSession(UserName string) string {
	str := fmt.Sprintf("%s%d", UserName, time.Now().Unix())
	return Md5Hash(str)
}

//创建目录,可多级创建
func makeDir(dirpath string) error {
	return os.MkdirAll(dirpath, 0700)
}
