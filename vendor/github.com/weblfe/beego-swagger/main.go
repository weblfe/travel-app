package beego_swagger

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/swag"
	"html/template"
	"net/http"
	"path"
	"sort"
	"strings"
	"sync"
)

const (
	defaultDocURL   = "doc.json"
	defaultIndex    = "index.html"
	routerKey       = ":splat"
	contentType     = "Content-Type"
	contentTypeHtml = "text/html; charset=utf-8"
	contentTypeJson = "application/json; charset=utf-8"
)

// Handler default
var Handler = New()

// Config stores SwaggerUI configuration variables
type Config struct {
	// Enable deep linking for tags and operations, default is true
	DeepLinking bool

	// Controls the default expansion setting for the operations and tags.
	// 'list' (default, expands only the tags),
	// 'full' (expands the tags and operations),
	// 'none' (expands nothing)
	DocExpansion string

	// Configuration information for OAuth2, optional if using OAuth2
	OAuth *OAuthConfig

	// Custom OAuth redirect URL
	OAuth2RedirectUrl string

	// URL pointing to API definition
	URL string
}

type OAuthConfig struct {
	// application name, displayed in authorization popup
	AppName string

	// ID of the client sent to the OAuth2 Provider, default is clientId
	ClientId string
}

// New returns custom handler
func New(config ...Config) beego.FilterFunc {
	cfg := Config{
		DeepLinking:  true,
		DocExpansion: "list",
	}

	if len(config) > 0 {
		cfg = config[0]
	}

	index, err := template.New("swagger_index.html").Parse(indexTmpl)
	if err != nil {
		panic("swagger: could not parse index template")
	}

	var (
		prefix string
		once   sync.Once
		fs     = http.FileServer(swaggerFiles.HTTP)
	)

	return func(c *context.Context) {
		// Set prefix
		once.Do(func() {
			prefix = parsePrefix(c.Input.Param(routerKey))
			// Set doc url
			if len(cfg.URL) == 0 {
				cfg.URL = path.Join(prefix, defaultDocURL)
			}
		})
		var (
			p      string
			output = c.Output
		)
		if p = parseParamPath(c, prefix); p != "" {
			c.Request.URL.Path = p
		} else {
			p = strings.TrimPrefix(c.Request.URL.Path, prefix)
			p = strings.TrimPrefix(p, "/")
		}

		switch p {
		case defaultIndex:
			c.Output.Header(contentType, contentTypeHtml)
			if err := index.Execute(c.ResponseWriter, cfg); err != nil {
				c.Abort(500, err.Error())
			}
			return
		case defaultDocURL:
			var docs, _ = swag.ReadDoc()
			output.Header(contentType, contentTypeJson)
			if err := output.Body([]byte(docs)); err != nil {
				c.Abort(500, err.Error())
			}
			return
		case "", "/":
			c.Redirect(302, path.Join(prefix, defaultIndex))
			return
		default:
			fs.ServeHTTP(c.ResponseWriter, c.Request)
			return
		}
	}
}

// 解析路由
func parseParamPath(c *context.Context, prefix string) string {
	var (
		params = c.Input.Params()
		size   = len(params)
	)
	if size <= 1 {
		return defaultIndex
	}
	var splat, ok = params[routerKey]
	if !ok || strings.HasSuffix(splat, "/") {
		return defaultIndex
	}
	var keys []string
	delete(params, ":splat")
	for k := range params {
		keys = append(keys, k)
	}
	var (
		count int
		_path = defaultIndex
	)
	if keys == nil {
		return _path
	}
	if count = len(keys); count <= 0 {
		return _path
	}
	sort.Strings(keys)
	if v, ok := params[keys[count-1]]; ok {
		return v
	}
	return _path
}

// 解析前缀
func parsePrefix(url string) string {
	if url == "" {
		return ""
	}
	if strings.Contains(url, "/") {
		var arr = strings.Split(url, "/")
		return arr[0]
	}
	return url
}
