package search

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"sfsgo/cache"
	"sfsgo/gstr"
	"sfsgo/mysql"
	"sfsgo/pubgo"
	"strconv"
	"strings"
	"time"
)

var (
	CacheSearchData, CacheNrData *cache.CacheData
)

type search struct {
	req  *http.Request
	Host string
	Path string //Url
	Uri  string //req.RequestURI
	Up   string //繁体参数 ?l=0
	Ift  bool   //是否转繁体
	Jlv  string
	P    string   //页码,定义string是为了方便截取字符串
	P1   string   //页码,P+1
	K    string   //数据库搜索词当Kw是繁体时，k=Kw的简体
	Kw   string   //用户搜索词，可能是繁体字
	Dir  string   //-用户选择查询
	Lwwh string   //选择查询的SQL语句
	Sql  string   //搜索的查询语句
	Jz   string   //精准搜索模式，默认是全搜索，只有组合查询时起作用
	Jp   string   //参数连接符，?、&
	Ks   []string //搜索词分解成的数组
	//csql   string              //组合查询各词的词频的sql语句
	//js      jonisql
	Tjd     []map[string]string         //组合查询各词的词频
	JoinTjd string                      //Tjd转换为字符串
	sdata   []map[string]string         //搜索查询数据结果
	Mkey    string                      //根据搜索词和参数组成 map 的key
	scaches map[string]*cache.CacheData //缓存
	Ct      int                         //返回记录数

	Rehtml template.HTML //查询结果并组成html//不转义必须定义为 template.HTML 类型
	Ashtml template.HTML //相关搜索词html//不转义必须定义为 template.HTML 类型
	Cidian template.HTML //相关搜索词html//不转义必须定义为 template.HTML 类型
	St     time.Duration //搜索用时
	Sst    time.Duration //数据结集用时
	Pgt    time.Duration //页面用时
	pgn    int           //--每页返回记录数

}

func Newsearch(req *http.Request) search {
	//静态变量返回值用指针*,非静态类则返回copy一份，取决于new
	return search{
		req: req,
		Ift: false,
		pgn: 21,
		scaches: map[string]*cache.CacheData{
			"sch": CacheSearchData,
			"nr":  CacheNrData,
		},
	}
}

func (s *search) Setto(uif *pubgo.Urlinfo) {
	s.Host = uif.Host
	s.Path = uif.Path
	s.Uri = uif.Uri
	s.Up = uif.Up
	s.Ift = uif.Ift
}

func (s *search) SetP(Path string, Ks []string) string {
	P := ""
	if strings.Contains(Path, "-") {
		P = gstr.RStr(Path, "-")
	} else {
		P = "1"
	}
	p, _ := strconv.Atoi(P)
	slen := len(Ks)
	if slen == 1 && p > 100 {
		P = "100" //--单词查询最多提供100页
	}
	if slen > 1 && p > 20 {
		P = "20" //--组合查询最多提供20页
	}
	return P
}
func (s *search) SetP1(P string) string { //-下一页
	intp, _ := strconv.Atoi(P)
	return strconv.Itoa(intp + 1)
}
func (s *search) SetUp() string {
	up := ""
	if s.req.URL.Query().Get("l") != "" {
		up = "?l=0"
	}
	return up
}
func (s *search) SetJz() string {
	return s.req.URL.Query().Get("jz")
}
func (s *search) SetJp() string {
	if strings.Contains(s.Uri, "?") {
		return "&"
	} else {
		return "?"
	}
}
func (s *search) SetIft() bool {
	return s.Up != ""
}

func (s *search) SetJlv() string {
	return gstr.Mstr(s.req.URL.Path, "/", "/")
}
func (s *search) SetHost() {
	s.Host = "http://" + s.req.Host
}

func (s *search) SetKw(Path string) string {
	Kw := ""
	if strings.Contains(Path, "-") {
		Kw = gstr.Do(Path, "/", "-", true, false)
		//s.P = GetRStr(s.Url, "-")
	} else {
		Kw = gstr.Do(Path, "/", "", true, true)
	}
	return strings.TrimSpace(Kw)
}

