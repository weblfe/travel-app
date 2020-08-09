package plugins

import (
		"encoding/json"
		"fmt"
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/logs"
		"github.com/qiniu/api.v7/v7/auth/qbox"
		"github.com/qiniu/api.v7/v7/storage"
		"time"
)

type OssPlugin struct {
		token         *TokenEntry
		configuration beego.M
		Policy        *storage.PutPolicy
}

type TokenEntry struct {
		ExpireAt int64
		Value    string
}

const (
		OssPluginName = "oss"
)

func newQinNiClient() {
		var accessKey, secretKey = "", ""
		bucket := "your bucket name"
		putPolicy := storage.PutPolicy{
				Scope: bucket,
		}
		mac := qbox.NewMac(accessKey, secretKey)
		token := putPolicy.UploadToken(mac)
		// 1小时
		fmt.Println(token)
}

func (this *OssPlugin) init() {
		this.token = new(TokenEntry)
		this.token.Value = ""
		this.token.ExpireAt = 0
		this.configuration = beego.M{}
}

// 授权
func (this *OssPlugin) Auth() error {
		this.getToken()
		return nil
}

// 获取token
func (this *OssPlugin) getToken() string {
		if this.token.ExpireAt > time.Now().Unix() {
				return this.token.Value
		}
		this.Lock()
		var token = this.getTokenSync()
		if token == nil {
				this.getMac()
		}

		return ""
}

func (this *OssPlugin) saveToken(token string) error {

		return nil
}

// 获取同步获取token
func (this *OssPlugin) getTokenSync() *TokenEntry {
		return &TokenEntry{

		}
}

func (this *OssPlugin) getInfoKey() string {
		return fmt.Sprintf("%s|%s|%v", this.GetAccessKey(), this.GetAccessKey(), this.GetPolicy())
}

// 获取锁
func (this *OssPlugin) Lock() {

}

// 解锁
func (this *OssPlugin) UnLock() {

}

func (this *OssPlugin) GetAccessKey() string {
		return ""
}

func (this *OssPlugin) GetSecretKey() string {
		return ""
}

func (this *OssPlugin) getMac() *qbox.Mac {
		return qbox.NewMac(this.GetAccessKey(), this.GetAccessKey())
}

func (this *OssPlugin) CreatePolicy(jsonData ...string) *OssPlugin {
		if len(jsonData) == 0 {
				this.Policy = &storage.PutPolicy{
						Scope: this.GetBucket(),
				}
				return this
		}
		this.Policy = &storage.PutPolicy{}
		var err = json.Unmarshal([]byte(jsonData[0]), this.Policy)
		if err != nil {
				logs.Error(err)
		}
		return this
}

func (this *OssPlugin) GetPolicy() *storage.PutPolicy {
		if this.Policy == nil {
				this.CreatePolicy()
		}
		return this.Policy
}

func (this *OssPlugin) GetBucket() string {
		return ""
}

func (this *OssPlugin) Register() {
		Plugin(this.PluginName(), this)
}

func (this *OssPlugin) PluginName() string {
		return OssPluginName
}
