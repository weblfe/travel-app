package repositories

import (
		"errors"
		"fmt"
		"github.com/weblfe/travel-app/models"
		"strconv"
		"strings"
		"sync"
		"unicode"
)

type (
	queryTypeRepository struct {
		typesDefinedData *TypeDefinedMgr
		once             sync.Once
	}

	TypeAlias      string
	typeDefinedMap map[TypeAlias][]int

	TypeDefinedMgr struct {
		data *typeDefinedMap
		all  map[int][]TypeAlias
	}
)

var (
	queryTypesParser = NewQueryTypeRepository()
)

func GetQueryTypesParser() *queryTypeRepository {
	return queryTypesParser
}

func (t TypeAlias) String() string {
	return string(t)
}

func (tM typeDefinedMap) Get(alias string) ([]int, bool) {
	if v, ok := tM[TypeAlias(alias)]; ok {
		return v, true
	}

	return []int{}, false
}

func (tM typeDefinedMap) Len() int {
	return len(tM)
}

func (tM typeDefinedMap) Exist(alias TypeAlias) bool {
	if _, ok := tM[alias]; ok {
		return true
	}
	return false
}

func NewTypeDefinedMgr() *TypeDefinedMgr {
	var (
		mgr  = new(TypeDefinedMgr)
		data = make(typeDefinedMap)
	)
	mgr.data = &data
	mgr.all = make(map[int][]TypeAlias)
	return mgr
}

func NewQueryTypeRepository() *queryTypeRepository {
	var rep = new(queryTypeRepository)
	rep.init()
	return rep
}

func (this *TypeDefinedMgr) Load(data typeDefinedMap) *TypeDefinedMgr {
	for k, v := range data {
		if len(v) <= 0 {
			continue
		}
		(*this.data)[k] = v
		for _, t := range v {
			this.all[t] = append(this.all[t], k)
		}
	}
	return this
}

func (this *TypeDefinedMgr) Get(k string) ([]int, bool) {
	if this.data == nil {
		return nil, false
	}
	if v, ok := this.data.Get(k); ok {
		return v, true
	}
	// 是否数字类型
	if !this.isDigit(k) {
		return nil, false
	}
	// typeInt
	ns, _ := strconv.Atoi(k)
	if _, ok := this.all[ns]; ok {
		return []int{ns}, true
	}
	return nil, false
}

func (TypeDefinedMgr) isDigit(value string) bool {
	for _, n := range []rune(value) {
		if !unicode.IsDigit(n) {
			return false
		}
	}
	return true
}

func (this *queryTypeRepository) init() {
	this.once = sync.Once{}
	this.typesDefinedData = NewTypeDefinedMgr()
	this.Boot()
}

func (this *queryTypeRepository) Boot() {
	//@todo 以后实现多类型注册
	this.once.Do(func() {
		this.typesDefinedData.Load(this.defaultTypes())
	})
}

func (this *queryTypeRepository) defaultTypes() typeDefinedMap {
	return map[TypeAlias][]int{
		models.ImageTypeCode:    {models.ImageType},
		models.VideoTypeCode:    {models.VideoType},
		models.ContentTypeCode:  {models.ContentType},
		models.StrategyTypeCode: {models.StrategyType},
		models.PostTypeCode:     {models.PostType},
		`vip_content`:           {models.PostType},
		`user_content`:          {models.ImageType, models.VideoType, models.ContentType},
		`image_content`:         {models.ImageType, models.StrategyType, models.PostType},
		`all`:                   {models.ImageType, models.VideoType, models.ContentType, models.StrategyType, models.PostType},
	}
}

func (this *queryTypeRepository) Parse(typeQuery string) ([]int, error) {
	if typeQuery == "" {
		return nil, errors.New(`empty query`)
	}
	var (
		types      []int
		typeFilter =make(map[int]interface{})
	)
	for _, typeName := range strings.Split(typeQuery, ",") {
		var (
			ok      bool
			typeArr []int
		)
		if typeArr, ok = this.typesDefinedData.Get(typeName); !ok {
			continue
		}
		for _, v := range typeArr {
			if _, ok = typeFilter[v]; ok {
				continue
			}
			typeFilter[v] = nil
			types = append(types, v)
		}
	}
	var size = len(types)
	if size <= 0 {
		return nil, errors.New(fmt.Sprintf(`parse undefined type %s`, typeQuery))
	}
	return types, nil
}
