package xianyan

import "time"

//该功能废弃
type Bws struct {
	Click, Day int
}

func NewBws() *Bws {
	return &Bws{}
}
func (b *Bws) Add() { //按天统计点击
	day := time.Now().Day()
	if day == b.Day {
		b.Click++
	} else {
		b.Day = day
		b.Click = 1
	}
}
