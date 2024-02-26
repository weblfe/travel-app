package plugins

import (
		"context"
		"encoding/json"
		"fmt"
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/logs"
		"github.com/qiniu/api.v7/v7/auth/qbox"
		"github.com/qiniu/api.v7/v7/storage"
		"io"
		url2 "net/url"
		"os"
		"time"
)

type OssPlugin struct {
		token         *TokenEntry
		configuration beego.M
		Policy        *storage.PutPolicy
}

type Oss interface {
		CreateUploader(params OssParams, configures ...func(*storage.Config)) func(
				context.Context, ...func(*storage.PutExtra)) (interface{}, error)
}

type TokenEntry struct {
		ExpireAt int64
		Value    string
}

const (
		OssPluginName     = "oss"
		QinNiuPrefixKey   = "QINNIU_PREFIX"
		QinNiuAccessKey   = "AK"
		QinNiuSecretKey   = "SK"
		QinNiuBucketImg   = "IMG"
		QinNiuBucketVideo = "VIDEO"
		QinNiuReturnBody  = "RETURN_BODY"
		QinNiuCallbackUrl = "CALLBACK_URL"
		OssProviderQinNiu = "QinNiu"
		OssDefaultExpires = 1800
)

var (
		_QinNiuPrefix   string
		_QinNiuInstance *OssPlugin
		_PropertiesKeys = []string{
				"CDN_DOMAIN", "CDN_IMG_DOMAIN",
				"CDN_VIDEO_DOMAIN", "AK",
				"SK", "IMG_BUCKET", "VIDEO_BUCKET",
		}
		_QinNiuProperties = make(map[string]string)
)

func OSS() *OssPlugin {
		var locker = getLock(OssPluginName)
		locker.Do(newQinNiClient)
		return _QinNiuInstance
}

