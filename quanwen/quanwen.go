package quanwen

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"sfsgo/gstr"
	"sfsgo/mysql"
	"sfsgo/pubgo"
	"strconv"
)

type quanwen struct {
	req      *http.Request
	Host     string
	Path     string
	Uri      string
	Up       string //繁体参数 ?l=0
	Ift      bool   //是否转繁体
	Wid      string
	Jinginfo []map[string]string //经的目录信息
	Jid      string              //--没did则是顶级经目录，否则是子目录。
	//Jcid     string //--顶级经目录没wid，所以为-1.否则=wid
	Jingming string
	Jtext    template.HTML
	Rhtml    template.HTML //--经文目录信息
	Hbhtml   template.HTML //--获取前后7部相关经典信息
	Tophtml  template.HTML //--顶级经目录信息
}

func Newquanwen(req *http.Request) quanwen {
	//静态变量返回值用指针*,非静态类则返回copy一份，取决于new
	return quanwen{req: req}
}

func (q *quanwen) Setto(uif *pubgo.Urlinfo) {
	q.Host = uif.Host
	q.Path = uif.Path
	q.Uri = uif.Uri
	q.Up = uif.Up
	q.Ift = uif.Ift
}
func (q *quanwen) SetWid(Path string) string {
	return gstr.Do(q.Path, "/", "", true, true)
}
func (q *quanwen) SetJinginfo(Wid string) []map[string]string {
	return mysql.Selects("SELECT jingming,fid,did,cd  FROM jingmulu WHERE tid=?", Wid)
}
func (q *quanwen) SetJid(Jinginfo0 map[string]string, Wid string) string {
	jinfo := Jinginfo0      //q.Jinginfo[0]
	if jinfo["did"] == "" { //--没did则是顶级经目录即是wid本事，否则是子目录。
		return Wid //--顶级目录只打开目录信息，子目录需要多加打开内容
	} else { //--顶级目录只打开目录信息，子目录需要多加打开内容
		return jinfo["fid"] //子目录时，同时cd=1，为有内容。cd=0时有可能是卷，但无内容
	}
}
func (q *quanwen) SetJtext(url string, Ift bool) template.HTML { //--打开某品、卷的内容
	res, err := http.Get(url)
	if err != nil {
		return template.HTML(string(err.Error()))
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return template.HTML(string(err.Error()))

	}
	res.Body.Close()
	bodys := (pubgo.Jf(string(body), Ift))
	return template.HTML(bodys)
}
func (q *quanwen) SetRhtml(Host, Wid, Jtext, Jid, Up string, Ift bool) template.HTML {
	jingming := ""
	rd := mysql.Selects("SELECT tid,jingming,cd,did from jingmulu WHERE fid=? ORDER BY did ", Jid)
	rs, hzk := "", ""
	for _, v := range rd {
		jingming = v["jingming"]
		jingming = pubgo.Jf(jingming, Ift)
		if v["cd"] == "1" {
			if v["tid"] != Wid /*q.Jcid*/ {
				hzk = "<i class=\"bi-chevron-down\" id='did" + v["did"] + "'></i>"
				rs += "<div class=\"card-footer\" ><a href='" + Host + "/quanwen/" + v["tid"] + Up + "'>" + jingming + "</a>  " + hzk + "</div>"
				rs += "<div class=\"card-footer\"  style=\"display: none; \"></div>"
			} else {
				hzk = "<a href='javascript:'  onclick=\"playv(document.getElementById('ld" + v["tid"] + "').innerText)\" ><i class=\"bi-caret-right-fill\">诵读</i></a>"
				rs += "<div class=\"card-footer\" id='jing" + v["tid"] + "'><b>" + jingming + "</b>   " + hzk + "</div>"
				rs += "<div class=\"card-footer\" id='ld" + v["tid"] + "'>" + string(Jtext) + "</div>"
			}
		} else {
			rs += "<div class=\"card-footer\" id='jing" + v["tid"] + "'><b>" + jingming + "</b></div>"
		}
	}
	return template.HTML(rs)
}
func (q *quanwen) SetHbhtml(Host, Jid, Up string, Ift bool) template.HTML {
	ijid, _ := strconv.Atoi(Jid)
	bjid := ijid - 2
	ejid := ijid + 7
	rstr, jingming, ico := "", "", ""
	var tid int
	rd := mysql.Selects("SELECT tid,jingming FROM jingmulu WHERE tid>=? and tid<=?", bjid, ejid)
	for _, v := range rd {
		if Jid == v["tid"] {
			ico = "<i class='bi-arrow-right'></i>"
			jingming = "<b>" + v["jingming"] + "</b>"
		} else {
			jingming = v["jingming"]
			tid, _ = strconv.Atoi(v["tid"])
			if tid < ijid {
				ico = "<i class='bi-arrow-up-short'></i>"
			} else {
				ico = "<i class='bi-arrow-down-short'></i>"
			}
		}
		rstr += "<div class=\"card-footer\">" + ico + "<a  href=\"" + Host + "/quanwen/" + v["tid"] + Up
		rstr += "\"   title='" + pubgo.Jf(jingming+"全文,原文完整版诵读", Ift) + "'  >" + pubgo.Jf(jingming, Ift) + "</a></div>"
	}
	return template.HTML(rstr)
}

func (q *quanwen) SetTophtml(Jid string, Ift bool) template.HTML {
	rstr, nr := "", ""
	rd := mysql.Selects("SELECT nr FROM dzj WHERE jid=? limit 5 ", Jid)
	for _, v := range rd {
		nr = v["nr"]
		if nr == "ml" {
			break
		}
		rstr += pubgo.Jf(nr, Ift) + "<br />"
	}
	return template.HTML(rstr)
}
func (q *quanwen) SetJingming() {
	q.Jingming = q.Jinginfo[0]["jingming"]
}