func (s *search) SetK(Kw string, Ift bool) string {
	return pubgo.Fj(Kw, Ift)
}
func (s *search) SetDir() string {
	return s.req.URL.Query().Get("dir")
}
func (s *search) SetLwwh(Dir string) string {
	Lwwh := ""
	if Dir != "" {
		if strings.Contains(s.Dir, "-") {
			Lwwh = " and Dir not like '" + strings.Replace(Dir, "-", "", -1) + "%' "
		} else {
			Lwwh = " and Dir like '" + Dir + "%' "
		}
	}
	return Lwwh
}
func (s *search) SetKs(K string) []string {
	if K == "" {
		fmt.Println("错误：搜索词未赋值!")
		return nil
	}
	k := pubgo.Sublen(K, 35)
	reg := regexp.MustCompile(`[\p{Han}]+`) // 查找连续的汉字
	ks := reg.FindAllString(k, -1)          //,并生成数组
	ks = pubgo.RemoveRepeatElement(ks)      //去重
	return ks
}
func (s *search) setcsql(Ks []string, Jlv string) []map[string]string {
	sql := ""
	for _, v := range Ks {
		//"(SELECT '"+v+"' ci,IFNULL((SELECT cishu FROM "+j.jlv+"cipin WHERE ci='"+v+"' ),0) cs)", " UNION ")
		sql += strings.Join([]string{"(SELECT '", v, "' ci,IFNULL((SELECT cishu FROM ", strings.Replace(Jlv, "bu", "", -1), "cipin WHERE ci='", v, "' ),0) cs)  UNION  "}, "")
	}
	//rd := Selects(sql)
	//fmt.Println(rd)
	sql = gstr.Do(sql, "", "UNION", false, true)
	sql = "SELECT * from (" + sql + ") abc ORDER BY 2 "
	return mysql.Selects(sql)
}
func (s *search) SetTjd() []map[string]string {
	return s.setcsql(s.Ks, s.Jlv)
}
func (s *search) SetSql(Ks []string, Jlv, Lwwh string, Tjd []map[string]string) string {
	ks := Ks
	intp, _ := strconv.Atoi(s.P) //字符串转int
	pgn := s.pgn                 //每页显示记录数
	pg := pgn * (intp - 1)
	limitsql := " limit " + strconv.Itoa(pg) + "," + strconv.Itoa(pgn)
	var gs searchresult
	if len(ks) == 1 { //单词查询//--动态拼接sql都不能用预处理方法
		gs = Newonesql(ks[0], strings.Replace(Jlv, "bu", "", -1), Lwwh)
	} else { //组合查询
		gs = Newjonisql(s.Jz, ks[0], strings.Replace(Jlv, "bu", "", -1), Lwwh, Tjd)
	}
	sql := gs.getsql()
	sql += limitsql
	return sql
}
func (s *search) SetJoinTjd(Tjd []map[string]string) string {
	rs := ""
	for _, v := range Tjd {
		rs += v["ci"] + " "
	}
	return strings.TrimSpace(rs)
}
func (s *search) SetMkey(Jlv, JoinTjd, P, Dir, Jz string) string {
	key := Jlv + " " + JoinTjd + " " + P
	if Dir != "" {
		key += " " + Dir
	}
	if Jz != "" {
		key += " " + Jz
	}
	return key
}
func (s *search) getcache(cname, Mkey string) interface{} {
	return s.scaches[cname].Get(Mkey)
}
func (s *search) addcache(cname, Mkey string, data interface{}) {
	if s.req.Referer() != "" { //屏蔽爬虫
		s.scaches[cname].Add(Mkey, data)
	}
}
func (s *search) Setsearchdata(Mkey string) {
	ts := pubgo.Newts()
	inrd := s.getcache("sch", Mkey)
	if inrd != nil {
		s.sdata = inrd.([]map[string]string)
	} else {
		Sql := s.SetSql(s.Ks, s.Jlv, s.Lwwh, s.Tjd)
		s.sdata = mysql.Selects(Sql)
		s.addcache("sch", Mkey, s.sdata)
	}
	s.St = ts.Gts()
	//fmt.Println(s.St)
}

