package models

import (
		"errors"
		"fmt"
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/cache"
		"github.com/astaxie/beego/logs"
		"github.com/globalsign/mgo"
		"github.com/globalsign/mgo/bson"
		"github.com/globalsign/mgo/txn"
		"math"
		"math/rand"
		"reflect"
		"strconv"
		"sync"
		"time"
)

type Model interface {
		Error() error
		Self() Model
		Bind(v interface{}) Model
		GetConn(name ...string) *mgo.Session
		SetProfile(key string, v interface{}) Model
		Db(db ...string) *mgo.Database
		Document() *mgo.Collection
		Collection(conn ...string) *mgo.Collection
		GetProfile(key string, defaults ...interface{}) interface{}
}

type IndexConstructor interface {
		CreateIndex(force ...bool)
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

var (
		GlobalMgoSessionContainer = sync.Map{}
)

type BaseModel struct {
		Err        error
		lock       sync.Mutex
		_Binder    Model
		_Scope     bson.M
		dbName     string
		document   string
		serverName string
		_Profiles  map[string]interface{}
}

// 配置
type ConnOption struct {
		Db       string `json:"db"`     // 数据库
		Server   string `json:"server"` // 配置类型
		Document string `json:"table"`  // 表,文档
}

// 事务上下文
type TxnContext struct {
		TxnOps    []txn.Op
		TxnRunner *txn.Runner
		TxnId     bson.ObjectId
		TxnResult interface{}
}

var (
		initLock      sync.Once
		readWriteLock sync.Mutex
		_baseProfiles *map[string]interface{}
		ErrEmptyData  = errors.New("empty data call")
)

func newBaseModel() *BaseModel {
		var base = new(BaseModel)
		return base
}

// 获取相关配置
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
		if this._Profiles == nil {
				this._Profiles = *GetModelProfiles()
		}
		if this._Binder != nil {
				if creator, ok := this._Binder.(IndexConstructor); ok {
						creator.CreateIndex()
				}
		}
}

func (this *BaseModel) getConnOptions(name string) *mgo.DialInfo {
		if name == "default" {
				return &mgo.DialInfo{
						Addrs:       []string{this.getString("db_host", "127.0.0.1:27017")},
						Source:      this.getString("db_source", ""),
						Username:    this.getString("db_username"),
						Password:    this.getString("db_password"),
						PoolLimit:   this.getInt("db_pool_limit"),
						Timeout:     this.getDuration("db_timeout"),
						PoolTimeout: this.getDuration("db_pool_timeout"),
				}
		}
		return &mgo.DialInfo{
				Addrs:       []string{this.getString(name+".db_host", "127.0.0.1:27017")},
				Source:      this.getString(name+".db_source", ""),
				Username:    this.getString(name + ".db_username"),
				Password:    this.getString(name + ".db_password"),
				PoolLimit:   this.getInt(name + ".db_pool_limit"),
				Timeout:     this.getDuration(name + ".db_timeout"),
				PoolTimeout: this.getDuration(name + ".db_pool_timeout"),
		}
}

