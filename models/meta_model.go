package models

import "github.com/astaxie/beego"

type Meta struct {
		More  bool `json:"more"`
		Page  int  `json:"page"`
		Count int  `json:"count"`
		Size  int  `json:"size"`
		Total int  `json:"total"`
}

func NewMeta() *Meta {
		var meta = new(Meta)
		meta.Init()
		return meta
}

func (this *Meta) Set(key string, v interface{}) *Meta {
		switch key {
		case "more":
				this.More = v.(bool)
		case "page":
				this.Page = v.(int)
		case "count":
				this.Count = v.(int)
		case "size":
				this.Size = v.(int)
		case "total":
				this.Total = v.(int)
		}
		return this
}

func (this *Meta) Init() *Meta {
		if this.Page == 0 {
				this.Page = 1
		}
		if this.Count == 0 {
				this.Count = 20
		}
		return this
}

func (this *Meta) Boot() {
		if this.Total > 0 {
				total := this.Total / this.Count
				left := this.Total % this.Count
				if total > this.Page {
						this.More = true
				}
				if total >= this.Page && left > 0 {
						this.More = true
				}
		}
}

func (this *Meta) M() beego.M {
		this.Boot()
		return beego.M{
				"more":  this.More,
				"count": this.Count,
				"page":  this.Page,
				"size":  this.Size,
				"total": this.Total,
		}
}
