package pubgo

//公共函数库
import (
	"io/ioutil"
	"sfsgo/sat"
)

func Jf(str string, s bool) string { //简体转繁体
	if s {
		dicter := sat.DefaultDict()
		return dicter.ReadReverse(str)
	}
	return str
}
func Fj(str string, s bool) string { //繁体转体简
	if s {
		dicter := sat.DefaultDict()
		return dicter.Read(str)
	}
	return str
}
func Of(fn string) string {
	filetext, _ := ioutil.ReadFile(fn)
	return string(filetext[:])
}
func RemoveRepeatElement(list []string) []string { //--数组去重
	/*
		可以利用go中，map数据类型的key唯一的属性，来对数组去重
		将strSlice数组中重复的元素去掉，使其中的元素唯一
	*/
	// 创建一个临时map用来存储数组元素
	temp := make(map[string]bool)
	index := 0
	for _, v := range list {
		// 遍历数组元素，判断此元素是否已经存在map中
		_, ok := temp[v]
		if ok {
			list = append(list[:index], list[index+1:]...)
		} else {
			temp[v] = true
		}
		index++
	}
	return list
}
func Sublen(str string, l int) string { //--截取前 l 个字符串
	runek := []rune(str) //包含中文必须如此才能得到正确的长度
	k := ""
	if len(runek) > l {
		k = string(runek[:l]) //截取35位
	} else {
		k = str
	}
	return k
}
