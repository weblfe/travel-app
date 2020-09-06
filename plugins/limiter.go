package plugins

import (
		"context"
		"github.com/astaxie/beego/cache"
		"github.com/astaxie/beego/logs"
		"os"
		"time"
)

const (
		LimiterPlugin      = "limiter"
		LimiterTokenPolicy = "token"
		PolicyKey          = "policy"
		GlobalPolicyEnvKey = "LIMIT_POLICY"
		TokenCtxValueKey   = "token"
		MacCtxValueKey     = "mac"
)

type LimitResult struct {
		Ok  bool   `json:"ok"`
		Msg string `json:"msg"`
}

type ContextLimit interface {
		GetPolicy() string
		Ctx() context.Context
		SetPolicy(name string) ContextLimit
		SetValue(name string, value interface{}) ContextLimit
}

type contextLimitImpl struct {
		ctx     context.Context
		cancel  context.CancelFunc
		timeout time.Duration
}

// 限流器
type Limiter interface {
		Get(name string, ctx ContextLimit) func() *LimitResult
		New(ContextLimit) func(ctx ContextLimit) *LimitResult
		SetProvider(name string, provider func(ctx ContextLimit) *LimitResult)
		PluginInterface
}

type limiterImpl struct {
		Name      string `json:"name"`
		Providers map[string]func(ctx ContextLimit) *LimitResult
}

var (
		limiterInstance *limiterImpl
)

func (this *limiterImpl) New(ctx ContextLimit) func(ctx ContextLimit) *LimitResult {
		var name = ctx.GetPolicy()
		if name == "" {
				return this.defaults()
		}
		handler, ok := this.Providers[name]
		if ok && handler != nil {
				return handler
		}
		return this.defaults()
}

// 空策略限流
func (this *limiterImpl) defaults() func(ctx ContextLimit) *LimitResult {
		return func(ctx ContextLimit) *LimitResult {
				return &LimitResult{
						true,
						"pass",
				}
		}
}

// 限流策略注册器
func (this *limiterImpl) SetProvider(name string, provider func(ctx ContextLimit) *LimitResult) {
		if _, ok := this.Providers[name]; ok {
				return
		}
		this.Providers[name] = provider
		return
}

// 获取限流处理
func (this *limiterImpl) Get(name string, ctx ContextLimit) func() *LimitResult {
		if _, ok := this.Providers[name]; !ok {
				return func() *LimitResult {
						return this.defaults()(ctx)
				}
		}
		return func() *LimitResult {
				var handler = this.Providers[name]
				return handler(ctx)
		}
}

func (this *limiterImpl) Register() {
		Plugin(this.PluginName(), this)
		this.Boot()
}

func (this *limiterImpl) PluginName() string {
		return this.Name
}

func (this *limiterImpl) Boot() {
		this.SetProvider(LimiterTokenPolicy, NewTokenLimiterProvider(this.getCache(), 100, 1*time.Second).Handler)
}

func (this *limiterImpl) getCache() cache.Cache {
		var cacheInstance, _ = cache.NewCache("memory", `{"interval":20}`)
		return cacheInstance
}

func GetLimiter() Limiter {
		if limiterInstance == nil {
				var locker = getLock(LimiterPlugin)
				locker.Do(newLimiter)
		}
		return limiterInstance
}

func newLimiter() {
		limiterInstance = new(limiterImpl)
		limiterInstance.Name = LimiterPlugin
		limiterInstance.Providers = make(map[string]func(ctx ContextLimit) *LimitResult, 10)
}

func NewContextLimit() ContextLimit {
		return newCtx()
}

func newCtx() ContextLimit {
		var impl = new(contextLimitImpl)
		impl.timeout = 60 * time.Second
		impl.ctx, impl.cancel = context.WithTimeout(context.Background(), impl.timeout)
		return impl
}

func (this *contextLimitImpl) Ctx() context.Context {
		return this.ctx
}

