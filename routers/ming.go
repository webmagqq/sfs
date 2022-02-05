package routers

import (
	"fmt"
	"html/template"
	"net/http"
	"sfsgo/gstr"
	"sfsgo/ming"
	"sfsgo/pubgo"
)

func Ming(w http.ResponseWriter, req *http.Request) {
	uif := pubgo.Newurlinfo(req)
	rd := ming.Newming(req)
	rd.Setto(&uif)
	rd.P = rd.SetP()
	rd.P1 = rd.SetP1(rd.P)
	rd.Kw = rd.SetKw(rd.Path)
	rd.K = rd.SetK(rd.Kw, rd.Ift)
	rd.Sid = rd.SeSid(rd.Path)
	rd.Sk = rd.SetSk(rd.K)
	rd.Osql = rd.SetOsql(rd.Sk)
	rd.Csql = rd.Setcsql(rd.Osql)
	rd.Rehtml = rd.SetRehtml(rd.Host, rd.Csql, rd.Up, rd.Ift)
	rd.Sid = rd.SetSid(rd.Rd)

	//fmt.Println(rd)
	//统计
	pubgo.Tj.Brows(gstr.Mstr(req.URL.Path, "/", "/"))
	//--组织模板数据
	TemplatesFiles := []string{
		"ming/ming.html",
		"tplpub/static.html",
		"tplpub/header.html",
		"tplpub/footer.html", // 多加的文件
	}
	//t, _ := template.ParseFiles(TemplatesFiles...)

	funcMap := template.FuncMap{ //--需要注册的函数
		"jf": pubgo.Jf, //简繁转换
	}
	t, _ := template.New("ming.html").Funcs(funcMap).ParseFiles(TemplatesFiles...)

	err := t.ExecuteTemplate(w, "ming.html", rd)
	if err != nil {
		fmt.Println(req.URL.Path, err)
	}
}
