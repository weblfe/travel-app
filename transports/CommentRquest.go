package transports

import (
		"encoding/gob"
		"github.com/astaxie/beego/context"
		"github.com/astaxie/beego/logs"
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/libs"
		"github.com/weblfe/travel-app/models"
)

// 评论请求
type CommentRequest struct {
		TargetType    string `json:"targetType" json:"type"` // 评论类型
		TargetId      string `json:"targetId"`               // 评论目表 ID
		Content       string `json:"content"`                // 评论内容
		transportImpl `json:",omitempty"`
}

// 评论
func NewCommentInstance() *CommentRequest {
		var comment = new(CommentRequest)
		comment.init()
		return comment
}

// 评论
func NewComment(ctx ...*context.BeegoInput) *CommentRequest {
		var comment = NewCommentInstance()
		if len(ctx) > 0 {
				comment.Load(ctx[0]).Init()
		}
		return comment
}

func (this *CommentRequest) init() {
		this.AppendInit(this.defaults)
		gob.Register(CommentRequest{})
}

func (this *CommentRequest) defaults() {
		if this.TargetType == "" {
				this.TargetType = "post"
		}
}

func (this *CommentRequest) Load(ctx *context.BeegoInput) *CommentRequest {
		var err = this.Decoder(ctx, this)
		if err != nil {
				logs.Error(err)
		}
		return this
}

func (this *CommentRequest) GobEncode() ([]byte, error) {
		return libs.Json().Marshal(this)
}

func (this *CommentRequest) filter(m bson.M) bson.M {
		return m
}

func (this *CommentRequest) Decode() *models.Comment {
		var comment = new(models.Comment)
		if this.Content == "" {
				return nil
		}
		comment.Content = this.Content
		comment.TargetType = this.TargetType
		comment.TargetId = this.TargetId
		comment.Defaults()
		return comment
}
