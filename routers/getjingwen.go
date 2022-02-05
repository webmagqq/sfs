package routers

import (
	"net/http"
	"sfsgo/mysql"
	"sfsgo/pubgo"
	"strings"
)

func Getjingwen(w http.ResponseWriter, req *http.Request) {
	tid := req.URL.Query().Get("tid")
	l := req.URL.Query().Get("l")

	text := gnr(tid, l)
	w.Write([]byte(text))

}
func gnr(tid, l string) string {
	rd := mysql.Selects("SELECT nr from dzj WHERE tid>? and tid<(SELECT tid from dzj WHERE tid > ? and (tn = 'a' or tn = 'b-ne' or nr like '下一部%')  limit 1)", tid, tid)
	str, nr := "", ""
	for _, v := range rd {
		//nr += msDataSet.Tables[0].Rows[i]["nr"].ToString().Replace("br", "<br>").Replace("p0", "<br>").Replace("p1", "<br>").Replace("zh1", " ");
		nr = pubgo.Jf(v["nr"], l != "")
		nr = strings.Replace(nr, "br", "<br>", -1)
		nr = strings.Replace(nr, "p0", "<br>", -1)
		nr = strings.Replace(nr, "p1", "<br>", -1)
		//nr = strings.Replace(nr, "p1", "<br>", -1)
		str += nr
	}
	return str
}
