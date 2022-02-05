package routers

import (
	"html/template"
	"net/http"
	"sfsgo/gstr"
	"sfsgo/pubgo"
	"sfsgo/wen"
)

func Wen(w http.ResponseWriter, req *http.Request) {
	//数据组织
	uif := pubgo.Newurlinfo(req)
	rd := wen.Newwen(req)
	rd.Setto(&uif)
	rd.SetWid(rd.Path)
	rd.Bt = rd.SetBt(rd.Path)
	rd.Jtext = rd.SetJtext(rd.Wid)
	rd.Jname = rd.SetJname()

	//统计
	pubgo.Tj.Brows(gstr.Mstr(req.URL.Path, "/", "/"))
	//--组织模板数据
	TemplatesFiles := []string{
		"wen/wen.html",
		"tplpub/static.html",
		"tplpub/header.html",
		"tplpub/footer.html", // 多加的文件
	}
	funcMap := template.FuncMap{ //--需要注册的函数
		"jf": pubgo.Jf, //简繁转换
	}
	t, _ := template.New("wen.html").Funcs(funcMap).ParseFiles(TemplatesFiles...)
	//--  wen.html必须是TemplatesFiles第一个文件名
	t.ExecuteTemplate(w, "wen.html", rd)

}
