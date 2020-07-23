package transports

import (
		"github.com/astaxie/beego/context"
		"time"
)

// 游记请求
type TravelNoteRequest struct {
		Title     string        `json:"title" bson:"title"`
		Content   string        `json:"content" bson:"content"`
		Type      int           `json:"type" bson:"type"`
		Images    []string      `json:"images,omitempty" bson:"images,omitempty"`
		UserId    string        `json:"userId" bson:"userId"`
		Videos    []string      `json:"videos,omitempty" bson:"videos,omitempty"`
		Group     string        `json:"group" bson:"group"`
		Tags      []string      `json:"tags" bson:"tags"`
		Status    int           `json:"status" bson:"status"`
		Address   string        `json:"address" bson:"address"`
		Privacy   int           `json:"privacy" bson:"privacy"`
		UpdatedAt time.Time     `json:"updatedAt" bson:"updatedAt"`
		CreatedAt time.Time     `json:"createdAt" bson:"createdAt"`
		DeletedAt int64         `json:"deletedAt" bson:"deletedAt"`
		transportImpl `json:",omitempty"`
}

func (this *TravelNoteRequest)Load(ctx *context.BeegoInput) *TravelNoteRequest  {

		return this
}

