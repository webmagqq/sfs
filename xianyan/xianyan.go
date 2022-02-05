package xianyan

import (
	"fmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"sfsgo/gstr"
	"sfsgo/pubgo"
	"sort"
	"strconv"
	"strings"
)

type xianyan struct {
	req  *http.Request
	Host string
	Path string
	Uri  string
	Up   string //繁体参数 ?l=0
	Ift  bool   //是否转繁体
	P    int
	dir  string
	fs   []fs.FileInfo
	Text template.HTML
	Pg   template.HTML
}

func Newxianyan(req *http.Request) xianyan {
	//静态变量返回值用指针*,非静态类则返回copy一份，取决于new
	return xianyan{req: req, dir: "./xianyan/wen"}
}

func (x *xianyan) Setto(uif *pubgo.Urlinfo) {
	x.Host = uif.Host
	x.Path = uif.Path
	x.Uri = uif.Uri
	x.Up = uif.Up
	x.Ift = uif.Ift
}
func (x *xianyan) SetP() {
	P := gstr.Do(x.Path, "/", "", true, false)
	if P == "" {
		x.P = 0
	}
	ip, err := strconv.Atoi(P)
	if err != nil {
		x.P = 0
	} else {
		x.P = ip - 1
	}
}
func (x *xianyan) Setfs() {
	x.fs, _ = ReadDir(x.dir, false)
}
func (x *xianyan) SetText() {
	if len(x.fs) == 0 || len(x.fs) <= x.P {
		return
	}
	filetext, err := ioutil.ReadFile(x.dir + "/" + x.fs[x.P].Name())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	ftex := pubgo.Jf(string(filetext), x.Ift)
	ftex = strings.Replace(ftex, "\r\n", "<br>", -1)
	x.Text = template.HTML(ftex)
}
func ReadDir(dirname string, st bool) ([]fs.FileInfo, error) { //--打开文件夹，修改ioutil.ReadDir(dirPth)的方法
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Slice(list, func(i, j int) bool {
		if st { //--升序
			return list[i].Name() < list[j].Name()
		} else { //倒序
			return list[i].Name() > list[j].Name()
		}

	})
	return list, nil
}
func (x *xianyan) SetPg() {
	if len(x.fs) < 2 {
		return
	}
	rs, p := "", ""
	for i, _ := range x.fs {
		p = strconv.Itoa(i + 1)
		if i == x.P {
			rs += "<a href=\"javascript:\"><span class=\"btn btn-outline-secondary\">" + p + "</span></a>"
		} else {
			rs += "<a href=\"" + x.Host + "/xianyan/" + p + "\"><span class=\"btn btn-outline-primary\">" + p + "</span></a>"
		}

	}
	x.Pg = template.HTML(rs)
}
