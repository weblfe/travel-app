package transports

import (
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/context"
		"github.com/weblfe/travel-app/models"
		"time"
)

// 游记请求
type TravelNoteRequest struct {
		Title         string    `json:"title" bson:"title"`
		Content       string    `json:"content" bson:"content"`
		Type          int       `json:"type" bson:"type"`
		Images        []string  `json:"images,omitempty" bson:"images,omitempty"`
		UserId        string    `json:"userId" bson:"userId"`
		Videos        []string  `json:"videos,omitempty" bson:"videos,omitempty"`
		Group         string    `json:"group" bson:"group"`
		Tags          []string  `json:"tags" bson:"tags"`
		Status        int       `json:"status" bson:"status"`
		Address       string    `json:"address" bson:"address"`
		Privacy       int       `json:"privacy" bson:"privacy"`
		UpdatedAt     time.Time `json:"updatedAt" bson:"updatedAt"`
		CreatedAt     time.Time `json:"createdAt" bson:"createdAt"`
		DeletedAt     int64     `json:"deletedAt" bson:"deletedAt"`
		transportImpl `json:",omitempty"`
}

func (this *TravelNoteRequest) Load(ctx *context.BeegoInput) *TravelNoteRequest {

		return this
}

func (this *TravelNoteRequest) M(m beego.M) beego.M {

		return m
}

func (this *TravelNoteRequest) Decode() *models.TravelNotes {
		var note = new(models.TravelNotes)
		note.UserId = this.UserId
		note.Type = this.Type
		note.Tags = this.Tags
		note.Address = this.Address
		note.Images = this.Images
		note.Videos = this.Videos
		note.Privacy = this.Privacy
		return note
}
