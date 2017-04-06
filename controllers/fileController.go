package controllers

import (
	"fmt"
	"net/http"
	"net/url"
	. "nuobelcloud/conf"
	"nuobelcloud/models"
	. "nuobelcloud/tools"

	"strings"
)

/*

文件操作控制器，负责文件的上传下载删除等等和文件相关的操作

*/
//getfilelist
func OnGetFileList(w http.ResponseWriter, r *http.Request) {

	sessionInfo := CheckSession(w, r)
	if sessionInfo == nil {
		Log(ERROR, "错误session")
		return
	}
	//
	fileType := r.FormValue("filetype")
	filepath := r.FormValue("filepath")
	var Prefix string = fmt.Sprintf("%s/%s/%d", FileRootPath, sessionInfo.Dir1, sessionInfo.Uid)
	if fileType == "all" {
		modelFile := &models.ModelFile{
			Dir:  filepath,
			Type: fileType,
		}
		data, err := modelFile.GetFileList(Prefix)
		var reply *ReplyMess
		if err != nil {
			reply = &ReplyMess{
				Code: 1,
				Mess: "not find file",
				Data: "",
			}
			w.Write(reply.MessToJson())
			return
		}

		reply = &ReplyMess{
			Code: 0,
			Mess: "",
			Data: models.FileListToJson(data),
		}

		w.Write(reply.MessToJson())
	} else {
		var fl []*models.FileList
		mapFiletype := make(map[string]bool, 5)
		if fileType == "music" {
			mapFiletype["mp3"] = true
			mapFiletype["wav"] = true
			mapFiletype["ogg"] = true
		} else if fileType == "video" {
			mapFiletype["mp4"] = true
			// mapFiletype["wmv"] = true   //html5播放视频map4是所有现代浏览器都支持的

		} else if fileType == "doc" {
			mapFiletype["xls"] = true
			mapFiletype["doc"] = true
			mapFiletype["ppt"] = true
			mapFiletype["pdf"] = true
		} else if fileType == "image" {
			mapFiletype["png"] = true
			mapFiletype["jpg"] = true
			mapFiletype["jpeg"] = true
			mapFiletype["gif"] = true
			mapFiletype["bmp"] = true
		}
		var filepath = "" //有类型的filepath就是该用户根目录
		fl = models.ListForFiletype(Prefix, filepath, mapFiletype)

		if len(fl) <= 0 {
			reply := &ReplyMess{
				Code: 0,
				Mess: "no data",
				Data: "",
			}
			w.Write(reply.MessToJson())
		} else {
			reply := &ReplyMess{
				Code: 0,
				Mess: "",
				Data: models.FileListToJson(fl),
			}
			w.Write(reply.MessToJson())
		}

	}

}

//uploadfile
func OnUploadFile(w http.ResponseWriter, r *http.Request) {
	sessionInfo := CheckSession(w, r)
	if sessionInfo == nil {
		Log(ERROR, "错误session")
		return
	}
	var Prefix string = fmt.Sprintf("%s/%s/%d", FileRootPath, sessionInfo.Dir1, sessionInfo.Uid)
	uploadpath := r.FormValue("path")

	if strings.ContainsAny(uploadpath, ".. || \\") {
		reply := &ReplyMess{
			Code: 1,
			Mess: "路径不能包含/..等特殊字符",
			Data: "",
		}
		w.Write(reply.MessToJson())
		return
	}

	f, h, err := r.FormFile("fielnames")
	if err != nil {
		reply := &ReplyMess{
			Code: 1,
			Mess: "file err",
			Data: "",
		}
		w.Write(reply.MessToJson())
		Log(ERROR, "file err ", err)
		return
	}
	defer f.Close()
	filename := h.Filename
	err = models.UploadFile(Prefix, uploadpath, filename, f)
	if err != nil {
		reply := &ReplyMess{
			Code: 1,
			Mess: "file upload err!",
			Data: "",
		}
		w.Write(reply.MessToJson())
		Log(ERROR, "file upload err ", err)
		return
	}
	w.Write([]byte("success"))

}

