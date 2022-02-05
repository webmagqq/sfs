package wen

import (
	"html/template"
	"net/http"
	"sfsgo/gstr"
	"sfsgo/mysql"
	"sfsgo/pubgo"
	"strconv"
	"strings"
)

type wen struct {
	req   *http.Request
	Host  string
	Path  string
	Uri   string
	Up    string //繁体参数 ?l=0
	Ift   bool   //是否转繁体
	Wid   string
	Bt    string
	Jid   string
	Jname string
	Jtext template.HTML
}

func Newwen(req *http.Request) wen {
	//静态变量返回值用指针*,非静态类则返回copy一份，取决于new
	return wen{req: req}
}

func (q *wen) Setto(uif *pubgo.Urlinfo) {
	q.Host = uif.Host
	q.Path = uif.Path
	q.Uri = uif.Uri
	q.Up = uif.Up
	q.Ift = uif.Ift
}
func (q *wen) SetWid(Path string) {
	if strings.Contains(Path, "_") { //兼容旧版
		q.Wid = gstr.Do(Path, "/", "_", true, false)
	}
	if strings.Contains(Path, "-") {
		q.Wid = gstr.Do(Path, "/", "-", true, false)
	}
	if q.Wid == "" { //兼容旧版
		q.Wid = gstr.Do(Path, "/", "", true, false)
	}
}
func (q *wen) SetBt(Path string) string {
	if strings.Contains(Path, "_") { //兼容旧版
		return gstr.Do(Path, "_", "", true, true)
	}
	if strings.Contains(Path, "-") { //兼容旧版
		return gstr.Do(Path, "-", "", true, true)
	}
	return ""
}

func (q *wen) SetJtext(Wid string) template.HTML {
	intwid, _ := strconv.Atoi(Wid)
	bid := intwid - 7
	eid := intwid + 21
	rd := mysql.Selects("SELECT jid,nr from dzj WHERE tid>? and tid<?", bid, eid)
	lenrd := len(rd)
	if lenrd == 0 {
		return ""
	}
	nr := ""
	for _, v := range rd {
		nr += v["nr"]

	}
	q.Jid = rd[lenrd-1]["jid"]
	nr = strings.Replace(nr, "br", "<br>", -1)
	nr = strings.Replace(nr, "p0", "<br>", -1)
	nr = strings.Replace(nr, "p1", "<br>", -1)
	nr = strings.Replace(nr, q.Bt, "<b>"+q.Bt+"</b>", -1)
	return template.HTML(nr)
}
func (q *wen) SetJname() string {
	rd := mysql.Selects("SELECT jingming from jingmulu WHERE  tid=?", q.Jid)
	if len(rd) == 0 {
		return ""
	}
	return rd[0]["jingming"]

}
