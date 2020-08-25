package services

import (
		"fmt"
		"github.com/astaxie/beego"
		"github.com/jordan-wright/email"
		"github.com/weblfe/travel-app/libs"
		"net/smtp"
		"os"
		"strings"
)

type EmailService interface {
		Sends(target EmailRequest) error
		Send(emailTo string, emailFrom string, content string) interface{}
		SendHtml(emailTo string, emailFrom string, content string) interface{}
		SendWithCallback(target EmailRequest, callback func(target EmailRequest, result interface{}))
}

type EmailRequest interface {
		To() []string
		From() string
		Content() string
		Extras() map[string]interface{}
}

type EmailRequestImpl struct {
		to      []string
		from    string
		content string
		extras  map[string]interface{}
}

type EmailServiceImpl struct {
		BaseService
		config map[string]string
}

const (
		EmailSmtpHost      = "smtp_host"
		DefaultHost        = "smtp.126.com"
		DefaultUserName    = "test"
		DefaultPassword    = ""
		DefaultSmtpPort    = "25"
		EmailSmtpUserName  = "smtp_username"
		EmailSmtpPassword  = "smtp_password"
		EmailSmtpPort      = "smtp_port"
		EmailSmtpFromEmail = "smtp_from"
		EmailIsHtml        = "isHtml"
		EmailAttachments   = "attachments"
		EmailFiles         = "files"
		EmailSubject       = "subject"
		EmailHeaders       = "headers"
		EmailContent       = "content"
		EmailTo            = "to"
		EmailFrom          = "from"
		EmailSender        = "sender"
)

var (
		DefaultTslSmtpPorts = []string{"465", "587"}
)

func NewEmailRequest() *EmailRequestImpl {
		var request = new(EmailRequestImpl)
		request.to = []string{}
		request.extras = make(map[string]interface{})
		return request
}

func EmailServiceOf() EmailService {
		var service = new(EmailServiceImpl)
		service.Init()
		return service
}

func (this *EmailRequestImpl) To() []string {
		return this.to
}

func (this *EmailRequestImpl) From() string {
		return this.from
}

func (this *EmailRequestImpl) Content() string {
		return this.content
}

func (this *EmailRequestImpl) Extras() map[string]interface{} {
		return this.extras
}

func (this *EmailRequestImpl) Set(key string, v string) *EmailRequestImpl {
		if key == EmailFrom {
				this.from = v
				return this
		}
		if key == EmailTo {
				this.to = append(this.to, v)
				return this
		}
		if key == EmailContent {
				this.content = v
				return this
		}
		if key == EmailHeaders {
				if _, ok := this.extras[key]; !ok {
						this.extras[key] = map[string]string{}
				}
				arr, ok := this.extras[key]
				if ok {
						headers, ok := arr.(map[string]string)
						if ok {
								headers[key] = v
						}
				}
				return this
		}
		this.extras[key] = v
		return this
}

func (this *EmailRequestImpl) AddFile(fs string) *EmailRequestImpl {
		if _, ok := this.extras[EmailFiles]; !ok {
				this.extras[EmailFiles] = []string{}
		}
		files := this.extras[EmailFiles]
		if !libs.IsExits(fs) {
				return this
		}
		if arr, ok := files.([]string); ok {
				arr = append(arr, fs)
				this.extras[EmailFiles] = arr
		}
		return this
}

func (this *EmailServiceImpl) Init() {
		if this.config == nil {
				this.config = make(map[string]string)
		}
		this.Constructor = func(args ...interface{}) interface{} {
				return EmailServiceOf()
		}
		if len(this.config) == 0 {
				this.load()
		}
}

func (this *EmailServiceImpl) load() {
		this.config[EmailSmtpHost] = beego.AppConfig.DefaultString(EmailSmtpHost, this.getEmailHost())
		this.config[EmailSmtpPort] = beego.AppConfig.DefaultString(EmailSmtpPort, this.getSmtpPort())
		this.config[EmailSmtpFromEmail] = beego.AppConfig.DefaultString(EmailSmtpFromEmail, this.getFrom())
		this.config[EmailSmtpUserName] = beego.AppConfig.DefaultString(EmailSmtpUserName, this.getUserName())
		this.config[EmailSmtpPassword] = beego.AppConfig.DefaultString(EmailSmtpPassword, this.getPassword())
}

func (this *EmailServiceImpl) Sends(target EmailRequest) error {
		return this.send(this.make(target))
}

func (this *EmailServiceImpl) Send(emailTo string, emailFrom string, content string) interface{} {
		var sender = email.NewEmail()
		sender.From = emailFrom
		sender.To = []string{emailTo}
		sender.Text = []byte(content)
		return this.send(sender)
}

