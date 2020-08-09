package models

import (
		"github.com/astaxie/beego"
		"strings"
)

type StrArrayEntry struct {
		Key   string
		Items []string
}

func NewStrArrayEntry() *StrArrayEntry {
		var entry = new(StrArrayEntry)
		entry.Init()
		return entry
}

func (this *StrArrayEntry) Init() {
		this.Items = make([]string, 2)
		this.Items = this.Items[:0]
}

func (this *StrArrayEntry) Cap() int {
		return cap(this.Items)
}

func (this *StrArrayEntry) Len() int {
		return len(this.Items)
}

func (this *StrArrayEntry) Count() int {
		return this.Len()
}

func (this *StrArrayEntry) Values() []string {
		return this.Items
}

func (this *StrArrayEntry) Included(v string, fold ...bool) bool {
		return StrArray(this.Items).Included(v, fold...)
}

func (this *StrArrayEntry) Search(v string, fold ...bool) int {
		return StrArray(this.Items).Search(v, fold...)
}

func (this *StrArrayEntry) Joins(joins ...[]string) *StrArrayEntry {
		for _, arr := range joins {
				this.Items = append(this.Items, arr...)
		}
		return this
}

func (this *StrArrayEntry) GetKey() string {
		return this.Key
}

func (this *StrArrayEntry) SetKey(key string) *StrArrayEntry {
		this.Key = key
		return this
}

func (this *StrArrayEntry) Foreach(each func(i int, value string) bool) {
		if this.Len() <= 0 {
				return
		}
		for i, v := range this.Items {
				if !each(i, v) {
						return
				}
		}
}

func (this *StrArrayEntry) Copy() *StrArrayEntry {
		var entry = NewStrArrayEntry()
		entry.Key = this.Key
		for _, v := range this.Items {
				entry.Items = append(entry.Items, v)
		}
		return entry
}

func (this *StrArrayEntry) Push(i int, value string) *StrArrayEntry {
		if this.Len() > i {
				this.Items[i] = value
		}
		return this
}

func (this *StrArrayEntry) Append(value string) *StrArrayEntry {
		this.Items = append(this.Items, value)
		return this
}

func (this *StrArrayEntry) Pop() string {
		var size = this.Len()
		if size <= 0 {
				return ""
		}
		var (
				end   = this.Len() - 2
				value = this.Items[this.Len()-1]
		)
		if end < 0 {
				end = 0
		}
		this.Items = this.Items[:end]
		return value
}

func (this *StrArrayEntry) ToMapper() beego.M {
		return beego.M{
				this.Key: this.Items,
		}
}

func (this StrArray) Search(v string, fold ...bool) int {
		if len(fold) == 0 {
				fold = append(fold, false)
		}
		for index, it := range this {
				if it == v {
						return index
				}
				if fold[0] && strings.EqualFold(it, v) {
						return index
				}
		}
		return -1
}

func (this StrArray) Included(v string, fold ...bool) bool {
		if len(fold) == 0 {
				fold = append(fold, false)
		}
		for _, it := range this {
				if it == v {
						return true
				}
				if fold[0] && strings.EqualFold(it, v) {
						return true
				}
		}
		return false
}

func (this StrArray) Foreach(each func(i int, value string) bool) {
		for i, v := range this {
				if !each(i, v) {
						break
				}
		}
}
