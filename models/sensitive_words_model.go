package models

import (
		"crypto/md5"
		"encoding/hex"
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/logs"
		"github.com/globalsign/mgo"
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/transforms"
		"sync"
		"time"
)

// 敏感词记录
type SensitiveWords struct {
		Id            bson.ObjectId `json:"id" bson:"_id"`              // ID
		App           string        `json:"app" bson:"app"`             // 规则名
		Type          string        `json:"type" bson:"type"`           // 分类名
		Word          string        `json:"word" bson:"word"`           // 词
		Hash          string        `json:"hash" bson:"hash"`           // 唯一
		Status        int           `json:"status" bson:"status"`       // 状态
		CreatedAt     time.Time     `json:"createdAt" bson:"createdAt"` // 创建时间
		UpdatedAt     time.Time     `json:"updatedAt" bson:"updatedAt"` // 更新时间
		dataClassImpl `bson:",omitempty"  json:",omitempty"`            // 工具类
}

type SensitiveWordsModel struct {
		BaseModel
}

// DFA 过滤
type (
		Null     struct{}
		dfaLogic struct {
				Set           map[string]Null
				locker        sync.Mutex
				InvalidWord   map[string]Null        // 空
				sensitiveWord map[string]interface{} //无效词汇，不参与敏感词汇判断直接忽略
		}
)

const (
		SensitiveWordsTable      = "sensitive_words"                                                           // 表名
		SensitiveWordsTypeGlobal = "global"                                                                    // 朋友关系
		InvalidWords             = " ,~,!,@,#,$,%,^,&,*,(,),_,-,+,=,?,<,>,.,—,，,。,/,\\,|,《,》,？,;,:,：,',‘,；,“," // 无效字符
)

var (
		_dfaInstance *dfaLogic
		_dfaLocker   = sync.Once{}
)

func NewSensitiveWords() *SensitiveWords {
		var data = new(SensitiveWords)
		data.Init()
		return data
}

func SensitiveWordsModelOf() *SensitiveWordsModel {
		var model = new(SensitiveWordsModel)
		return model.init()
}

func (this *SensitiveWords) data() beego.M {
		return beego.M{
				"id":        this.Id.Hex(),
				"app":       this.App,
				"type":      this.Type,
				"status":    this.Status,
				"word":      this.Word,
				"createdAt": this.CreatedAt.Unix(),
				"updatedAt": this.UpdatedAt.Unix(),
		}
}

func (this *SensitiveWords) getHash() string {
		if this.Hash == "" {
				if this.Word == "" {
						return ""
				}
				var ins = md5.New()
				ins.Write([]byte(this.Word))
				return hex.EncodeToString(ins.Sum(nil))
		}
		return this.Hash
}

func (this *SensitiveWords) Init() {
		this.AddFilters(transforms.FilterEmpty)
		this.SetProvider(DataProvider, this.data)
		this.SetProvider(SaverProvider, this.save)
		this.SetProvider(DefaultProvider, this.defaults)
		this.SetProvider(AttributesProvider, this.setAttributes)
}

func (this *SensitiveWords) defaults() {
		if this.Id == "" {
				this.Id = bson.NewObjectId()
		}
		if this.UpdatedAt.IsZero() {
				this.UpdatedAt = time.Now().Local()
		}
		if this.CreatedAt.IsZero() {
				this.CreatedAt = time.Now().Local()
		}
		if this.Status == 0 {
				this.Status = 1
		}
		if this.App == "" {
				this.App = NewAppInfo().GetAppName()
		}
		if this.Type == "" {
				this.Type = SensitiveWordsTypeGlobal
		}
		if this.Hash == "" && this.Word != "" {
				this.Hash = this.getHash()
		}
}

func (this *SensitiveWords) save() error {
		var (
				model = UserRelationModelOf()
				data  = model.GetByUnique(this.data())
		)
		if data == nil {
				this.InitDefault()
				return model.Add(this)
		}
		return model.Update(bson.M{"_id": data.Id}, this.M(func(m beego.M) beego.M {
				delete(m, "id")
				delete(m, "createdAt")
				m["updatedAt"] = time.Now().Local()
				return m
		}))
}

