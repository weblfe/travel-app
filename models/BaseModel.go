package models

import (
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/logs"
		"github.com/globalsign/mgo"
		"github.com/globalsign/mgo/bson"
		"sync"
		"time"
)

type Model interface {
		Error() error
		GetConn(name ...string) *mgo.Session
		SetProfile(key string, v interface{}) Model
		Collection(conn ...string) *mgo.Collection
		Self() Model
		GetProfile(key string, defaults ...interface{}) interface{}
}

type IndexConstructor interface {
		CreateIndex()
}

type ListsParams interface {
		Page() int
		Count() int
		Skip() int
		SetTotal(int) ListsParams
		More() bool
}

type ListsParamImpl struct {
		page  int
		size  int
		total int
}

func NewListParam(page, size int) ListsParams {
		var param = new(ListsParamImpl)
		param.size = size
		param.page = page
		return param
}

type TableNameAble interface {
		TableName() string
}

type BaseModel struct {
		Err         error
		Connections map[string]*mgo.Session
		lock        sync.Mutex
		_Profiles   map[string]interface{}
		_Self       Model
		_Sess       []*mgo.Session
}

var (
		initLock      sync.Once
		readWriteLock sync.Mutex
		_baseProfiles *map[string]interface{}
)

func newBaseModel() *BaseModel {
		var base = new(BaseModel)
		return base
}

func GetModelProfiles() *map[string]interface{} {
		if _baseProfiles == nil {
				initLock.Do(func() {
						_baseProfiles = &map[string]interface{}{}
				})
		}
		return _baseProfiles
}

func SetProfile(key string, v interface{}) {
		info := GetModelProfiles()
		readWriteLock.Lock()
		(*info)[key] = v
		readWriteLock.Unlock()
}

func GetProfile(key string, defaults ...interface{}) interface{} {
		info := GetModelProfiles()
		readWriteLock.Lock()
		defer readWriteLock.Unlock()
		if v, ok := (*info)[key]; ok {
				return v
		}
		if len(defaults) == 0 {
				defaults = append(defaults, nil)
		}
		return defaults[0]
}

func (this *ListsParamImpl) Page() int {
		return this.page
}

func (this *ListsParamImpl) Count() int {
		return this.size
}

func (this *ListsParamImpl) Skip() int {
		return (this.page - 1) * this.size
}

func (this *ListsParamImpl) SetTotal(total int) ListsParams {
		this.total = total
		return this
}

func (this *ListsParamImpl) More() bool {
		page := this.total / this.size
		return page > this.page || this.total%this.size > 0 && page+1 > this.page
}

func (this *BaseModel) Init() {
		if this.Connections == nil {
				this.Connections = map[string]*mgo.Session{}
		}
		if this._Profiles == nil {
				this._Profiles = *GetModelProfiles()
		}
		if this._Sess == nil {
				this._Sess = make([]*mgo.Session, 2)
		}
		if this._Self != nil {
				if creator, ok := this._Self.(IndexConstructor); ok {
						creator.CreateIndex()
				}
		}
}

func (this *BaseModel) getConnUrl(name string) *mgo.DialInfo {
		if name == "default" {
				return &mgo.DialInfo{
						Addrs:    []string{this.getString("db_host", "127.0.0.1:27017")},
						Source:   this.getString("db_source", ""),
						Username: this.getString("db_username"),
						Password: this.getString("db_password"),
				}
		}
		return &mgo.DialInfo{
				Addrs:    []string{this.getString(name+".db_host", "127.0.0.1:27017")},
				Source:   this.getString(name+".db_source", ""),
				Username: this.getString(name + ".db_username"),
				Password: this.getString(name + ".db_password"),
		}
}

func (this *BaseModel) getString(key string, defaults ...string) string {
		if len(defaults) == 0 {
				defaults = append(defaults, "")
		}
		if v, ok := this._Profiles[key]; ok {
				if str, ok := v.(string); ok {
						return str
				}
		}
		return defaults[0]
}

func (this *BaseModel) GetConn(conn ...string) *mgo.Session {
		if len(conn) == 0 {
				conn = append(conn, "default")
		}
		if conn, ok := this.Connections[conn[0]]; ok {
				return conn
		}
		return this.conn(conn[0])
}

func (this *BaseModel) Self() Model {
		if this._Self == nil {
				return this
		}
		return this._Self
}

func (this *BaseModel) Collection(conn ...string) *mgo.Collection {
		if len(conn) == 0 {
				conn = append(conn, "default")
		}
		var sess = this.GetConn(conn...)
		if sess != nil {
				ins := sess.Copy()
				if table, ok := this.Self().(TableNameAble); ok {
						name := table.TableName()
						conn = append(conn, name)
				}
				if len(conn) < 2 {
						conn = append(conn, "default")
				}
				coll := ins.DB(this.GetDatabaseName(conn[0])).C(conn[1])
				// 添加到回收sess
				if coll != nil {
						this._Sess = append(this._Sess, ins)
				}
				return coll
		}
		return nil
}

