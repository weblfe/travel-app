package kernel

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
	"strings"
)

type cdnDomain struct {
	collection string
	conn       string
	db         string
}

func NewCdnDomain() *cdnDomain {
	var domain = new(cdnDomain)
	domain.db = envOr("APP_DATABASE_NAME", "travel")
	domain.collection = envOr("ATTACHMENT_COLLECTION_NAME", "attachments")
	return domain
}

func (domain *cdnDomain) SetDomainUrl(url string) int64 {
	var (
		connection = GetDbMgr().GetDb(domain.conn)
		db         = connection.DB(domain.db)
		doc        = make(bson.M)
		query      = bson.M{}
		collection = db.C(domain.collection)
	)
	var (
		result int64
		iter   = collection.Find(query).Batch(10).Iter()
	)
	for iter.Next(&doc) {
		var (
			cdnUrl, _ = doc["cdnUrl"]
			newUrl    = domain.replaceDomainUrl(fmt.Sprintf("%v", cdnUrl), url)
		)
		fmt.Println(fmt.Sprintf("update: id:%v", doc["_id"]))
		doc["cdnUrl"] = newUrl
		if err := collection.Update(bson.M{"_id": doc["_id"]}, doc); err != nil {
			fmt.Println("error:", err.Error())
		}
		result++
	}
	return result
}

func (domain *cdnDomain) replaceDomainUrl(url string, newDomainUrl string) string {
	if url == "" {
		return ""
	}
	var (
		arr    = strings.Split(url, "/storage")
		newUrl = strings.Replace(url, arr[0], newDomainUrl, 1)
	)
	fmt.Println(fmt.Sprintf("oldUrl:%s, newUrl:%s", url, newUrl))
	return newUrl
}
