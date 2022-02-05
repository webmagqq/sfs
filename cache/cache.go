package cache

import (
	"container/list"
	"strconv"
	"sync"
	"time"
)

type CacheData struct {
	mt    sync.RWMutex
	CdL   *list.List               //--双向链表，方便移动
	cache map[string]*list.Element //用元素绑定链表值，方便查找
	maxct int                      //cache最大个数
	cct   int                      //当前个数
	suc   int64                    //命中使用次数
	day   int                      //--按天统计命中
}
type entry struct {
	key   string //保存cache map[string]*list.Element 的key值，用于反过来删除map
	value interface{}
	suc   int64 //命中使用次数
}

//var m sync.RWMutex

func NewCache(maxct int) *CacheData {
	if maxct < 1 {
		panic("maxct 必须大于 1")
	}
	return &CacheData{
		CdL:   list.New(),
		cache: make(map[string]*list.Element),
		maxct: maxct,
		cct:   0,
		suc:   0,
	}
}

func (c *CacheData) Add(k string, v interface{}) {
	c.mt.Lock()
	defer c.mt.Unlock()
	if _, ok := c.cache[k]; !ok {
		//c.CdL.MoveToFront(vl) //MoveToFront将元素e移动到列表l的最前面。如果e不是l的元素，则不修改列表。元素不能为零。
		//} else {
		cv := c.CdL.PushFront(&entry{k, v, 0}) //在列表l的前面插入一个值为v的新元素e并返回e。
		c.cache[k] = cv
		c.cct += 1
	}
	for c.maxct < c.cct {
		c.del() //--删除链表值//删除map
	}
}

func (c *CacheData) Get(k string) interface{} {
	c.mt.Lock()
	defer c.mt.Unlock()
	if l, ok := c.cache[k]; ok {
		c.CdL.MoveToFront(l)   //MoveToFront将元素e移动到列表l的最前面。如果e不是l的元素，则不修改列表。元素不能为零。
		vl := l.Value.(*entry) //interface都用这种方式转换类型
		day := time.Now().Day()
		if day == c.day {
			c.suc += 1
		} else {
			c.day = day
			c.suc = 1
		}

		vl.suc += 1
		return vl.value //MoveToFront将元素e移动到列表l的最前面。如果e不是l的元素，则不修改列表。元素不能为零。
	}
	return nil
}
func (c *CacheData) del() {
	ele := c.CdL.Back() //返回链表的最后一个元素，返回类型元列表中的元素。
	ev := ele.Value.(*entry)
	delete(c.cache, ev.key) //删除map
	c.CdL.Remove(ele)       //如果e是列表l的元素，则删除从l中删除e。它返回元素值e.Value。元素不能为零。
	c.cct = c.cct - 1
}
func (c *CacheData) Gsuc() int64 {
	c.mt.RLock() //--纯读，用RLock，因为读的时候存在其他进程修改的可能性
	defer c.mt.RUnlock()
	return c.suc //命中使用次数
}

func (c *CacheData) Gets() string { //--测试函数
	rstr, key := "", ""
	var suc int64
	max := 100
	i := 0
	c.mt.RLock() //--纯读，用RLock，因为读的时候存在其他进程修改的可能性
	defer c.mt.RUnlock()
	for e := c.CdL.Front(); e != nil; e = e.Next() {
		if i >= max {
			break
		}
		key = e.Value.(*entry).key
		suc = e.Value.(*entry).suc
		rstr += strconv.Itoa(i) + "." + key + ":" + strconv.FormatInt(suc, 10) + "\n"
		i++
	}
	return rstr
}

/*
func (c *CacheData) GetLen() int64 { //--测试函数
	c.mt.RLock() //--纯读，用RLock，因为读的时候存在其他进程修改的可能性
	defer c.mt.RUnlock()
	rstr := ""
	var slen int64
	for e := c.CdL.Front(); e != nil; e = e.Next() {
		rstr = e.Value.(*entry).value.(string)
		slen += int64(len(rstr))
	}
	return slen
}
*/
