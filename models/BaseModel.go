package models

import (
		"github.com/astaxie/beego/logs"
		"github.com/globalsign/mgo"
		"sync"
)

type Model interface {
		Error() error
		GetConn(name ...string) *mgo.Session
		SetProfile(key string, v interface{}) Model
		GetProfile(key string, defaults ...interface{}) interface{}
}

type BaseModel struct {
		Err         error
		Connections map[string]*mgo.Session
		lock        sync.Mutex
		_Profiles   map[string]interface{}
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
						_baseProfiles = new(map[string]interface{})
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

func (this *BaseModel) Init() {
		if this.Connections == nil {
				this.Connections = map[string]*mgo.Session{}
		}
		if this._Profiles == nil {
				this._Profiles = *GetModelProfiles()
		}
}

func (this *BaseModel) getConnUrl(name string) string {
		return ""
}

func (this *BaseModel) GetConn(name ...string) *mgo.Session {
		if len(name) == 0 {
				name = append(name, "default")
		}
		this.lock.Lock()
		if conn, ok := this.Connections[name[0]]; ok {
				return conn
		}
		this.lock.Unlock()
		return this.conn(name[0])
}

func (this *BaseModel) conn(name string) *mgo.Session {
		this.lock.Lock()
		defer this.lock.Unlock()
		if sess, err := mgo.Dial(this.getConnUrl(name)); err == nil {
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