func (s *search) SetRehtml(Ks []string, Ift bool, K, P, Host, Up string) template.HTML {
	ts := pubgo.Newts()
	pgn := s.pgn
	ks := Ks
	rdata := s.sdata
	s.Ct = len(rdata)
	//fmt.Println(rdata)
	if s.Ct == 0 {
		return template.HTML(pubgo.Fj("没有找到相关经文。", Ift))
		//fmt.Println("没有找到数据", s.sql)
	}
	tids := ""
	for _, v := range rdata {
		//jpids += v["jpid"] + ","
		tids += v["tid"] + ","
	}
	//jpids = jpids[:len(jpids)-1]
	tids = tids[:len(tids)-1]
	//sid = tidDataSet.Tables[0].Rows[tct - 1]["tid"].ToString(); //--记录最好最大的id值
	//muludata := mysql.Selects("select tid,jingming from jingmulu WHERE tid in (" + jpids + ")")
	cerd := s.getcache("nr", s.Mkey) //CacheNrData.Get(s.mkey)
	var dzjdata []map[string]string
	if cerd == nil {
		dzjdata = mysql.Selects("select a.tid,a.jid,a.nr,b.jingming from dzj a,jingmulu b WHERE a.tid in (" + tids + ") AND a.jid=b.tid")
		s.addcache("nr", s.Mkey, dzjdata)
	} else {
		dzjdata = cerd.([]map[string]string)
	}

	//fmt.Println(muludata)
	//fmt.Println(dzjdata)
	html, tid, jid, nr, jingming, bt := "", "", "", "", "", ""
	tjid, tjingming := "", ""
	//var jml map[string]string
	uptid, downtid := 0, 0
	for i, v := range dzjdata {
		tid = v["tid"]
		uptid, _ = strconv.Atoi(tid)
		downtid = uptid
		jid = v["jid"]
		nr = v["nr"]
		jingming = v["jingming"]
		//jml = gjm(tid, rdata, muludata)
		tjid, tjingming = getjm(tid, rdata)
		bt = gbt(ks[0], nr)
		nr = strings.Replace(nr, "br", "<br>", -1)
		nr = strings.Replace(nr, "p0", "", -1)
		nr = strings.Replace(nr, "p1", "", -1)
		nr = allb(ks, nr) //--将所有关键词加粗
		html += "<div class=\"card-footer\" data-kw='" + K + "' data-jid='" + jid + "'>"
		intp, _ := strconv.Atoi(P)
		intp = pgn*(intp-1) + i + 1
		html += "<span class='badge bg-primary'>" + strconv.Itoa(intp) + "</span> <a href='javascript:' class='fs-5'><i class='bi-arrow-up-circle' data-id='" + strconv.Itoa(uptid) + "'></i></a><span class='nrs'>" + pubgo.Jf(nr, Ift) + "</span>... <a href='javascript:' class=\"fs-5\"><i class=\"bi-arrow-down-circle\"  data-id='" + strconv.Itoa(downtid) + "'></i></a></div>"
		html += "<div  class='card-header text-end'> "
		html += "<a class='text-secondary' href='" + Host + "/quanwen/" + jid + Up + "' title='" + pubgo.Jf(jingming+"全文,原文完整版诵读", Ift) + "' >" + pubgo.Jf(jingming, Ift) + "</a> "
		html += "<i class='bi-forward-fill'></i><a class='text-secondary' href='" + Host + "/quanwen/" + tjid + Up + "' title='" + pubgo.Jf(tjingming+"全文,原文完整版诵读", Ift) + "' >" + pubgo.Jf(tjingming, Ift) + "</a> "
		html += "<i class='bi-plus'></i><a class='text-secondary' href='" + Host + "/" + strings.Replace(s.Jlv, "bu", "", -1) + "/" + tid + "_" + bt + Up + "' title='" + pubgo.Jf(bt, Ift) + "' >" + pubgo.Jf("摘录", Ift) + "</a></div>" //" + Jf(bt, s.Ift) + "
	}
	s.Sst = ts.Gts()
	return template.HTML(html)

}

func getjm(tid string, a []map[string]string) (string, string) {
	jpid, jname := "", ""
	for _, v := range a {
		if tid == v["tid"] {
			jpid = v["jpid"]
			break
		}
	}
	jname = Mlmap[jpid]
	return jpid, jname
}

func gbt(k, n string) string {
	reg := regexp.MustCompile(`[\p{Han}]+`) // 查找连续的汉字
	ks := reg.FindAllString(n, -1)          //,并生成数组
	for _, v := range ks {
		if strings.Contains(v, k) {
			return v
		}
	}
	return ""
}

func allb(ks []string, text string) string {
	str := text
	for _, v := range ks {
		str = strings.Replace(str, v, "<b>"+v+"</b>", -1)
	}
	return str
}
func (s *search) SetAshtml(Ks0, Jlv, Up string, Ift bool) template.HTML {
	ks := Ks0 //s.Ks[0]
	bks := []rune(ks)
	kslen := len(bks)
	if kslen > 2 {
		ks = string(bks[:kslen/2])
	} else {
		ks = string(bks)
	}

	sql := "select ci from " + strings.Replace(Jlv, "bu", "", -1) + "cipin where ci like '" + ks + "%' limit 108"
	//fmt.Println(sql)
	rd := mysql.Selects(sql)
	ci, str := "", ""
	for _, v := range rd {
		ci = v["ci"]
		ci = pubgo.Jf(ci, Ift)
		//<span class="badge bg-primary">21</span>
		str += "<a href=\"" + s.Host + "/" + Jlv + "/" + ci + Up + "\" title='什么是" + ci + "," + ci + "是什么意思'><span class='badge bg-secondary'>" + ci + "</span></a>  "
	}
	return template.HTML(str)
}
func (s *search) SetCidian(K, Kw string, Ift bool) template.HTML {
	rd := mysql.Selects("select zhushi from cidian where ci=?", K) //? ，不能管类型，只需?号，加''反而不对
	zs := "无此词条!"
	if len(rd) > 0 {
		zs = rd[0]["zhushi"]
		zs = pubgo.Jf(zs, Ift)
		zs = strings.Replace(zs, "\r\n", "<br>", -1)
		zs = "<b>" + Kw + "</b>" + "<br>" + zs
		return template.HTML(zs)
	} else {
		return template.HTML(pubgo.Jf(zs, Ift))
	}
}