func newQinNiClient() {
		_QinNiuInstance = new(OssPlugin)
		_QinNiuInstance.init()
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

// Lock 获取锁
func (this *OssPlugin) Lock() {

}

// UnLock 解锁
func (this *OssPlugin) UnLock() {

}

func (this *OssPlugin) GetAccessKey() string {
		return GetQinNiuProperty(QinNiuAccessKey)
}

func (this *OssPlugin) GetSecretKey() string {
		return GetQinNiuProperty(QinNiuSecretKey)
}

func (this *OssPlugin) getMac() *qbox.Mac {
		return qbox.NewMac(this.GetAccessKey(), this.GetAccessKey())
}

func (this *OssPlugin) CreatePolicy(jsonData ...string) *OssPlugin {
		if len(jsonData) == 0 {
				this.Policy = &storage.PutPolicy{
						Scope: this.GetBucket(QinNiuBucketImg),
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

func (this *OssPlugin) GetBucket(typ string) string {
		var key = typ + "_BUCKET"
		return GetQinNiuProperty(key)
}

func (this *OssPlugin) Register() {
		Plugin(this.PluginName(), this)
}

func (this *OssPlugin) PluginName() string {
		return OssPluginName
}

// CreateUploader 构建上传函数
func (this *OssPlugin) CreateUploader(params *OssParams, configures ...func(*storage.Config)) func(
		context.Context, ...func(*storage.PutExtra)) (interface{}, error) {

		var (
				key       = params.Key
				bucket    = params.Bucket
				Expires   = params.Expires
				accessKey = this.GetAccessKey()
				secretKey = this.GetSecretKey()
		)

		if bucket == "" && params.TypeName != "" {
				bucket = this.GetBucket(params.TypeName)
		}
		var putPolicy storage.PutPolicy
		if params.Provider == "" {
				params.Provider = OssProviderQinNiu
		}
		if params.PutPolicy.Scope == "" {
				params.PutPolicy.Scope = bucket
		}
		if params.PutPolicy.Expires == 0 {
				params.PutPolicy.Expires = OssDefaultExpires
		}
		if params.PutPolicy != nil {
				putPolicy = *params.PutPolicy
		} else {
				putPolicy = storage.PutPolicy{
						Scope:   bucket,
						Expires: Expires,
				}
		}
		if putPolicy.CallbackURL == "" {
				putPolicy.CallbackURL = GetQinNiuProperty("CALLBACK_URL")
		}
		if putPolicy.ReturnBody == "" {
				putPolicy.ReturnBody = GetQinNiuProperty("RETURN_BODY")
		}
		var (
				cfg = params.Storage
				mac = qbox.NewMac(accessKey, secretKey)
				// @todo cache
				upToken = putPolicy.UploadToken(mac)
		)
		if cfg.Zone == nil {
				cfg.Zone = &storage.ZoneHuanan
		}
		// 最近配置
		if len(configures) > 0 {
				for _, fn := range configures {
						fn(&cfg)
				}
		}

		var (
				ret          = params.Result
				formUploader = storage.NewFormUploader(&cfg)
		)
		if ret == nil {
				ret = &map[string]interface{}{}
				params.Result = ret
		}
		var putExtra *storage.PutExtra

		// 可选配置
		if params.Extras != nil {
				putExtra = params.Extras
		} else {
				putExtra = &storage.PutExtra{
						Params: map[string]string{
								"x:app":       os.Getenv("APP_NAME"),
								"x:timestamp": fmt.Sprintf("%v", time.Now().Unix()),
						},
				}
		}

		return func(ctx context.Context, extrasHandler ...func(*storage.PutExtra)) (interface{}, error) {
				if len(extrasHandler) > 0 {
						for _, fn := range extrasHandler {
								fn(putExtra)
						}
				}
				//	fmt.Println(params.Result)
				//	fmt.Println(params.PutPolicy.Scope)
				// 	fmt.Println(params.PutPolicy.ReturnBody)
				if params.File == "" && params.Reader != nil {
						return ret, formUploader.Put(ctx, ret, upToken, key, params.Reader, params.Size, putExtra)
				}
				return ret, formUploader.PutFile(ctx, ret, upToken, key, params.File, putExtra)
		}

}

func GetQinNiuProperties() map[string]string {
		if len(_QinNiuProperties) > 0 {
				return _QinNiuProperties
		}
		if _QinNiuPrefix == "" {
				_QinNiuPrefix = os.Getenv(QinNiuPrefixKey)
				if _QinNiuPrefix == "" {
						_QinNiuPrefix = "QINNIU_"
				}
		}
		for _, key := range GetQinNiuPropertiesKeys() {
				property := _QinNiuPrefix + key
				v := os.Getenv(property)
				_QinNiuProperties[property] = v
		}
		return _QinNiuProperties
}

func GetQinNiuPropertiesKeys() []string {
		return _PropertiesKeys
}

func GetQinNiuProperty(key string, defaults ...string) string {
		var properties = GetQinNiuProperties()
		if v, ok := properties[_QinNiuPrefix+key]; ok {
				return v
		}
		return arrayFirst(defaults)
}

func arrayFirst(arr []string) string {
		if len(arr) == 0 {
				return ""
		}
		return arr[0]
}

func GetOSS() *OssPlugin {
		return OSS()
}

func AppendCertificate(Url string, expire int64) string {
		if expire == 0 {
				expire = time.Now().Add(30 * time.Minute).Unix()
		}
		var info, err = url2.Parse(Url)
		if err != nil {
				logs.Error(err)
				return Url
		}
		var (
				ak          = GetOSS().GetAccessKey()
				sk          = GetOSS().GetSecretKey()
				mac         = qbox.NewMac(ak, sk)
				domain, key = info.Scheme + "://" + info.Host, info.Path
		)
		if key[0] == '/' {
				key = key[1:]
		}
		return storage.MakePrivateURL(mac, domain, key, expire)
}

func GetOssAccessUrl(key string, ossName string, bucket string) string {
		if ossName != OssProviderQinNiu {
				return ""
		}
		var ty = getBucketTypeByBucket(bucket)
		if ty == "" {
				return ""
		}
		var host = GetQinNiuProperty("CDN_" + ty + "_DOMAIN")
		if host == "" {
				return ""
		}
		if key[0] == '/' {
				return host + key
		}
		return host + "/" + key
}

func getBucketTypeByBucket(bucket string) string {
		if GetQinNiuProperty(QinNiuBucketImg+"_BUCKET") == bucket {
				return QinNiuBucketImg
		}
		if GetQinNiuProperty(QinNiuBucketVideo+"_BUCKET") == bucket {
				return QinNiuBucketVideo
		}
		return ""
}

type OssParams struct {
		Provider  string             `json:"provider"`          // 服务提供
		TypeName  string             `json:"type"`              // 图片 ｜ 视频类型 IMG , VIDEO 按类型自动获取 bucket
		Token     string             `json:"token,omitempty"`   // 上传令牌
		Bucket    string             `json:"bucket,omitempty"`  // 存储桶 ，有对应类型时可选
		Storage   storage.Config     `json:"storage"`           // 存储配置 [ 区域，https, cdn加速域名]
		Result    interface{}        `json:"result"`            // 上传结果对象
		Reader    io.Reader          `json:",omitempty"`        // 上传io流
		Size      int64              `json:"size,omitempty"`    // 流大小
		File      string             `json:"file,omitempty"`    // 上传本地文件地址
		Key       string             `json:"key"`               // 上传文件保存名字
		Expires   uint64             `json:"expires,omitempty"` // token 有效时间
		Extras    *storage.PutExtra  `json:"extras"`            // 上传扩展信息
		PutPolicy *storage.PutPolicy `json:"put_policy"`        // 上传策略
}
