package routers

import (
	"html/template"
	"net/http"
	"sfsgo/gstr"
	"sfsgo/mapp"
	"sfsgo/pubgo"
)

func Mapp(w http.ResponseWriter, req *http.Request) {
	//数据组织
	uif := pubgo.Newurlinfo(req)
	rd := mapp.Newmapp(req)
	rd.Setto(&uif)

	//统计
	pubgo.Tj.Brows(gstr.Mstr(req.URL.Path, "/", "/"))
	//--组织模板数据
	TemplatesFiles := []string{
		"mapp/mapp.html",
		"tplpub/static.html",
		"tplpub/header.html",
		"tplpub/footer.html", // 多加的文件
	}
	funcMap := template.FuncMap{ //--需要注册的函数
		"jf": pubgo.Jf, //简繁转换
	}
	t, _ := template.New("mapp.html").Funcs(funcMap).ParseFiles(TemplatesFiles...)
	//--  quanwen.html必须是TemplatesFiles第一个文件名
	t.ExecuteTemplate(w, "mapp.html", rd)

}
