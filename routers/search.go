package routers

import (
	"fmt"
	"html/template"
	"net/http"
	"sfsgo/gstr"
	"sfsgo/pubgo"
	"sfsgo/search"
)

func Search(w http.ResponseWriter, req *http.Request) {
	ts := pubgo.Newts()
	rd := search.Newsearch(req)
	uif := pubgo.Newurlinfo(req)
	rd.Setto(&uif)

	rd.Up = rd.SetUp()
	rd.Ift = rd.SetIft()
	rd.Jlv = rd.SetJlv()
	rd.Kw = rd.SetKw(rd.Path)
	if rd.Kw == "" {
		http.Redirect(w, req, "/", http.StatusFound)
		return
	}
	rd.K = rd.SetK(rd.Kw, rd.Ift)
	rd.Dir = rd.SetDir()
	rd.Lwwh = rd.SetLwwh(rd.Dir)
	rd.Jp = rd.SetJp()
	rd.Ks = rd.SetKs(rd.K)
	rd.P = rd.SetP(rd.Path, rd.Ks)
	rd.P1 = rd.SetP1(rd.P)
	rd.Tjd = rd.SetTjd()
	rd.JoinTjd = rd.SetJoinTjd(rd.Tjd)
	rd.Jz = rd.SetJz()
	rd.Mkey = rd.SetMkey(rd.Jlv, rd.JoinTjd, rd.P, rd.Dir, rd.Jz)
	rd.Setsearchdata(rd.Mkey)
	rd.Rehtml = rd.SetRehtml(rd.Ks, rd.Ift, rd.K, rd.P, rd.Host, rd.Up)
	rd.Ashtml = rd.SetAshtml(rd.Ks[0], rd.Jlv, rd.Up, rd.Ift)
	rd.Cidian = rd.SetCidian(rd.K, rd.Kw, rd.Ift)
	//fmt.Println(req.Referer())
	//统计
	pubgo.Tj.Brows(gstr.Mstr(req.URL.Path, "/", "/"))
	//--组织模板数据
	TemplatesFiles := []string{
		"search/search.html",
		"tplpub/static.html",
		"tplpub/header.html",
		"tplpub/footer.html", // 多加的文件
	}
	//t, _ := template.ParseFiles(TemplatesFiles...)

	funcMap := template.FuncMap{ //--需要注册的函数
		"jf": pubgo.Jf, //简繁转换
	}
	t, _ := template.New("search.html").Funcs(funcMap).ParseFiles(TemplatesFiles...)
	rd.Pgt = ts.Gts()
	err := t.ExecuteTemplate(w, "search.html", rd)
	if err != nil {
		fmt.Println(req.URL.Path, err)
	}
}
