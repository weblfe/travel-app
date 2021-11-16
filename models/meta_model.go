package models

import "github.com/astaxie/beego"

type Meta struct {
	HasMore bool `json:"more"`
	P       int  `json:"page"`
	C       int  `json:"count"`
	Size    int  `json:"size"`
	Total   int  `json:"total"`
}

func (this *Meta) Page() int {
	return this.P
}

func (this *Meta) Count() int {
	return this.C
}

func (this *Meta) More() bool {
	return this.HasMore
}

func NewMeta() *Meta {
	var meta = new(Meta)
	meta.Init()
	return meta
}

func (this *Meta) Skip() int {
	panic("implement me")
}

func (this *Meta) SetTotal(i int) ListsParams {
	var m = NewMeta()
	m.Total = i
	return m
}

func (this *Meta) Set(key string, v interface{}) *Meta {
	switch key {
	case "more":
		this.HasMore = v.(bool)
	case "page":
		this.P = v.(int)
	case "count":
		this.C = v.(int)
	case "size":
		this.Size = v.(int)
	case "total":
		this.Total = v.(int)
	}
	return this
}

func (this *Meta) Init() *Meta {
	if this.P == 0 {
		this.P = 1
	}
	if this.C == 0 {
		this.C = 20
	}
	return this
}

func (this *Meta) Boot() {
	if this.Total > 0 {
		total := this.Total / this.C
		left := this.Total % this.C
		if total > this.P {
			this.HasMore = true
		}
		if total >= this.P && left > 0 {
			this.HasMore = true
		}
	}
}

func (this *Meta) M() beego.M {
	this.Boot()
	return beego.M{
		"more":  this.More(),
		"count": this.Count(),
		"page":  this.Page(),
		"size":  this.Size,
		"total": this.Total,
	}
}
