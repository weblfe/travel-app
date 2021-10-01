package kernel

type cdnDomain struct {
	collection string
	conn  string
	db    string
}

func NewCdnDomain() *cdnDomain  {
	var domain = new(cdnDomain)
	domain.db = envOr("APP_DATABASE_NAME","travel")
	domain.collection = envOr("ATTACHMENT_COLLECTION_NAME","attachments")
	return domain
}

func (domain *cdnDomain)SetDomainUrl(url string) int64 {
	var (
		connection = GetDbMgr().GetDb(domain.conn)
	    db = connection.DB(domain.db)
	)
	iter:=db.C(domain.collection).NewIter(connection,nil,100,nil)
	_ = iter.For(nil, func() error {

		return nil
	})
	return 0
}
