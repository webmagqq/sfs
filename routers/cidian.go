package routers

import (
	"fmt"
	"net/http"
	"sfsgo/mysql"
	"sfsgo/pubgo"
	"strings"
)

func Cidian(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query().Get("q")
	a := req.URL.Query().Get("a")
	dir := req.URL.Query().Get("dir")
	l := req.URL.Query().Get("l")
	q = pubgo.Sublen(q, 7)
	q = pubgo.Fj(q, l != "")
	sql := ""
	if a != "foxuecidian" {
		ar0 := []rune(q)[0]
		num := int(ar0) % 21
		if a != "ming" {
			a = fmt.Sprintf("%sidx%d", strings.Replace(a, "bu", "", -1), num)
		} else {
			a = "mingidx"
		}
		lmwh := ""
		if dir != "" {
			//lmwh = (dir.IndexOf("-") == -1) ? " and dir like '" + dir + "%' " : " and dir not like '" + dir + "%' ";
			if !strings.Contains(dir, "-") { //--排除符号
				lmwh = " and dir like '" + dir + "%' "
			} else {
				lmwh = " and dir not like '" + dir + "%' "
			}
		}
		sql = "SELECT DISTINCT ci from " + a + " where ci LIKE '" + q + "%' " + lmwh + " LIMIT 9"
	} else {
		sql = "SELECT ci  from cidian WHERE ci like '" + q + "%' limit 9"
	}
	e0 := strings.Index(sql, "*")
	e1 := strings.Index(sql, "/")
	e2 := strings.Index(sql, "(")
	e3 := strings.Index(sql, ")")
	e4 := strings.Index(sql, "-")
	if e0+e1+e2+e3+e4 != -5 { //防止sql注入
		w.Write([]byte(""))
		return
	}
	rd := mysql.Selects(sql)
	rstr := ""
	for _, v := range rd {
		rstr += "<option>" + v["ci"] + "</option>"
	}
	w.Write([]byte(rstr))
}