func (this *SensitiveWords) Set(key string, v interface{}) *SensitiveWords {
		switch key {
		case "id":
				this.SetObjectId(&this.Id, v)
		case "app":
				this.SetString(&this.App, v)
		case "type":
				this.SetString(&this.Type, v)
		case "word":
				this.SetString(&this.Word, v)
		case "status":
				this.SetNumInt(&this.Status, v)
		case "hash":
				this.SetString(&this.Hash, v)
		case "createdAt":
				this.SetTime(&this.CreatedAt, v)
		case "updatedAt":
				this.SetTime(&this.UpdatedAt, v)
		}
		return this
}

func (this *SensitiveWords) setAttributes(data map[string]interface{}, safe ...bool) {
		for key, v := range data {
				if !safe[0] {
						if this.Excludes(key) {
								continue
						}
						if this.IsEmpty(v) {
								continue
						}
				}
				this.Set(key, v)
		}
}

func (this *SensitiveWordsModel) TableName() string {
		return SensitiveWordsTable
}

func (this *SensitiveWordsModel) CreateIndex(force ...bool) {
		this.createIndex(this.getCreateIndexHandler(), force...)
}

func (this *SensitiveWordsModel) getCreateIndexHandler() func(*mgo.Collection) {
		return func(doc *mgo.Collection) {
				this.logs(doc.EnsureIndex(mgo.Index{
						Key:    []string{"app", "hash"},
						Unique: true,
						Sparse: false,
				}))
				this.logs(doc.EnsureIndexKey("word"))
				this.logs(doc.EnsureIndexKey("type", "status"))
		}
}

func (this *SensitiveWordsModel) init() *SensitiveWordsModel {
		this.Bind(this)
		this.Init()
		return this
}

func (this *SensitiveWordsModel) GetByUnique(m beego.M) *SensitiveWords {
		var (
				err   error
				data  = NewSensitiveWords()
				query = bson.M{"app": "", "hash": ""}
		)
		for key := range query {
				v, ok := m[key]
				if !ok {
						return nil
				}
				str, ok := v.(string)
				if ok && str != "" {
						query[key] = str
						continue
				}
				return nil
		}
		err = this.NewQuery(query).One(data)
		if err == nil {
				return data
		}
		return nil
}

func (this *SensitiveWordsModel) Foreach(each func(it *SensitiveWords) bool, limit ...int) {
		var (
				count = this.GetTotal()
				items = make([]*SensitiveWords, 100)
				query = bson.M{"app": NewAppInfo().GetAppName(), "status": 1}
		)
		if len(limit) == 0 {
				limit = append(limit, count)
		}
		var (
				Query = this.NewQuery(query)
				iter  = Query.Limit(limit[0]).Iter()
		)
		items = items[:0]
		// 遍历
		for {
				m := NewSensitiveWords()
				if !iter.Next(&m) {
						break
				}
				items = append(items, m)
		}
		for _, it := range items {
				if !each(it) {
						break
				}
		}
}

func (this *SensitiveWordsModel) Filters(words string) string {
		return GetDfaInstance().ChangeSensitiveWords(words)
}

// 批量添加
func (this *SensitiveWordsModel) Adds(words []string, ty string, app ...string) error {
		if len(app) == 0 {
				app = append(app, NewAppInfo().GetAppName())
		}
		var (
				it      = NewSensitiveWords()
				updates = make([]interface{}, 1)
				items   = make([]interface{}, len(words))
		)
		items = items[:0]
		updates = updates[:0]
		for _, word := range this.uniqueArray(words) {
				data := NewSensitiveWords()
				data.Word = word
				data.Hash = getHash(word)
				data.Type = ty
				data.App = app[0]
				err := this.FindOne(bson.M{"app": data.App, "hash": data.Hash}, it)
				if err == nil {
						if it.Type == ty {
								continue
						}
						data.Type = ty
						data.Status = 1
						data.CreatedAt = it.CreatedAt
						data.UpdatedAt = time.Now().Local()
						data.CreatedAt = it.CreatedAt
						data.Id = it.Id
						updates = append(updates, data)
						continue
				}
				data.defaults()
				items = append(items, data)
		}
		if len(updates) > 0 {
				table := this.Document()
				defer this.destroy(table)
				info, err := table.UpdateAll(bson.M{}, updates)
				if err != nil {
						logs.Error(err)
				}
				if info.Updated <= 0 {
						logs.Info(SensitiveWordsTable + " update failed ")
				}
		}
		if len(items) < 0 {
				return nil
		}
		return this.Inserts(items)
}

