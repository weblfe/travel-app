package kernel

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/globalsign/mgo"
	"os"
	"runtime"
	"sync"
)

type DbMgr struct {
	dbs    map[string]*mgo.Session
	locker sync.RWMutex
}

var (
	dbMgr = newDbMgr()
)

func newDbMgr() *DbMgr {
	var mgr = new(DbMgr)
	mgr.dbs = make(map[string]*mgo.Session)
	mgr.locker = sync.RWMutex{}
	mgr.init()
	return mgr
}

func GetDbMgr() *DbMgr {
	return dbMgr
}

func (dbMgr *DbMgr) init() {
	runtime.SetFinalizer(dbMgr, (*DbMgr).destroy)
}

func (dbMgr *DbMgr) GetDb(name ...string) *mgo.Session {
	name = append(name, "")
	var (
		url  = dbMgr.getUrl(name[0])
		hash = hashUrl(url)
	)
	dbMgr.locker.Lock()
	if v, ok := dbMgr.dbs[hash]; ok {
		dbMgr.locker.Unlock()
		return v
	}
	dbMgr.locker.Unlock()
	var db, err = dbMgr.openDb(url)
	if err != nil {
		panic(err)
	}
	return db
}

func (dbMgr *DbMgr) getUrl(name string) string {
	var dbUrl = envOr(envPrefix("DB_CONN_URL", name), "127.0.0.1:27017")
	return dbUrl
}

func (dbMgr *DbMgr) openDb(url string) (*mgo.Session, error) {
	var (
		hash         = hashUrl(url)
		session, err = mgo.Dial(url)
	)
	if err != nil {
		fmt.Println("OpenDb:", err)
		return nil, err
	}
	if session == nil {
		fmt.Println("nil session")
		return nil, errors.New("nil session")
	}
	dbMgr.append(hash, session)
	return session, nil
}

func (dbMgr *DbMgr) append(key string, db *mgo.Session) bool {
	dbMgr.locker.Lock()
	defer dbMgr.locker.Unlock()
	if _, ok := dbMgr.dbs[key]; ok {
		return ok
	}
	dbMgr.dbs[key] = db
	return true
}

func (dbMgr *DbMgr) destroy() {
	var keys []string
	defer runtime.SetFinalizer(dbMgr, nil)
	for k, v := range dbMgr.dbs {
		v.Close()
		fmt.Println(fmt.Sprintf("key: %s,db.close", k))
		keys = append(keys, k)
	}
	for _, v := range keys {
		delete(dbMgr.dbs, v)
	}
}

func hashUrl(url string) string {
	var (
		m      = md5.New()
		_, err = m.Write([]byte(url))
	)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return fmt.Sprintf("%x", m.Sum(nil))
}

func envOr(key string, value string) string {
	var v = os.Getenv(key)
	if v == "" {
		return value
	}
	return v
}

func envPrefix(key string, prefix string) string {
	if prefix != "" {
		return fmt.Sprintf("%s_%s", prefix, key)
	}
	return key
}