//
func OnDownloadFile(w http.ResponseWriter, r *http.Request) {
	sessionInfo := CheckSession(w, r)
	if sessionInfo == nil {
		Log(ERROR, "错误session")
		return
	}
	var Prefix string = fmt.Sprintf("%s/%s/%d", FileRootPath, sessionInfo.Dir1, sessionInfo.Uid)
	var filepath = r.FormValue("path")
	filepath = Prefix + filepath
	var filename = r.FormValue("filename")
	if len(filepath) < 2 {
		reply := &ReplyMess{
			Code: 1,
			Mess: "param err",
			Data: "",
		}
		w.Write(reply.MessToJson())
		return
	}
	filepath += "/" + filename
	w.Header().Set("Content-Type", "application/octet-stream")
	//urlencode编码
	filename = url.QueryEscape(filename)
	w.Header().Set("Content-Disposition", "attachment; filename*=UTF-8''"+filename) //文件名utf8

	w.Header().Set("Content-Length", "1024")
	Log(INFO, "filepath=", filepath)

	http.ServeFile(w, r, filepath)
}

//
func OnDeleteFile(w http.ResponseWriter, r *http.Request) {
	sessionInfo := CheckSession(w, r)
	if sessionInfo == nil {
		Log(ERROR, "错误session")
		return
	}
	var Prefix string = fmt.Sprintf("%s/%s/%d", FileRootPath, sessionInfo.Dir1, sessionInfo.Uid)
	var filepath = r.FormValue("path")
	filepath = Prefix + filepath
	var filename = r.FormValue("filename")
	if len(filepath) < 1 {
		reply := &ReplyMess{
			Code: 1,
			Mess: "param err",
			Data: "",
		}
		w.Write(reply.MessToJson())
		return
	}
	filename = filepath + "/" + filename

	err := models.DeleteFile(filename)
	if err != nil {
		reply := &ReplyMess{
			Code: 1,
			Mess: "delete err!",
			Data: "",
		}
		w.Write(reply.MessToJson())
		Log(ERROR, "delete file err ", err)
		return
	}
	reply := &ReplyMess{
		Code: 0,
		Mess: "delete success!",
		Data: "",
	}
	w.Write(reply.MessToJson())
}

//
func OnMakeDir(w http.ResponseWriter, r *http.Request) {
	sessionInfo := CheckSession(w, r)
	if sessionInfo == nil {
		Log(ERROR, "错误session")
		return
	}
	var Prefix string = fmt.Sprintf("%s/%s/%d", FileRootPath, sessionInfo.Dir1, sessionInfo.Uid)
	path := r.FormValue("path")
	path = Prefix + path
	newdirname := r.FormValue("newdirname")
	path += "/" + newdirname
	if len(path) < 2 {
		reply := &ReplyMess{
			Code: 1,
			Mess: "param err",
			Data: "",
		}
		w.Write(reply.MessToJson())
		return
	}
	if strings.ContainsAny(path, ".. || \\") {
		reply := &ReplyMess{
			Code: 1,
			Mess: "文件夹名称不能包含..等特殊字符",
			Data: "",
		}
		w.Write(reply.MessToJson())
		return
	}

	err := models.MakeNewDir(path)
	if err != nil {
		reply := &ReplyMess{
			Code: 1,
			Mess: "make dir err",
			Data: "",
		}
		w.Write(reply.MessToJson())
		return
	}
	reply := &ReplyMess{
		Code: 0,
		Mess: "创建成功",
		Data: "",
	}
	w.Write(reply.MessToJson())
}

//重命名
func OnRenameFileHandle(w http.ResponseWriter, r *http.Request) {
	sessionInfo := CheckSession(w, r)
	if sessionInfo == nil {
		Log(ERROR, "错误session")
		return
	}
	filepath := r.FormValue("filepath")
	var Prefix string = fmt.Sprintf("%s/%s/%d", FileRootPath, sessionInfo.Dir1, sessionInfo.Uid)
	filepath = Prefix + filepath
	oldname := r.FormValue("oldname")
	newname := r.FormValue("newname")

	if len(oldname) < 1 || len(newname) < 1 {
		reply := &ReplyMess{
			Code: 1,
			Mess: "name error",
			Data: "",
		}
		w.Write(reply.MessToJson())
		return
	}

	old := filepath + "/" + oldname
	new := filepath + "/" + newname

	err := models.RenameFilename(old, new)
	if err != nil {
		reply := &ReplyMess{
			Code: 1,
			Mess: "rename fail",
			Data: "",
		}
		Log(ERROR, "rename fail", err)
		w.Write(reply.MessToJson())
		return
	}
	reply := &ReplyMess{
		Code: 0,
		Mess: "重命名成功",
		Data: "",
	}
	w.Write(reply.MessToJson())
}
