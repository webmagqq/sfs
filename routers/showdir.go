package routers

import (
	"net/http"
	"sfsgo/mysql"
	"sfsgo/pubgo"
	"strconv"
	"strings"
)

func Showdir(w http.ResponseWriter, req *http.Request) {
	tid := req.URL.Query().Get("tid")
	l := req.URL.Query().Get("l")

	text := gethtml("http://"+req.Host, tid, l)
	w.Write([]byte(text))

}
func gethtml(path, id, l string) string {
	html := `
	<li class='list-group-item d-flex justify-content-between align-items-start'>
                    <div class='ms-2 me-auto'>
                        <a href='javascript:' onclick='showdir(this)'  data-tid='(tid)' data-dir='(dir),(tid)' data-jlv='(jlv)' data-title='(jingming)'>(jingming)</a> 
                    </div>
                   <a href='javascript:'><span class='badge bg-primary rounded-pill' data-dir='(dir),(tid)' data-title='(jingming)' data-jlv='(jlv)' onclick='showseldir(this)'>选择</span></a>
                    <a href='javascript:'><span class='badge bg-danger rounded-pill'  data-dir='(dir),(tid)' data-title='(jingming)' data-jlv='(jlv)' onclick='showseldir(this,1)'>排除</span></a>
                </li>
	`
	html1 := `
	<li class='list-group-item d-flex justify-content-between align-items-start'>
                    <div class='ms-2 me-auto'>
					<a href="(Host)/quanwen/(tid)">(jingming)</a>
                    </div>
                  </li>
	`
	restr, str, tid, jingming, dir, jlv, odir := "", "", "", "", "", "", ""
	itid := 0
	rd := mysql.Selects("SELECT tid,jingming,dir from jingmulu WHERE fid=?", id)
	for _, v := range rd {
		jingming = v["jingming"]
		if strings.Contains(jingming, "目录") {
			continue
		}
		jingming = pubgo.Jf(jingming, l != "")
		dir = v["dir"]
		odir = dir[:1]
		switch odir {
		case "1":
			jlv = "jingbu"
		case "2":
			jlv = "lvbu"
		case "3":
			jlv = "lunbu"
		}
		tid = v["tid"]
		itid, _ = strconv.Atoi(tid)

		if itid < 1833 {
			//restr += html.Replace("(dir)", dir).Replace("(jingming)", jingming).Replace("(tid)", tid).Replace("(jlv)", jlv)
			str = strings.Replace(html, "(dir)", dir, -1)
			str = strings.Replace(str, "(jingming)", jingming, -1)
			str = strings.Replace(str, "(tid)", tid, -1)
			str = strings.Replace(str, "(jlv)", jlv, -1)
			restr += str
		} else {
			str = strings.Replace(html1, "(jingming)", jingming, -1)
			str = strings.Replace(str, "(tid)", tid, -1)
			str = strings.Replace(str, "(Host)", path, -1)
			restr += str
		}

	}
	return restr
}
