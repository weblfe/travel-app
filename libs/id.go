package libs

import (
		"github.com/astaxie/beego/logs"
		"github.com/globalsign/mgo"
		"github.com/globalsign/mgo/bson"
		"sync"
		"time"
)

type CommitId struct {
		storage   sync.Map
		Db        string
		TableName string
		Conn      *mgo.Session
		_Self     *mgo.Session
		txnId     bson.ObjectId
		locker    sync.Locker
		id        *ids
}

var (
		locks sync.Map
)

type ids struct {
		Id         bson.ObjectId `json:"_id" bson:"_id"`
		Collection string        `json:"collection" bson:"collection"`
		NumId      int64         `json:"inc_id" bson:"inc_id"`
		State      int           `json:"state" bson:"state"`
		CreatedAt  time.Time     `json:"created_at" bson:"created_at"`
		EndAt      int64         `json:"end_at" bson:"end_at"`
}

func NewId(tableName string, db string, conn *mgo.Session) *CommitId {
		var commitId = new(CommitId)
		commitId.TableName = tableName
		commitId.Conn = conn
		commitId.Db = db
		commitId.init()
		return commitId
}

func getLock(name string) sync.Locker {
		if v, ok := locks.Load(name); ok {
				if ins, ok := v.(sync.Locker); ok {
						return ins
				}
		}
		var lock = &sync.Mutex{}
		locks.Store(name, lock)
		return lock
}

// 用户表自增ID
func GetId(db, collection string, sess *mgo.Session) int64 {
		txn := NewId(collection, db, sess)
		defer func() {
				err := txn.Commit()
				if err != nil {
						logs.Error(err)
				}
		}()
		return txn.GetId()
}

func (this *CommitId) init() {
		this.txnId = bson.NewObjectId()
		this.locker = getLock(this.TableName)
		ins := &ids{}
		this._Self = this.Conn.Copy()
		this.locker.Lock()
		defer this.locker.Unlock()
		err := this.collection().Find(bson.M{"collection": this.TableName}).Sort("-_id").One(ins)
		if err == nil || err == mgo.ErrNotFound {
				ins.Id = this.txnId
				ins.NumId = ins.NumId + 1
				ins.State = 0
				ins.Collection = this.TableName
				ins.CreatedAt = time.Now()
				ins.EndAt = 0
				_ = this.collection().Insert(ins)
				this.id = ins
				return
		}
		panic(err)
}

func (this *CommitId) GetId() int64 {
		return this.id.NumId
}

func (this *CommitId) Commit() error {
		defer this._Self.Close()
		this.id.State = 1
		this.id.EndAt = time.Now().Unix()
		return this.collection().Update(bson.M{"_id": this.txnId}, this.id)
}

func (this *CommitId) RollBack() error {
		defer this._Self.Close()
		this.id.State = -1
		this.id.EndAt = time.Now().Unix()
		return this.collection().Update(bson.M{"_id": this.txnId}, this.id)
}

func (this *CommitId) collection() *mgo.Collection {
		return this._Self.DB(this.Db).C("commit_id")
}
