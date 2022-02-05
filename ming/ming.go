package ming

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"sfsgo/gstr"
	"sfsgo/mysql"
	"sfsgo/pubgo"
	"strconv"
	"strings"
)

//搜索经名
type ming struct {
	req    *http.Request
	Host   string
	Path   string //Url
	Uri    string //req.RequestURI
	Up     string //繁体参数 ?l=0
	Ift    bool   //是否转繁体
	Sid    string //-搜索初始 id
	P      string //页码,定义string是为了方便截取字符串
	P1     string //页码,P+1
	K      string //数据库搜索词当Kw是繁体时，k=Kw的简体
	Kw     string //用户搜索词，可能是繁体字
	Sk     string //解析后的搜索词，提取用户搜索词前后2位，和中间一位字符的组合
	Csql   string //查询语句
	Osql   string //-查询的字词频排序语句
	Rd     []map[string]string
	Ct     int           //返回记录数
	Rehtml template.HTML //查询结果并组成html//不转义必须定义为 template.HTML 类型

}

func Newming(req *http.Request) ming {
	//静态变量返回值用指针*,非静态类则返回copy一份，取决于new
	return ming{req: req, Ift: false}
}
func (s *ming) Setto(uif *pubgo.Urlinfo) {
	s.Host = uif.Host
	s.Path = uif.Path
	s.Uri = uif.Uri
	s.Up = uif.Up
	s.Ift = uif.Ift
}

func (s *ming) SetP() string {
	P := ""
	if strings.Contains(s.Path, "-") {
		P = gstr.Do(s.Path, "-", "", true, true)
	} else {
		P = "1"
	}
	if p, _ := strconv.Atoi(s.P); p > 100 {
		P = "100" //--最大提供30页
	}
	return P
}
func (s *ming) SetP1(P string) string {
	intp, _ := strconv.Atoi(P)
	return strconv.Itoa(intp + 1)
}

func (s *ming) SetKw(Path string) string {
	Kw := ""
	if strings.Contains(Path, "-") {
		Kw = gstr.Do(Path, "/", "-", true, false)
	} else {
		Kw = gstr.Do(Path, "/", "", true, true)
	}
	return strings.TrimSpace(Kw)
}
func (s *ming) SeSid(Path string) string {
	if strings.Contains(Path, "-") {
		return gstr.Mstr(Path, "-", "-")
	} else {
		return "0"
	}
}
func (s *ming) SetK(Kw string, Ift bool) string {
	return pubgo.Fj(Kw, Ift)
}

func (s *ming) SetSk(K string) string {
	k := pubgo.Sublen(K, 35)
	reg := regexp.MustCompile(`[\p{Han}]+`) // 查找连续的汉字
	ks := reg.FindAllString(k, -1)          //,并生成数组
	ksj := strings.Join(ks, "")
	bks := []rune(ksj)
	lks := len(bks)
	if lks >= 5 { //提取用户搜索词前后2位，和中间一位字符的组合
		return string(bks[:2]) + string(bks[lks/2]) + string(bks[lks-2:])
	} else {
		return ksj
	}

}
func (s *ming) SetOsql(Sk string) string {
	sk := []rune(Sk)
	l := len(sk)
	rs := ""
	for i, v := range sk {
		rs += "'" + string(v) + "'"
		if i < l-1 {
			rs += ","
		}
	}
	return "SELECT ci from mingcipin WHERE ci in(" + rs + ")  ORDER BY ct "
}
func (s *ming) Setcsql(Osql string) string {
	if Osql == "" {
		fmt.Println(Osql, "未设置")
		return ""
	}
	rd := mysql.Selects(Osql)
	rdlen := len(rd)
	sfsql, csql, ci := "", "", ""
	for i, v := range rd {
		ci = v["ci"]
		if i < rdlen-1 {
			sfsql = "SELECT jid from ming WHERE ci like '" + ci + "%' AND jid IN (<rep>) "
			if i == 0 {
				csql = sfsql
			} else {
				csql = strings.Replace(csql, "<rep>", sfsql, -1) //csql.Replace("<rep>", sfsql);
			}
		} else { //--最后一条
			sfsql = " SELECT jid from ming WHERE ci like '" + ci + "%' "
		}
	}
	if rdlen > 1 {
		csql = strings.Replace(csql, "<rep>", sfsql, -1) //csql.Replace("<rep>", sfsql)
	} else {
		csql = sfsql
	}
	csql = strings.Replace(csql, "WHERE", "WHERE jid>"+s.Sid+" and ", -1) //csql.Replace("WHERE", "WHERE jid>" + sid + " and ")
	csql = strings.Replace(csql, "SELECT", "select DISTINCT", -1)
	csql += " LIMIT  21"
	return csql
	//fmt.Println(csql)
}
func (s *ming) getids() string {
	if s.Csql == "" {
		fmt.Println(s.Csql, "未设置")
		return ""
	}
	rd := mysql.Selects(s.Csql)
	l := len(rd)
	if l == 0 {
		s.Ct = 0
		return ""
	}
	s.Ct = l
	rs := ""
	for i, v := range rd {
		rs += v["jid"]
		if i < l-2 {
			rs += ","
		}
	}
	return rs
}
func setrd(Csql string) []map[string]string {
	return mysql.Selects(Csql)
}
func (s *ming) SetRehtml(Host, Csql, Up string, Ift bool) template.HTML {
	if Csql == "" {
		fmt.Println(Csql, "未设置")
		return ""
	}
	rstr, astr := "", ""
	ids := s.getids()
	if ids == "" {
		return template.HTML(pubgo.Jf("没有找到！", Ift))
	}
	csql := "SELECT tid,jingming FROM jingmulu WHERE tid IN (" + ids + ")"
	rd := setrd(csql)
	s.Rd = rd
	for _, v := range rd {
		astr = v["jingming"] + "<a style=\"float: right\" href=\"" + Host + "/quanwen/" + v["tid"]
		astr += Up + "\"   title='" + pubgo.Jf(v["jingming"]+"全文,原文完整版诵读", Ift) + "'  >(查看<i class=\"bi-caret-right-fill\"></i>)</a>"
		rstr += "<div class=\"card-footer\">" + astr + "</div>"

	}

	return template.HTML(rstr)
}
func (s *ming) SetSid(rd []map[string]string) string {
	rdlen := len(rd)
	if rdlen > 0 {
		return rd[rdlen-1]["tid"]
	}
	return "0"
}
