package main

import (
	"net/http"
	"nuobelcloud/controllers"
)

//入口

func main() {
	//静态文件
	http.Handle("/static/", http.FileServer(http.Dir("")))
	//登录
	http.HandleFunc("/", controllers.OnHome)
	http.HandleFunc("/nuobel", controllers.OnIndex)
	//关于我
	http.HandleFunc("/aboutme", controllers.OnAboutme)
	http.HandleFunc("/nuobel/login", controllers.OnLogin)
	//register
	http.HandleFunc("/nuobel/register", controllers.OnRegister)
	//index
	http.HandleFunc("/nuobel/index", controllers.OnIndex)
	//getuserinfo
	http.HandleFunc("/nuobel/getuserinfo", controllers.OnGetUserInfo)
	//getfilelist
	http.HandleFunc("/nuobel/getfilelist", controllers.OnGetFileList)
	//uploadfile
	http.HandleFunc("/nuobel/uploadfile", controllers.OnUploadFile)
	//download
	http.HandleFunc("/downloadfile", controllers.OnDownloadFile)
	//deletefile
	http.HandleFunc("/nuobel/deletefile", controllers.OnDeleteFile)
	//makedir
	http.HandleFunc("/nuobel/makedir", controllers.OnMakeDir)
	//exit
	http.HandleFunc("/nuobel/exit", controllers.OnExit)
	//rename
	http.HandleFunc("/nuobel/rename", controllers.OnRenameFileHandle)
	//监听
	http.ListenAndServe(":8080", nil)
}
