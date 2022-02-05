package tongji

import (
	"fmt"
	"html/template"
	"net/http"
	"sfsgo/mysql"
	"sfsgo/pubgo"
	"strconv"
)

type tongji struct {
	req   *http.Request
	Host  string
	Path  string
	Uri   string
	Up    string //繁体参数 ?l=0
	Ift   bool   //是否转繁体
	P     int
	Ct    string
	d     string
	Rhtml template.HTML
}

func Newtongji(req *http.Request) tongji {
	//静态变量返回值用指针*,非静态类则返回copy一份，取决于new
	return tongji{req: req}
}

func (t *tongji) Setto(uif *pubgo.Urlinfo) {
	t.Host = uif.Host
	t.Path = uif.Path
	t.Uri = uif.Uri
	t.Up = uif.Up
	t.Ift = uif.Ift
}
func (t *tongji) SetP() {
	P := t.req.URL.Query().Get("p")
	if P == "" {
		t.P = 1
	}
	ip, err := strconv.Atoi(P)
	if err != nil {
		t.P = 1
	} else {
		t.P = ip
	}
}
func (t *tongji) Trtb() {
	t.d = t.req.URL.Query().Get("d")
	if t.d != "" {
		rd := mysql.Exesql("truncate table llantji")
		if rd == -1 {
			fmt.Println("错误")
		}
	}
}
func (t *tongji) SetRhtml() {
	pg := 100
	lm := (t.P - 1) * pg
	rstr := ""
	var rd []map[string]string
	wl := t.req.URL.Query().Get("wl")
	if wl == "" {
		rd = mysql.Selects("SELECT lxing,idkey,llan,purl,ip,rq from llantji where tid<=(select tid from llantji ORDER BY tid DESC limit ?,1) ORDER BY tid DESC limit 100", lm)
	} else {
		rd = mysql.Selects("SELECT lxing,idkey,llan,purl,ip,rq from llantji where tid<=(select tid from llantji WHERE lxing=? ORDER BY tid DESC limit ?,1) and  lxing=?  ORDER BY tid DESC limit 100", wl, lm, wl)
	}
	for _, v := range rd {
		rstr += "<tr>"

		rstr += "<td  style=\"max-width: 300px\">" + v["lxing"] + "</td>"
		rstr += "<td  style=\"max-width: 300px\">" + v["idkey"] + "</td>"
		rstr += "<td  style=\"max-width: 300px\">" + v["llan"] + "</td>"
		rstr += "<td  style=\"max-width: 300px\">" + v["purl"] + "</td>"
		rstr += "<td  style=\"max-width: 300px\">" + v["ip"] + "</td>"
		rstr += "<td  style=\"max-width: 300px\">" + v["rq"] + "</td>"

		rstr += "</tr>"
	}
	t.Rhtml = template.HTML(rstr)
}
