package routers

import (
	"html/template"
	"net/http"
	"sfsgo/pubgo"
	"sfsgo/tongji"
)

func Tongji(w http.ResponseWriter, req *http.Request) {
	//数据组织
	uif := pubgo.Newurlinfo(req)
	rd := tongji.Newtongji(req)
	rd.Setto(&uif)
	rd.SetP()
	rd.SetRhtml()
	rd.Trtb()

	//--组织模板数据
	TemplatesFiles := []string{
		"tongji/tongji.html",
		"tplpub/static.html",
		"tplpub/header.html",
		"tplpub/footer.html", // 多加的文件
	}
	funcMap := template.FuncMap{ //--需要注册的函数
		"jf": pubgo.Jf, //简繁转换
	}
	t, _ := template.New("tongji.html").Funcs(funcMap).ParseFiles(TemplatesFiles...)
	//--  quanwen.html必须是TemplatesFiles第一个文件名
	t.ExecuteTemplate(w, "tongji.html", rd)

}