func (this *contextLimitImpl) GetPolicy() string {
		var v = this.ctx.Value(PolicyKey)
		if v == nil || v == "" {
				// 全局限流规则
				return os.Getenv(GlobalPolicyEnvKey)
		}
		return v.(string)
}

func (this *contextLimitImpl) SetPolicy(name string) ContextLimit {
		this.ctx = context.WithValue(this.ctx, PolicyKey, name)
		return this
}

func (this *contextLimitImpl) SetValue(name string, value interface{}) ContextLimit {
		this.ctx = context.WithValue(this.ctx, name, value)
		return this
}

type TokenLimiterProvider interface {
		Handler(ctx ContextLimit) *LimitResult
}

type tokenLimiterProviderImpl struct {
		storage      cache.Cache
		maxTimes     int
		timeInterval time.Duration
}

func NewTokenLimiterProvider(storage cache.Cache, max int, interval time.Duration) TokenLimiterProvider {
		var provider = new(tokenLimiterProviderImpl)
		provider.maxTimes = max
		provider.storage = storage
		provider.timeInterval = interval
		return provider
}

func (this *tokenLimiterProviderImpl) Handler(ctx ContextLimit) *LimitResult {
		var (
				ctxObj = ctx.Ctx()
				token  = ctxObj.Value(TokenCtxValueKey)
		)
		// 访问令牌
		if token != nil && token != "" {
				return this.limitByToken(token.(string))
		}
		// 临时身份机器码
		var mac = ctxObj.Value(MacCtxValueKey)
		// 机器人访问限制
		if mac == nil || mac == "" {
				// @todo robot
				return &LimitResult{
						Ok:  true,
						Msg: "robot limit!",
				}
		}
		return this.limitByToken(mac.(string))
}

func (this *tokenLimiterProviderImpl) limitByToken(token string) *LimitResult {
		if token == "" {
				return &LimitResult{
						true, "pass empty token",
				}
		}
		var (
				err          error
				timestampNow = time.Now().UnixNano()
				keyLastAt    = this.keyLastAt(token)
				keyTime      = this.keyTimes(token)
		)
		var (
				pass     = true
				lastTime = this.storage.Get(keyLastAt)
				times    = this.storage.Get(keyTime)
		)
		// 计数器 重置,过期, 第一次
		if lastTime == nil {
				err = this.storage.Put(keyTime, 0, time.Minute)
				this.logs(err)
		}
		// 访问时间间隔检查
		if lastTime != nil && this.timeAccessLimit(timestampNow, lastTime.(int64)) {
				pass = false
		}
		// 访问次数检查
		if pass && times != nil && this.accessTimesLimit(times.(int)) {
				pass = false
		}
		// 更新范围次数
		if times != nil {
				err = this.storage.Incr(keyTime)

		} else {
				err = this.storage.Put(keyTime, 1, time.Minute)
		}
		this.logs(err)
		// 更新访问时间
		err = this.storage.Put(keyLastAt, timestampNow, time.Minute)
		this.logs(err)
		if pass {
				return &LimitResult{
						pass, "pass token",
				}
		}
		return &LimitResult{
				pass, "assess token frequently!",
		}
}

// 访问间隔判断
func (this *tokenLimiterProviderImpl) timeAccessLimit(now, last int64) bool {
		var long = int64(this.timeInterval)
		if now > last && now-last >= long {
				return false
		}
		logs.Info("limit timer ")
		return true
}

// 访问次数
func (this *tokenLimiterProviderImpl) accessTimesLimit(times int) bool {
		if this.maxTimes >= times {
				return false
		}
		logs.Info("limit times ")
		return true
}

func (this *tokenLimiterProviderImpl) keyLastAt(token string) string {
		return "token_access_lastAt_" + token
}

func (this *tokenLimiterProviderImpl) keyTimes(token string) string {
		return "token_access_unit_times_" + token
}

func (this *tokenLimiterProviderImpl) logs(err error) {
		if err == nil {
				return
		}
		logs.Error(err)
}
