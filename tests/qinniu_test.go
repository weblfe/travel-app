package test

import (
		cc "context"
		"fmt"
		"github.com/astaxie/beego/logs"
		"github.com/joho/godotenv"
		"github.com/qiniu/api.v7/v7/storage"
		. "github.com/smartystreets/goconvey/convey"
		"github.com/weblfe/travel-app/plugins"
		"os"
		"path/filepath"
		"testing"
)

func init() {
		var pwd, err = os.Getwd()
		if err != nil {
				logs.Error(err)
				return
		}
		file := filepath.Join(pwd, ".env")
		err = godotenv.Load(file)
		if err != nil {
				logs.Error(err)
				return
		}
}

func TestGetQinNiuProperties(t *testing.T) {
		Convey("test QinNiu", t, func() {
				var p = plugins.GetQinNiuProperties()
				So(len(p) > 0, ShouldBeTrue)
		})
}

func TestGetOSS(t *testing.T) {
		Convey("test QinNiu OSS", t, func() {
				var p = plugins.GetOSS()
				So(p != nil, ShouldBeTrue)
		})
}

func TestOSSPlugin_Upload(t *testing.T) {
		var (
				pwd, _ = os.Getwd()
				key    = "1.jpg"
				file   = filepath.Join(pwd, "static/storage/2020-7-20/1.jpg")
		)
		fmt.Println(file)
		Convey("test QinNiu OSS", t, func() {
				var uploader = plugins.GetOSS().CreateUploader(&plugins.OssParams{
						TypeName: plugins.QinNiuBucketImg,
						File:     file,
						Key:      key,
				})
				So(uploader != nil, ShouldBeTrue)
				var ctx = cc.Background()
				var res, err = uploader(ctx, func(extra *storage.PutExtra) {
						fmt.Println(extra)
				})
				fmt.Println(res)
				So(err == nil, ShouldBeTrue)
				So(res != nil, ShouldBeTrue)
		})
}