func (this *BaseModel) GetDatabaseName(conn ...string) string {
		if len(conn) == 0 {
				conn = append(conn, "default")
		}
		var name = conn[0]
		if name == "default" || name == "" {
				return this.getString("db_name", "mongodb")
		}
		return this.getString(name+".db_name", "mongodb")
}

func (this *BaseModel) conn(name string) *mgo.Session {
		this.lock.Lock()
		defer this.lock.Unlock()
		if sess, err := mgo.DialWithInfo(this.getConnUrl(name)); err == nil {
				this.Connections[name] = sess
				return sess
		}
		return nil
}

func (this *BaseModel) Error() error {
		return this.Err
}

func (this *BaseModel) Destroy() {
		if this.Err != nil {
				logs.Error(this.Err)
				this.Err = nil
		}
		if this.Connections != nil && len(this.Connections) != 0 {
				for key, conn := range this.Connections {
						conn.Close()
						delete(this.Connections, key)
				}
				this.Connections = map[string]*mgo.Session{}
		}
		if this._Profiles == nil {
				this._Profiles = nil
		}
}

func (this *BaseModel) SetProfile(key string, v interface{}) Model {
		this.lock.Lock()
		this._Profiles[key] = v
		this.lock.Unlock()
		return this
}

func (this *BaseModel) GetProfile(key string, defaults ...interface{}) interface{} {
		this.lock.Lock()
		if v, ok := this._Profiles[key]; ok {
				return v
		}
		this.lock.Unlock()
		if len(defaults) == 0 {
				defaults = append(defaults, nil)
		}
		return defaults[0]
}

// 回收session
func (this *BaseModel) destroy() {
		for _, sess := range this._Sess {
				if sess != nil {
						sess.Close()
				}
		}
		this._Sess = this._Sess[0:0]
}

func (this *BaseModel) Add(docs interface{}) error {
		table := this.Collection()
		defer this.destroy()
		return table.Insert(docs)
}

func (this *BaseModel) Insert(docs interface{}) error {
		table := this.Collection()
		defer this.destroy()
		return table.Insert(docs)
}

func (this *BaseModel) Inserts(docs []interface{}) error {
		if len(docs) == 0 {
				return nil
		}
		table := this.Collection()
		defer this.destroy()
		return table.Insert(docs...)
}

func (this *BaseModel) GetByKey(key string, v interface{}, result interface{}) error {
		table := this.Collection()
		defer this.destroy()
		return table.Find(bson.M{key: v}).One(result)
}

func (this *BaseModel) GetById(id string, result interface{}, selects ...interface{}) error {
		table := this.Collection()
		defer this.destroy()
		if len(selects) > 0 {
				return table.Find(beego.M{
						"_id": bson.ObjectIdHex(id),
				}).Select(selects[0]).One(result)
		}
		return table.Find(beego.M{
				"_id": bson.ObjectIdHex(id),
		}).One(result)
}

func (this *BaseModel) Update(query interface{}, data interface{}) error {
		table := this.Collection()
		defer this.destroy()
		return table.Update(query, data)
}

func (this *BaseModel) Updates(query interface{}, data interface{}) (*mgo.ChangeInfo, error) {
		table := this.Collection()
		defer this.destroy()
		return table.UpdateAll(query, data)
}

func (this *BaseModel) All(query interface{}, result interface{}, selects ...interface{}) error {
		table := this.Collection()
		defer this.destroy()
		if len(selects) > 0 {
				return table.Find(query).Select(selects[0]).All(result)
		}
		return table.Find(query).All(result)
}

func (this *BaseModel) Remove(query interface{}, softDelete ...bool) error {
		table := this.Collection()
		defer this.destroy()
		if len(softDelete) == 0 {
				softDelete = append(softDelete, true)
		}
		if softDelete[0] {
				return this.Update(query, beego.M{"deleted_at": time.Now().Unix()})
		}
		return table.Remove(query)
}

func (this *BaseModel) Deletes(query interface{}, softDelete ...bool) (*mgo.ChangeInfo, error) {
		table := this.Collection()
		defer this.destroy()
		if len(softDelete) == 0 {
				softDelete = append(softDelete, true)
		}
		if softDelete[0] {
				return this.Updates(query, beego.M{"deleted_at": time.Now().Unix()})
		}
		return table.RemoveAll(query)
}

func (this *BaseModel) Lists(query interface{}, result interface{}, limit ListsParams, selects ...interface{}) (int, error) {
		table := this.Collection()
		var (
				size, skip = limit.Count(), limit.Skip()
		)
		defer this.destroy()
		total, err := table.Find(query).Count()
		if err != nil {
				return 0, err
		}
		if total == 0 {
				return 0, nil
		}
		if len(selects) > 0 {
				return total, table.Find(query).Select(selects[0]).Limit(size).Skip(skip).All(result)
		}
		return total, table.Find(query).Limit(limit.Count()).Skip(skip).All(result)
}
