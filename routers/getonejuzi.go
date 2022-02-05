package routers

import (
	"net/http"
	"regexp"
	"sfsgo/mysql"
	"sfsgo/pubgo"
	"strconv"
	"strings"
)

func Getonejuzi(w http.ResponseWriter, req *http.Request) {
	tid := req.URL.Query().Get("tid")
	jid := req.URL.Query().Get("jid")
	kw := req.URL.Query().Get("kw")
	updown := req.URL.Query().Get("updown")
	l := req.URL.Query().Get("l")
	//int cid = (updown == "0") ? int.Parse(tid) - 1 : int.Parse(tid) + 1;
	cid, _ := strconv.Atoi(tid)
	if updown == "0" {
		cid = cid - 1
	} else {
		cid = cid + 1
	}
	rd := mysql.Selects("SELECT nr,tn from dzj WHERE tid=? and jid=?", cid, jid)
	nr := rd[0]["nr"]
	tn := rd[0]["tn"]
	if tn == "a" {
		nr = ""
	}
	if tn == "p0" {
		nr = ""
	}
	nr = strings.Replace(nr, "br", "", -1)
	nr = strings.Replace(nr, "zh1", " ", -1)
	if nr != "" {
		nr = pubgo.Jf(nr, l != "")
		nr = jcu(kw, nr, l)
	}
	w.Write([]byte(nr))

}
func jcu(kw, text, l string) string { //加粗
	reg := regexp.MustCompile(`[\p{Han}]+`) // 查找连续的汉字
	ks := reg.FindAllString(kw, -1)         //,并生成数组
	if len(ks) > 1 {
		ks = pubgo.RemoveRepeatElement(ks) //去重
	}
	str := text
	for _, v := range ks {
		str = strings.Replace(str, v, "<b>"+pubgo.Jf(v, l != "")+"</b>", -1)
	}
	return str
}
