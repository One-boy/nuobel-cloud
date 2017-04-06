package models

import (
	//"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	. "nuobelcloud/tools"
	"os"
	"regexp"
	"strings"
)

/*

文件操作相关数据处理

*/

//要获取的文件信息
type FileList struct {
	Name     string
	Path     string
	Size     float32
	Lasttime string
	Type     string
}

//要获取文件需要的信息
type ModelFile struct {
	Dir  string
	Type string
}

func (c *ModelFile) GetFileList(prefix string) ([]*FileList, error) {

	var filelist []*FileList
	fileinfo, err := ioutil.ReadDir(prefix + c.Dir)
	if err != nil {
		Log(ERROR, "dir=", c.Dir, "type=", c.Type, err)
		return nil, err
	}
	for _, v := range fileinfo {
		fl := new(FileList)
		fl.Name = v.Name()

		fl.Lasttime = v.ModTime().Format("2006-01-02 15:04:05")
		if v.IsDir() {
			fl.Size = 0
			fl.Type = "dir"
		} else {
			fl.Size = float32(v.Size())
			//提取后缀名
			reg := regexp.MustCompile(`\.[A-Za-z0-9]+$`)
			filetype := reg.FindAllString(v.Name(), -1)
			if len(filetype) > 0 {
				fl.Type = strings.Replace(filetype[0], ".", "", 1)
			} else {
				fl.Type = "unknown"
			}
		}
		fl.Path = c.Dir
		filelist = append(filelist, fl)
	}
	return filelist, nil
}

//文件列表转字符串
func FileListToJson(data []*FileList) string {
	var str string
	var datalen = len(data)
	for k, v := range data {
		if k < datalen && str != "" {
			str += fmt.Sprintf(",{\"name\":\"%s\",\"type\":\"%s\",\"time\":\"%s\",\"size\":%f,\"path\":\"%s\"}", v.Name, v.Type, v.Lasttime, v.Size, v.Path)
		} else {
			str += fmt.Sprintf("{\"name\":\"%s\",\"type\":\"%s\",\"time\":\"%s\",\"size\":%f,\"path\":\"%s\"}", v.Name, v.Type, v.Lasttime, v.Size, v.Path)
		}
	}
	str = fmt.Sprintf("[%s]", str)
	return str
}

//上传处理

func UploadFile(prefix string, filepath string, filename string, file multipart.File) error {
	Log(INFO, "上传路径=", filepath, "上传文件名:", filename)
	t, err := os.Create(prefix + filepath + "/" + filename)
	if err != nil {
		return err
	}
	defer file.Close()
	defer t.Close()
	_, err = io.Copy(t, file)
	if err != nil {
		return err
	}
	return nil
}

//新建文件夹
func MakeNewDir(dirpath string) error {
	return makeDir(dirpath)
}

//删除文件
func DeleteFile(path string) error {

	return os.Remove(path)
}

//文件更名
func RenameFilename(oldpath, newpath string) error {

	return os.Rename(oldpath, newpath)
}

//列出制定类型文件
func ListForFiletype(prefix string, filepath string, searchtype map[string]bool) []*FileList {
	var filelist []*FileList
	fileinfo, err := ioutil.ReadDir(prefix + filepath)
	if err != nil {
		Log(ERROR, "ioutil readdir err , handle ListForFiletype,filepath=", filepath)
		filelist = nil
		return nil
	}
	for _, v := range fileinfo {
		if v.IsDir() {
			filelist = append(filelist, ListForFiletype(prefix, filepath+"/"+v.Name(), searchtype)...)
		} else {
			//提取后缀名
			reg := regexp.MustCompile(`\.[A-Za-z0-9]+$`)
			filetype := reg.FindAllString(v.Name(), -1)
			if len(filetype) > 0 {
				if _, ok := searchtype[strings.Replace(filetype[0], ".", "", 1)]; ok {
					fl := new(FileList)
					fl.Name = v.Name()
					fl.Size = float32(v.Size())
					fl.Lasttime = v.ModTime().Format("2006-01-02 15:04:05")
					fl.Path = filepath
					fl.Type = strings.Replace(filetype[0], ".", "", 1)
					filelist = append(filelist, fl)
				}
			}

		}
	}
	return filelist
}
