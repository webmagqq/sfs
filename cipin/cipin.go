package cipin

import (
	"html/template"
	"net/http"
	"sfsgo/gstr"
	"sfsgo/mysql"
	"sfsgo/pubgo"
	"strconv"
)

type cipin struct {
	req  *http.Request
	Host string
	Path string
	Uri  string
	Up   string //繁体参数 ?l=0
	Ift  bool   //是否转繁体

	P      string
	Kwhtml template.HTML
	PgHtml template.HTML
}

func Newcipin(req *http.Request) cipin {
	//静态变量返回值用指针*,非静态类则返回copy一份，取决于new
	return cipin{req: req}
}

func (c *cipin) Setto(uif *pubgo.Urlinfo) {
	c.Host = uif.Host
	c.Path = uif.Path
	c.Uri = uif.Uri
	c.Up = uif.Up
	c.Ift = uif.Ift
}
func (c *cipin) SetP() {
	c.P = gstr.Do(c.Path, "/", "", true, true)
}
func (c *cipin) SetKwhtml() {
	rstr, ci := "", ""
	intp, _ := strconv.Atoi(c.P)
	bl := (intp - 1) * 100
	rd := mysql.Selects("SELECT tid,ci,cishu from jingcipin WHERE tid>=? LIMIT 300", bl)
	for _, v := range rd {
		// "<button type=\"button\" class=\"btn btn-light\"><a  href=\"" + url + "/jingbu/" + fenci + lp +
		// "\"  title='什么是"+fenci+","+fenci+"在经中是什么意思'>" + fenci + "</a> <span class=\"badge bg-secondary\">" + pinci + "</span></button> ";
		ci = v["ci"]
		ci = pubgo.Jf(ci, c.Ift)
		rstr += "<button type=\"button\" class=\"btn btn-light\"><a  href=\"" + c.Host + "/jingbu/" + ci + c.Up
		rstr += "\"  title='什么是" + ci + "," + ci + "是什么意思'>" + ci + "</a>"
		rstr += "<span class=\"badge bg-secondary\">" + v["cishu"] + "</span></button> "
	}
	c.Kwhtml = template.HTML(rstr)
}
func (c *cipin) SetPgHtml() {
	b, _ := strconv.Atoi(c.P)
	if b == 1 {
		b = 2
	}
	e := b + 15
	if e > 2467 {
		e = 2467
	}
	pages := ""
	for b <= e {
		pages += "<a  href=\"" + c.Host + "/cipin/" + strconv.Itoa(b) + c.Up + "\"><span class=\"btn btn-outline-secondary\">" + strconv.Itoa(b) + "</span></a>"
		b++
	}
	c.PgHtml = template.HTML(pages)
}