func (this *BaseModel) getDuration(key string, defaults ...time.Duration) time.Duration {
		if len(defaults) == 0 {
				defaults = append(defaults, 0)
		}
		if v, ok := this._Profiles[key]; ok {
				if d, ok := v.(time.Duration); ok {
						return d
				}
				d, ok := v.(string)
				if ok {
						t, err := time.ParseDuration(d)
						if err == nil {
								return t
						}
						logs.Error(err)
				}
				if d, ok := v.(int64); ok {
						return time.Duration(d)
				}
				if d, ok := v.(int); ok {
						return time.Duration(d)
				}
		}
		return defaults[0]
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

func (this *BaseModel) getInt(key string, defaults ...int) int {
		if len(defaults) == 0 {
				defaults = append(defaults, 0)
		}
		if v, ok := this._Profiles[key]; ok {
				if num, ok := v.(int); ok {
						return num
				}
				if num, ok := v.(string); ok {
						if n, err := strconv.Atoi(num); err == nil {
								return n
						}
				}
		}
		return defaults[0]
}

func (this *BaseModel) GetConn(conn ...string) *mgo.Session {
		return this.conn(arrayFirst(conn, this.getServerName()))
}

func (this *BaseModel) Self() Model {
		if this._Binder == nil {
				return this
		}
		return this._Binder
}

func (this *BaseModel) Collection(doc ...string) *mgo.Collection {
		if len(doc) > 0 {
				var (
						db         = this.Db()
						collection = db.C(arrayFirst(doc, this.getDoc()))
				)
				if collection == nil {
						panic(errors.New("collection error"))
				}
				return collection
		}
		return this.Document()
}

func (this *BaseModel) Db(db ...string) *mgo.Database {
		var dbName = arrayFirst(db, this.getDb())
		return this.Server().DB(dbName)
}

func (this *BaseModel) getDb() string {
		if this.dbName == "" {
				return this.getString("db_name", "default")
		}
		return this.dbName
}

func (this *BaseModel) Server(server ...string) *mgo.Session {
		return this.getServer(arrayFirst(server, this.getServerName()))
}

func (this *BaseModel) getServer(serverName string) *mgo.Session {
		if this.serverName == "" {
				this.serverName = serverName
		}
		var server = this.getSess(serverName)
		if server != nil {
				return server
		}
		return this.connection(serverName)
}

func (this *BaseModel) connection(server string) *mgo.Session {
		var sess, err = mgo.DialWithInfo(this.getConnOptions(server))
		if err != nil {
				panic(err)
		}
		GlobalMgoSessionContainer.Store(server, sess)
		defer this.waiting(server)
		return sess
}

func (this *BaseModel) getServerName() string {
		if this.serverName == "" {
				return this.getString("server", "default")
		}
		return this.serverName
}

func (this *BaseModel) Document() *mgo.Collection {
		var docName = this.getDoc()
		if docName == "" {
				panic(errors.New("document name empty"))
		}
		var (
				db         = this.Db()
				collection = db.C(docName)
		)
		if collection == nil {
				panic(errors.New("collection error"))
		}
		return collection
}

func (this *BaseModel) getDoc() string {
		var doc = this.document
		if doc == "" {
				if this._Binder == nil {
						return "documents"
				}
				if table, ok := this._Binder.(TableNameAble); ok {
						doc = table.TableName()
				}
		}
		if doc == "" {
				doc = this.getString("db_default_documents", "documents")
		}
		return doc
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

func (this *BaseModel) getSess(name string) *mgo.Session {
		var conn, ok = GlobalMgoSessionContainer.Load(name)
		if ok && conn != nil {
				if sess := conn.(*mgo.Session); sess != nil {
						return sess
				}
		}
		return this.connection(name)
}

func (this *BaseModel) conn(servers ...string) *mgo.Session {
		var (
				server   = arrayFirst(servers, "default")
				conn, ok = GlobalMgoSessionContainer.Load(server)
		)
		if ok && conn != nil {
				if sess, ok := conn.(*mgo.Session); ok {
						return sess
				}
				panic(errors.New("connection type error"))
		}
		var (
				info      = this.getConnOptions(server)
				sess, err = mgo.DialWithInfo(info)
		)
		if err != nil {
				panic(err)
		}
		return sess
}

func (this *BaseModel) Error() error {
		return this.Err
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
func (this *BaseModel) destroy(document *mgo.Collection) {
		this.done(document.Name)
		this.Release()
}

func (this *BaseModel) done(name string) {
		var (
				key  = this.getRefCountName(name)
				n, _ = GlobalMgoSessionContainer.Load(key)
		)
		if n == 0 {
				return
		}
		if n == nil {
				GlobalMgoSessionContainer.Store(key, 0)
				return
		}
		GlobalMgoSessionContainer.Store(key, n.(int)-1)
}

func (this *BaseModel) waiting(name string) {
		var (
				key  = this.getRefCountName(name)
				n, _ = GlobalMgoSessionContainer.Load(key)
		)
		if n == nil {
				GlobalMgoSessionContainer.Store(key, 1)
				return
		}
		GlobalMgoSessionContainer.Store(key, n.(int)+1)
}

// 释放
func (this *BaseModel) Release() {
		if this.serverName == "" {
				return
		}
		var (
				refCountName = this.getRefCountName(this.serverName)
				refCount, _  = GlobalMgoSessionContainer.Load(refCountName)
				db, _        = GlobalMgoSessionContainer.Load(this.serverName)
		)
		// db 不存在
		if db == nil {
				// 计数清库
				if refCount != nil {
						GlobalMgoSessionContainer.Delete(refCountName)
				}
				return
		}
		count := refCount.(int)
		if count > 0 {
				return
		}
		GlobalMgoSessionContainer.Delete(refCountName)
		GlobalMgoSessionContainer.Delete(this.serverName)
		this.close(db)
}

func (this *BaseModel) close(db interface{}) {
		if closer, ok := db.(*mgo.Session); ok {
				closer.Close()
		}
}

func (this *BaseModel) getRefCountName(server string) string {
		return server + "_" + "ref_count"
}

func (this *BaseModel) Add(docs interface{}) error {
		var table = this.Document()
		return table.Insert(docs)
}

func (this *BaseModel) Insert(docs interface{}) error {
		var table = this.Document()
		defer this.destroy(table)
		if m, ok := docs.(MapperAble); ok {
				return table.Insert(m.M())
		}
		return table.Insert(docs)
}

func (this *BaseModel) Inserts(docs []interface{}) error {
		if len(docs) == 0 {
				return nil
		}
		var table = this.Document()
		defer this.destroy(table)
		return table.Insert(docs...)
}

func (this *BaseModel) GetByKey(key string, v interface{}, result interface{}) error {
		var table = this.Document()
		defer this.destroy(table)
		defer this.resetScopeQuery()
		return table.Find(this.UseScopeQuery(bson.M{key: v})).One(result)
}

func (this *BaseModel) GetById(id string, result interface{}, selects ...interface{}) error {
		if id == "" {
				return mgo.ErrNotFound
		}
		var table = this.Document()
		defer this.destroy(table)
		defer this.resetScopeQuery()
		if len(selects) > 0 {
				return table.Find(this.UseScopeQuery(bson.M{
						"_id": bson.ObjectIdHex(id),
				})).Select(selects[0]).One(result)
		}
		return table.Find(this.UseScopeQuery(bson.M{
				"_id": bson.ObjectIdHex(id),
		})).One(result)
}

// 使用局部查询
func (this *BaseModel) UseScopeQuery(m bson.M) bson.M {
		if this._Scope == nil || len(this._Scope) == 0 {
				return m
		}
		for key, v := range this._Scope {
				m[key] = v
		}
		return m
}

// 添加局部查询
func (this *BaseModel) AddScopeQuery(m bson.M) {
		if len(m) == 0 {
				return
		}
		if this._Scope == nil {
				this._Scope = m
				return
		}
		for key, v := range m {
				this._Scope[key] = v
		}
		return
}

// 释放局部查询条件
func (this *BaseModel) resetScopeQuery() {
		this._Scope = nil
}

func (this *BaseModel) GetByObjectId(id bson.ObjectId, result interface{}, selects ...interface{}) error {
		var table = this.Document()
		defer this.destroy(table)
		defer this.resetScopeQuery()
		if len(selects) > 0 {
				return table.Find(this.UseScopeQuery(bson.M{
						"_id": id,
				})).Select(selects[0]).One(result)
		}
		return table.Find(this.UseScopeQuery(bson.M{
				"_id": id,
		})).One(result)
}

func (this *BaseModel) Update(query interface{}, data interface{}) error {
		var table = this.Document()
		data = this.setUpdate(data)
		defer this.destroy(table)
		return table.Update(query, data)
}

// 更新
func (this *BaseModel) setUpdate(data interface{}) interface{} {
		var newData = make(beego.M)
		if m, ok := data.(beego.M); ok {
				if _, ok := m["$set"]; ok {
						return data
				}
				newData["$set"] = data
				return newData
		}
		if m, ok := data.(bson.M); ok {
				if _, ok := m["$set"]; ok {
						return data
				}
				newData["$set"] = data
				return newData
		}
		if m, ok := data.(map[string]interface{}); ok {
				if _, ok := m["$set"]; ok {
						return data
				}
				newData["$set"] = data
				return newData
		}
		var t = reflect.TypeOf(data)
		if t.Kind() == reflect.Struct || t.Elem().Kind() == reflect.Struct {
				newData["$set"] = data
		}
		return data
}

func (this *BaseModel) UpdateById(id string, data interface{}) error {
		var table = this.Document()
		defer this.destroy(table)
		defer this.resetScopeQuery()
		data = this.setUpdate(data)
		return table.Update(this.UseScopeQuery(bson.M{"_id": bson.ObjectIdHex(id)}), data)
}

func (this *BaseModel) Updates(query interface{}, data interface{}) (*mgo.ChangeInfo, error) {
		var table = this.Document()
		defer this.destroy(table)
		data = this.setUpdate(data)
		return table.UpdateAll(query, data)
}

func (this *BaseModel) All(query interface{}, result interface{}, selects ...interface{}) error {
		var table = this.Document()
		defer this.destroy(table)
		query = this.WrapperScopeQuery(query)
		defer this.resetScopeQuery()
		if len(selects) > 0 {
				return table.Find(query).Select(selects[0]).All(result)
		}
		return table.Find(query).All(result)
}

func (this *BaseModel) Remove(query interface{}, softDelete ...bool) error {
		var table = this.Document()
		defer this.destroy(table)
		if len(softDelete) == 0 {
				softDelete = append(softDelete, true)
		}
		if softDelete[0] {
				return this.Update(query, beego.M{"deleted_at": time.Now().Unix()})
		}
		return table.Remove(query)
}

func (this *BaseModel) Deletes(query interface{}, softDelete ...bool) (*mgo.ChangeInfo, error) {
		var table = this.Document()
		defer this.destroy(table)
		if len(softDelete) == 0 {
				softDelete = append(softDelete, true)
		}
		if softDelete[0] {
				return this.Updates(query, beego.M{"deleted_at": time.Now().Unix()})
		}
		return table.RemoveAll(query)
}

func (this *BaseModel) Lists(query interface{}, result interface{}, limit ListsParams, selects ...interface{}) (int, error) {
		var table = this.Document()
		var (
				size, skip = limit.Count(), limit.Skip()
		)
		defer this.destroy(table)
		query = this.WrapperScopeQuery(query)
		defer this.resetScopeQuery()
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

func (this *BaseModel) ListsQuery(query interface{}, limit ListsParams, selects ...interface{}) *mgo.Query {
		var table = this.Document()
		var (
				size, skip = 0, 0
		)
		if limit != nil {
				size, skip = limit.Count(), limit.Skip()
		}
		query = this.WrapperScopeQuery(query)
		defer this.resetScopeQuery()
		if len(selects) > 0 {
				if limit == nil {
						return table.Find(query).Select(selects[0])
				}
				return table.Find(query).Select(selects[0]).Limit(size).Skip(skip)
		}
		if limit == nil {
				return table.Find(query)
		}
		return table.Find(query).Limit(limit.Count()).Skip(skip)
}

func (this *BaseModel) FindOne(query interface{}, result interface{}, selects ...interface{}) error {
		var table = this.Document()
		query = this.WrapperScopeQuery(query)
		defer this.destroy(table)
		defer this.resetScopeQuery()
		if len(selects) > 0 {
				return table.Find(query).Select(selects[0]).One(result)
		}
		return table.Find(query).One(result)
}

// 查询作用域
func (this *BaseModel) WrapperScopeQuery(query interface{}) interface{} {
		if this._Scope == nil || len(this._Scope) == 0 {
				return query
		}
		switch query.(type) {
		case bson.M:
				return this.UseScopeQuery(query.(bson.M))
		case beego.M:
				return this.UseScopeQuery(bson.M(query.(beego.M)))
		case map[string]interface{}:
				return this.UseScopeQuery(query.(map[string]interface{}))
		}
		return query
}

func (this *BaseModel) NewQuery(query bson.M) *mgo.Query {
		var table = this.Document()
		return table.Find(query)
}

func (this *BaseModel) Gets(query interface{}, result interface{}, selects ...interface{}) error {
		var table = this.Document()
		defer this.destroy(table)
		query = this.WrapperScopeQuery(query)
		defer this.resetScopeQuery()
		if len(selects) > 0 {
				return table.Find(query).Select(selects[0]).All(result)
		}
		return table.Find(query).All(result)
}

func (this *BaseModel) Count(query interface{}) int {
		var table = this.Document()
		defer this.destroy(table)
		defer this.resetScopeQuery()
		query = this.WrapperScopeQuery(query)
		n, err := table.Find(query).Count()
		if err == nil {
				return n
		}
		return 0
}

// [
//  '$match' => $map],
//   [
//   '$group' => [
//        '_id' => null,
//        'total_money' => ['$sum' => '$money'],
//        'total_money_usd' => ['$sum' => '$money_usd']
//    ]
// ]
func (this *BaseModel) Sum(query bson.M, sum string) int {
		var table = this.Document()
		defer this.destroy(table)
		query = this.UseScopeQuery(query)
		var (
				resultPipe struct {
						ID bson.ObjectId `bson:"_id"`
						C  int           `bson:"c"`
				}
				pipe = []bson.M{
						{"$match": query},
						{
								"$group": bson.M{
										"_id": nil,
										"c": bson.M{
												"$sum": "$" + sum + "",
										},
								},
						},
				}
		)

		err := table.Pipe(pipe).One(&resultPipe)
		if err == nil {
				return resultPipe.C
		}
		return 0
}

func (this *BaseModel) Exists(query interface{}) bool {
		var table = this.Document()
		defer this.destroy(table)
		query = this.WrapperScopeQuery(query)
		var tmp beego.M
		if err := table.Find(query).One(&tmp); err == nil {
				if _, ok := tmp["_id"]; ok {
						return true
				}
				return len(tmp) > 0
		}
		return false
}

func (this *BaseModel) IncrBy(query interface{}, incr interface{}) error {
		var table = this.Document()
		defer this.destroy(table)
		return table.Update(query, bson.M{"$inc": incr})
}

// 记录不存在异常
func (this *BaseModel) IsNotFound(err error) bool {
		return IsNotFound(err)
}

// 游标异常
func (this *BaseModel) IsErrCursor(err error) bool {
		return IsErrCursor(err)
}

// 使用软删除查询
func (this *BaseModel) UseSoftDelete(status ...int64) {
		if len(status) == 0 {
				this.AddScopeQuery(bson.M{"deletedAt": 0})
				return
		}
		if len(status) == 1 {
				var deletedAt = status[0]
				if deletedAt == 1 {
						this.AddScopeQuery(bson.M{"deletedAt": bson.M{"$gt": 0}})
						return
				}
				this.AddScopeQuery(bson.M{"deletedAt": deletedAt})
				return
		}
		this.AddScopeQuery(bson.M{"deletedAt": bson.M{"$in": status}})
		return
}

// 移出软删除条件
func (this *BaseModel) UnUseSoftDel() {
		if this._Scope == nil && len(this._Scope) <= 0 {
				return
		}
		delete(this._Scope, "deletedAt")
}

// 是重复异常
func (this *BaseModel) IsDuplicate(err error) bool {
		return IsDuplicate(err)
}

// 获取锁
func (this *BaseModel) getRedisLocker(name string, duration ...time.Duration) string {
		var (
				err      error
				value    = time.Now().UnixNano()
				cacheIns = this.getLocker()
		)
		if len(duration) == 0 {
				duration = append(duration, time.Minute)
		}
		if cacheIns == nil {
				cacheIns = cache.NewMemoryCache()
		}
		if locker, ok := cacheIns.(cache.Cache); ok {
				// 检查
				if locker.IsExist(name) {
						return ""
				}
				// 上锁
				err = locker.Put(name, value, duration[0])
				if err == nil {
						return name
				}
				// 检查
				v := locker.Get(name)
				if v == nil || v != value {
						return ""
				}
				return name
		}
		return ""
}

// 解锁
func (this *BaseModel) unLocker(name string) bool {
		var err = this.getLocker().Delete(name)
		if err == nil {
				return true
		}
		return false
}

// 高级查询
func (this *BaseModel) Pipe(handler func(pipe *mgo.Pipe) interface{}) interface{} {
		return handler(this.Document().Pipe(nil))
}

// 执行事务
// 示例：
//  ops := []txn.Op{{
//				C:      "accounts",
//				Id:     "aram",
//				Assert: bson.M{"balance": bson.M{"$gte": 100}},
//				Update: M{"$inc": M{"balance": -100}},
//		}, {
//				C:      "accounts",
//				Id:     "ben",
//				Assert: M{"valid": true},
//				Update: M{"$inc": M{"balance": 100}},
//		}}
//	  runner.Run(ops, id, nil)
func (this *BaseModel) Txn(handler func(runner *txn.Runner, txnId bson.ObjectId) error) error {
		var (
				table         = this.Document()
				txnId, runner = this.StartTxn(table)
		)
		defer this.destroy(table)
		return handler(runner, txnId)
}

// 绑定数据对象
func (this *BaseModel) Bind(v interface{}) Model {
		if v == nil {
				return this
		}
		if tab, ok := v.(Model); ok {
				this._Binder = tab
		}
		return this
}

// 提交事务
func (this *BaseModel) Commit(ctx TxnContext) error {
		if ctx.TxnRunner == nil {
				return errors.New("empty txn runner")
		}
		return ctx.TxnRunner.Run(ctx.TxnOps, ctx.TxnId, ctx.TxnResult)
}

// 创建事务
func (this *BaseModel) StartTxn(docs ...*mgo.Collection) (bson.ObjectId, *txn.Runner) {
		if len(docs) == 0 {
				docs = append(docs, this.Document())
		}
		var (
				id     = bson.NewObjectId()
				runner = txn.NewRunner(docs[0])
		)
		return id, runner
}

// 创建索引
func (this *BaseModel) createIndex(handler func(doc *mgo.Collection), force ...bool) {
		if len(force) == 0 || handler == nil {
				return
		}
		if force[0] {
				handler(this.Document())
		}
		return
}

// 日志记录
func (this *BaseModel) logs(msg interface{}) {
		if msg == nil {
				return
		}
		switch msg.(type) {
		case error:
				logs.Error(msg)
		case string:
				logs.Info(msg)
		case map[string]interface{}:
				logs.Debug(fmt.Sprintf("info: %v ", msg))
		default:
				logs.Debug(msg)
		}
}

// 获取锁
func (this *BaseModel) getLocker() cache.Cache {
		var (
				cacheIns = this.GetProfile("redisLocker")
		)
		if cacheIns == nil {
				cacheIns = cache.NewMemoryCache()
				this.SetProfile("redisLocker", cacheIns)
		}
		if locker, ok := cacheIns.(cache.Cache); ok {
				return locker
		}
		return nil
}

func (this *BaseModel) wait() {
		var s = float64(randInt(100, 200))
		time.Sleep(time.Duration(s) * time.Millisecond)
}

func randInt(min, max int) int {
		rand.Seed(time.Now().UnixNano())
		min, max = int(math.Max(float64(min), float64(max))), int(math.Min(float64(min), float64(max)))
		var n = rand.Intn(max) + min
		if n > max {
				return max
		}
		return n
}

func IsNotFound(err error) bool {
		return mgo.ErrNotFound == err
}

func IsErrCursor(err error) bool {
		return mgo.ErrCursor == err
}

func IsDuplicate(err error) bool {
		return mgo.IsDup(err)
}

func GetDate() int64 {
		var (
				bitM, bitD   = "", ""
				year, m, day = time.Now().Local().Date()
		)
		if int(m) < 10 {
				bitM = "0"
		}
		if day < 10 {
				bitD = "0"
		}
		date, _ := strconv.Atoi(fmt.Sprintf("%d%s%d%s%d", year, bitM, int(m), bitD, day))
		return int64(date)
}

func arrayFirst(arr []string, defaults ...string) string {
		if len(defaults) == 0 {
				defaults = append(defaults, "")
		}
		if len(arr) == 0 {
				return defaults[0]
		}
		return arr[0]
}

func ArrayFirst(arr []interface{}) interface{} {
		if len(arr) == 0 {
				return nil
		}
		return arr[0]
}