// 更新过滤器
func (this *SensitiveWordsModel) updateDfaMapper(words []string) {
		var dfa = GetDfaInstance()
		for _, word := range words {
				dfa.AddSensitiveToMap(map[string]Null{word: {}})
		}
}

func (this *SensitiveWordsModel) uniqueArray(words []string) []string {
		var (
				results []string
				cache   = make(map[string]int)
		)
		for _, it := range words {
				if _, ok := cache[it]; ok {
						continue
				}
				cache[it] = 1
				results = append(results, it)
		}
		return results
}

func (this *SensitiveWordsModel) GetTotal() int {
		return this.Count(bson.M{"status": 1, "app": NewAppInfo().GetAppName()})
}

func getHash(word string) string {
		var ins = md5.New()
		ins.Write([]byte(word))
		return hex.EncodeToString(ins.Sum(nil))
}

func newDfa() {
		_dfaInstance = new(dfaLogic)
		_dfaInstance.init()
}

func GetDfaInstance() *dfaLogic {
		if _dfaInstance == nil {
				_dfaLocker.Do(newDfa)
		}
		return _dfaInstance
}

func (this *dfaLogic) init() *dfaLogic {
		if this.InvalidWord == nil || len(this.InvalidWord) == 0 {
				this.InvalidWord = map[string]Null{}
				this.loadInvalidWord()
		}
		if this.Set == nil || len(this.Set) == 0 {
				this.Set = map[string]Null{}
				this.sensitiveWord = map[string]interface{}{}
				this.loadSensitives()
		}
		this.locker = sync.Mutex{}
		return this
}

func (this *dfaLogic) loadSensitives() {
		SensitiveWordsModelOf().Foreach(func(it *SensitiveWords) bool {
				this.AddSensitiveToMap(map[string]Null{it.Word: {}})
				return true
		})
}

func (this *dfaLogic) loadInvalidWord() *dfaLogic {
		var arr = []rune(InvalidWords)
		for _, ch := range arr {
				this.InvalidWord[string(ch)] = struct{}{}
		}
		return this
}

//生成违禁词集合
func (this *dfaLogic) AddSensitiveToMap(set map[string]Null) {
		this.locker.Lock()
		defer this.locker.Unlock()
		for key := range set {
				str := []rune(key)
				nowMap := this.sensitiveWord
				for i := 0; i < len(str); i++ {
						char := string(str[i])
						_, ok := nowMap[char]
						// 存在
						if ok {
								nowMap = nowMap[char].(map[string]interface{})
						} else {
								//如果该key不存在
								thisMap := make(map[string]interface{})
								thisMap["isEnd"] = false
								nowMap[char] = thisMap
								nowMap = thisMap
						}
						if i == len(str)-1 {
								nowMap["isEnd"] = true
						}
				}
		}
}

//敏感词汇转换为*
func (this *dfaLogic) ChangeSensitiveWords(txt string, sensitives ...map[string]interface{}) (word string) {
		var (
				start  = -1
				tag    = -1
				str    = []rune(txt)
				nowMap map[string]interface{}
		)
		if len(sensitives) <= 0 {
				sensitives = append(sensitives, this.sensitiveWord)
		}
		nowMap = sensitives[0]
		for i := 0; i < len(str); i++ {
				ch := string(str[i])
				_, ok := this.InvalidWord[ch]
				if ok || string(str[i]) == "," {
						continue
				}
				value, ok := nowMap[ch]
				if ok {
						thisMap := value.(map[string]interface{})
						tag++
						if tag == 0 {
								start = i
						}
						isEnd, _ := thisMap["isEnd"].(bool)
						if !isEnd {
								value = nowMap[ch]
								nowMap = value.(map[string]interface{})
								continue
						}
						for y := start; y < i+1; y++ {
								str[y] = 42
						}
						nowMap = sensitives[0]
						start = -1
						tag = -1
						continue
				}
				if start != -1 {
						i = start + 1
				}
				nowMap = sensitives[0]
				start = -1
				tag = -1
		}
		return string(str)
}
