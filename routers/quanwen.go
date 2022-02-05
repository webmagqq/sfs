package routers

import (
	"html/template"
	"net/http"
	"sfsgo/gstr"
	"sfsgo/pubgo"
	"sfsgo/quanwen"
	"strings"
)

func Quanwen(w http.ResponseWriter, req *http.Request) {
	//数据组织
	uif := pubgo.Newurlinfo(req)
	rd := quanwen.Newquanwen(req)
	rd.Setto(&uif)

	rd.Wid = rd.SetWid(rd.Host)
	rd.Wid = strings.Replace(rd.Wid, "jing-", "", -1) //wid.Replace("jing-", "").Replace("lv-", "").Replace("lun-", "");//兼容前面版本
	rd.Wid = strings.Replace(rd.Wid, "lv-", "", -1)
	rd.Wid = strings.Replace(rd.Wid, "lun-", "", -1)
	rd.Jinginfo = rd.SetJinginfo(rd.Wid)
	rd.SetJingming()
	rd.Jid = rd.SetJid(rd.Jinginfo[0], rd.Wid)
	if rd.Jinginfo[0]["cd"] == "1" {
		rd.Jtext = rd.SetJtext(rd.Host+"/getjingwen/?tid="+rd.Jinginfo[0]["did"]+rd.Up, rd.Ift) //qjubliang.url + "/getjingwen.aspx?tid="+ yid+ fajian.UrlPara
	}
	rd.Rhtml = rd.SetRhtml(rd.Host, rd.Wid, string(rd.Jtext), rd.Jid, rd.Up, rd.Ift)
	rd.Hbhtml = rd.SetHbhtml(rd.Host, rd.Jid, rd.Up, rd.Ift)
	rd.Tophtml = rd.SetTophtml(rd.Jid, rd.Ift)

	//统计
	pubgo.Tj.Brows(gstr.Mstr(req.URL.Path, "/", "/"))

	//--组织模板数据
	TemplatesFiles := []string{
		"quanwen/quanwen.html",
		"tplpub/static.html",
		"tplpub/header.html",
		"tplpub/footer.html", // 多加的文件
	}
	funcMap := template.FuncMap{ //--需要注册的函数
		"jf": pubgo.Jf, //简繁转换
	}
	t, _ := template.New("quanwen.html").Funcs(funcMap).ParseFiles(TemplatesFiles...)
	//--  quanwen.html必须是TemplatesFiles第一个文件名
	t.ExecuteTemplate(w, "quanwen.html", rd)

}
