package test

import (
		"fmt"
		"github.com/weblfe/travel-app/plugins"
		"testing"
		"time"
)

func TestGetConfigureCentreRepositoryMangerInstance(t *testing.T) {
		var (
				err    error
				data   map[string]string
				manger = plugins.GetConfigureCentreRepositoryMangerInstance()
		)
		manger.InitDef()
		var provider = manger.Get(plugins.EtcdProvider)
		err = provider.Put("test", time.Now().String())
		if err != nil {
				t.Error(err)
		}
		data, err = provider.Get("test")
		if err != nil {
				t.Error(err)
		}
		fmt.Printf("%v\n", data)
}