func (this *EmailServiceImpl) SendHtml(emailTo string, emailFrom string, content string) interface{} {
		var sender = email.NewEmail()
		sender.From = emailFrom
		sender.To = []string{emailTo}
		sender.HTML = []byte(content)
		return this.send(sender)
}

func (this *EmailServiceImpl) SendWithCallback(target EmailRequest, callback func(target EmailRequest, result interface{})) {
		var res = this.send(this.make(target))
		callback(target, res)
}

func (this *EmailServiceImpl) make(target EmailRequest) *email.Email {
		var (
				from          = target.From()
				extras        = target.Extras()
				sender        = email.NewEmail()
				attachmentArr []*email.Attachment
		)
		sender.To = target.To()
		if from == "" {
				sender.From = this.getFrom()
		} else {
				sender.From = from
		}
		h, ok := extras[EmailIsHtml]
		if ok {
				str, ok := h.(string)
				if ok && (str == "true" || str == "1") {
						sender.HTML = []byte(target.Content())
				}
		}
		if len(sender.HTML) == 0 {
				sender.Text = []byte(target.Content())
		}
		attachments, ok := extras[EmailAttachments]
		if ok {
				if att, ok := attachments.([]*email.Attachment); ok {
						sender.Attachments = att
				}
		}
		if len(sender.Attachments) == 0 {
				files, ok := extras[EmailFiles]
				if !ok {
						return sender
				}
				fsArr, ok := files.([]string)
				if !ok {
						return sender
				}
				for _, fs := range fsArr {
						a, err := sender.AttachFile(fs)
						if err == nil {
								attachmentArr = append(attachmentArr, a)
						}
				}
		}
		subject, ok := extras[EmailSubject]
		if ok {
				sender.Subject = subject.(string)
		}
		headers, ok := extras[EmailHeaders]
		if ok {
				if mapper, ok := headers.(map[string]string); ok {
						for key, v := range mapper {
								sender.Headers.Set(key, v)
						}
				}
		}
		if len(attachmentArr) != 0 && len(sender.Attachments) == 0 {
				sender.Attachments = append(sender.Attachments, attachmentArr...)
		}
		// 发送者的名字
		senderName, ok := extras[EmailSender]
		if ok && senderName != nil {
				sender.Sender = senderName.(string)
		}
		return sender
}

func (this *EmailServiceImpl) getFrom() string {
		if from, ok := this.config[EmailSmtpFromEmail]; ok && from != "" {
				return from
		}
		if from := os.Getenv(strings.ToUpper(EmailSmtpFromEmail)); from != "" {
				return from
		}
		return ""
}

func (this *EmailServiceImpl) send(sender *email.Email) error {
		return sender.Send(this.getSmtpHost(), this.getAuth())
}

func (this *EmailServiceImpl) IsTsl() bool {
		var p = this.getSmtpPort()
		for _, port := range DefaultTslSmtpPorts {
				if p == port {
						return true
				}
		}
		return false
}

func (this *EmailServiceImpl) getAuth() smtp.Auth {
		return smtp.PlainAuth("", this.getUserName(), this.getPassword(), this.getEmailHost())
}

func (this *EmailServiceImpl) getEmailHost() string {
		if host, ok := this.config[EmailSmtpHost]; ok && host != "" {
				return host
		}
		if host := os.Getenv(strings.ToUpper(EmailSmtpHost)); host != "" {
				return host
		}
		return DefaultHost
}

func (this *EmailServiceImpl) getUserName() string {
		if username, ok := this.config[EmailSmtpUserName]; ok && username != "" {
				return username
		}
		if username := os.Getenv(strings.ToUpper(EmailSmtpUserName)); username != "" {
				return username
		}
		return DefaultUserName
}

func (this *EmailServiceImpl) getPassword() string {
		if pass, ok := this.config[EmailSmtpPassword]; ok && pass != "" {
				return pass
		}
		if pass := os.Getenv(strings.ToUpper(EmailSmtpPassword)); pass != "" {
				return pass
		}
		return DefaultPassword
}

func (this *EmailServiceImpl) getSmtpPort() string {
		if port, ok := this.config[EmailSmtpPort]; ok && port != "" {
				return port
		}
		if port := os.Getenv(strings.ToUpper(EmailSmtpPort)); port != "" {
				return port
		}
		return DefaultSmtpPort
}

func (this *EmailServiceImpl) getSmtpHost() string {
		return fmt.Sprintf("%s:%s", this.getEmailHost(), this.getSmtpPort())
}
