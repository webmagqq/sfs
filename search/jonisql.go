package search

import (
	"fmt"
	"sfsgo/gstr"
	"sfsgo/pubgo"
	"strconv"
	"strings"
)

type searchresult interface {
	getsql() string
}

type onesql struct {
	ak, Jlv, lwwh string
}

func Newonesql(ak, Jlv, lwwh string) *onesql {
	return &onesql{ak: ak, Jlv: Jlv, lwwh: lwwh}
}
func (o *onesql) getsql() string {
	ks := o.ak
	runek := []rune(ks)
	idx := runek[0]       //获取第一字的int
	tidx := int(idx) % 21 //--定位表
	kc := pubgo.Sublen(ks, 7)
	return "SELECT tid,jpid  from " + o.Jlv + "idx" + strconv.Itoa(tidx) + "  WHERE  ci like '" + kc + "%'  " + o.lwwh
}

//--------------------------------------------
type jonisql struct {
	ak, jz, Jlv, lwwh string
	Tjd               []map[string]string
}

func Newjonisql(jz, ak, Jlv, lwwh string, Tjd []map[string]string) *jonisql {
	return &jonisql{jz: jz, ak: ak, Jlv: Jlv, lwwh: lwwh, Tjd: Tjd}
}

func gpp(rd []map[string]string, jz string) (int, string) {
	//--根据搜索词量调节搜索方式，以保证不出现慢查询
	maxloop := len(rd)
	if maxloop > 4 {
		maxloop = 4 //--最多4个组合查询，最多没有实际意义。
	}

	wid := "tid" //--精准搜索 where条件
	if jz != "" {
		wid = "jpid" //--全搜索
	}
	ics, _ := strconv.Atoi(rd[0]["cs"])

	if ics > 5000 && maxloop > 3 { //第1个词大于10000的词的情况下，
		maxloop = 3
	}
	if ics > 10000 && maxloop > 2 { //第1个词大于10000的词的情况下，
		maxloop = 2
	}
	if ics > 3000 { //第1个词大于3000的词的情况下，必须如此，否则导致慢查询
		wid = "tid" //--where条件限制为同一句子上，这样精准且快
	}

	return maxloop, wid
}

func (j *jonisql) getsql() string {
	rd := j.Tjd
	maxloop, wid := gpp(rd, j.jz)
	ci, bpar, wh, whand := "", "", " where ", ""
	kNo := 0
	for i, v := range rd {
		if i >= maxloop {
			break
		}
		if v["ci"] == "" {
			continue
		}
		ci = v["ci"]
		if ci == j.ak {
			kNo = i //--用户搜索第一个词 在排序后的位置
		}
		ci = pubgo.Sublen(ci, 7) //最大长度7
		runek := []rune(ci)
		tbNo := int(runek[0]) % 21 //--定位表
		/*
			select DISTINCT a2.tid,a0.jpid from jingidx6 a0,jingidx19 a1,jingidx18 a2
			where a0.ci like '嗔%' and a1.ci like '痴%' and a2.ci like '贪%'
			and a1.jpid = a0.jpid and a2.jpid = a0.jpid limit 0,21
		*/
		//j.jlv + "idx" + tbNo.ToString() + " a" + i.ToString() + ","
		bpar += fmt.Sprintf("%sidx%d a%d,", j.Jlv, tbNo, i)
		//wh += " a" + i.ToString() + ".ci like '" + ci + "%' " + lmwh.Replace("dir", " a" + i.ToString()+".dir") + " and ";

		wh += fmt.Sprintf(" a%d.ci like '%s%s' %s and ", i, ci, "%", strings.Replace(j.lwwh, "dir", " a"+strconv.Itoa(i)+".dir", -1))
		if i > 0 {
			//jon += "  a" + i.ToString() + ".dlid = a0.dlid  and a" + i.ToString() + ".tid > a0.tid-" + dllen.ToString() + "  and a" + i.ToString() + ".tid < a0.tid+" + dllen.ToString() + "  and ";
			//jon += "  a" + i.ToString() + ".dlid = a0.dlid  and ";
			//"  a" + i.ToString() + ".jpid = a0.jpid  and "
			whand += fmt.Sprintf(" a%d.%s = a0.%s and ", i, wid, wid)
		}
	}
	sql := fmt.Sprintf("select DISTINCT a%d.tid,a%d.jpid from ", kNo, kNo)
	sql = sql + bpar[:len(bpar)-1] + wh + gstr.Do(whand, "", "and", false, true)
	return sql

}
