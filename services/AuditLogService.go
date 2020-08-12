package services

import (
		"errors"
		"github.com/weblfe/travel-app/models"
)

type AuditLogService interface {
		Adds(userId, typ, comment string, ids []string) error
}

type auditLogServiceImpl struct {
		BaseService
		model *models.AuditLogModel
}

func AuditLogServiceOf() AuditLogService {
		var service = newAuditLogService()
		return service
}

func newAuditLogService() *auditLogServiceImpl {
		var service = new(auditLogServiceImpl)
		service.Init()
		return service
}

func (this *auditLogServiceImpl) Init() {
		this.ClassName = "auditLogServiceImpl"
		this.Constructor = func(args ...interface{}) interface{} {
				return AuditLogServiceOf()
		}
		this.model = models.AuditLogModelOf()
		this.init()
}

func (this *auditLogServiceImpl) Adds(userId, typ, comment string, ids []string) error {
		var docs []interface{}
		for _, id := range ids {
				if id == "" {
						continue
				}
				log := models.NewAuditLog()
				log.UserId = userId
				log.Comment = comment
				log.AuditType = typ
				log.PostId = id
				log.InitDefault()
				docs = append(docs, log)
		}
		if len(docs) <= 0 {
				return errors.New("errors empty adds")
		}
		return this.model.Inserts(docs)
}
