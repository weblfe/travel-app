package transports

import (
		"github.com/astaxie/beego/context"
		"github.com/globalsign/mgo/bson"
)

// 评论请求
type SearchRequest struct {
		Query         string `json:"query"`  // 评论类型
		Sort          string `json:"sort"`   // 搜索排序
		Search        bson.M `json:"search"` // 搜索对象
		transportImpl `json:",omitempty"`
}

func NewSearchInstance() *SearchRequest {
		var search = new(SearchRequest)
		return search
}

func NewSearchQuery(ctx ...*context.BeegoInput) *SearchRequest {
		if len(ctx) > 0 {
				var instance = NewSearchInstance().Load(ctx[0])
				instance.Init()
				return instance
		}
		return NewSearchInstance()
}

func (this *SearchRequest) Load(ctx *context.BeegoInput) *SearchRequest {

		return this
}

func (this *SearchRequest) filter(m bson.M) bson.M {

		return m
}
