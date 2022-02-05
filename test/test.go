package test

import (
	"sfsgo/pubgo"
	"sfsgo/search"
	"strconv"
)

type gif func() string

type ifos struct {
	Mdo map[string]gif
}

func Newifos() *ifos {
	Mdo := make(map[string]gif)
	Mdo["1"] = getCacheSearchData
	Mdo["2"] = gettongji
	Mdo["3"] = getCacheSearchsuc
	return &ifos{Mdo: Mdo}
}
func getCacheSearchData() string {
	ch := search.CacheSearchData
	return ch.Gets()
}
func getCacheSearchsuc() string {
	ch := search.CacheSearchData
	return "缓存命中：" + strconv.FormatInt(ch.Gsuc(), 10)
}
func gettongji() string {
	ifo := ""
	sum := 0
	for k, v := range pubgo.Tj.Tjs {
		sum += v.Bws
		ifo += k + ":" + strconv.Itoa(v.Bws) + "\n"
	}
	return ifo + "总计：" + strconv.Itoa(sum)
}
