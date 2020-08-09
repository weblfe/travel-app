package plugins

import (
		"errors"
		"fmt"
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/logs"
		"github.com/boombuler/barcode"
		"github.com/boombuler/barcode/qr"
		"github.com/skip2/go-qrcode"
		qrcode2 "github.com/tuotoo/qrcode"
		"image/color"
		"image/png"
		"io"
		"net/http"
		"os"
		"path/filepath"
		"regexp"
		"time"
)

type Qrcode struct {
		Storage string  `json:"storage"`
		Config  beego.M `json:"config"`
}

type Options struct {
		Size            int                  `json:"size"`            // 二维码 （宽 X 高）
		FileName        string               `json:"filename"`        // 文件保存名路径
		BackgroundColor color.Color          `json:"backgroundColor"` // 背景色
		ForegroundColor color.Color          `json:"foregroundColor"` // 颜色
		Level           qrcode.RecoveryLevel `json:"level"`           // 质量
		Auto            bool                 `json:"auto"`            // 是否自动扩展 减少白边
}

const (
		QrcodePluginName = "qrcode"
)

var (
		_QrcodeIns *Qrcode
)

func GetQrcode() *Qrcode {
		if _QrcodeIns == nil {
				var locker = getLock(QrcodePluginName)
				locker.Do(newQrcode)
		}
		return _QrcodeIns
}

func newQrcode() {
		_QrcodeIns = new(Qrcode)
		_QrcodeIns.init()
}

func (this *Qrcode) init() {
		this.Config = beego.M{}
		this.Storage = ""
}

func (this *Qrcode) Register() {
		Plugin(this.PluginName(), this)
}

func (this *Qrcode) PluginName() string {
		return QrcodePluginName
}

func (this *Qrcode) Create(content string, options ...int) ([]byte, error) {
		return qrcode.Encode(content, this.getLevel(options), this.getSize(options))
}

func (this *Qrcode) Save(content string, options *Options) error {
		if options.FileName == "" {
				return errors.New("filename empty error")
		}
		options.Init()
		// 是否有颜色
		if options.Auto && !options.hasColor() {
				return this.createQRCodeByBoom(content, options.GetLevel(), options.Size, options.FileName)
		}
		if options.BackgroundColor == nil && options.ForegroundColor == nil {
				return qrcode.WriteFile(content, options.Level, options.Size, options.FileName)
		}
		return qrcode.WriteColorFile(content, options.Level, options.Size, options.BackgroundColor, options.ForegroundColor, options.FileName)
}

func (this *Qrcode) createQRCodeByBoom(content string, quality qr.ErrorCorrectionLevel, size int, dest string) (err error) {
		qrCode, err := qr.Encode(content, quality, qr.Auto)
		if err != nil {
				return
		}
		// Scale the barcode to 200x200 pixels
		qrCode, err = barcode.Scale(qrCode, size, size)
		if err != nil {
				return
		}
		// create the output file
		file, err := os.Create(dest)
		if err != nil {
				return
		}
		// 关闭
		defer func() {
				err := file.Close()
				if err != nil {
						logs.Error(err)
				}
		}()
		// encode the barcode as png
		err = png.Encode(file, qrCode)
		if err != nil {
				return
		}
		return
}

func (this *Qrcode) Decode(filename string) (*qrcode2.Matrix, error) {
		var reader, err = this.getReader(filename)
		if err != nil {
				return nil, err
		}
		m, err := qrcode2.Decode(reader)
		if err == nil {
				return m, nil
		}
		logs.Error(err)
		return nil, err
}

func (this *Qrcode) getReader(filename string) (io.ReadCloser, error) {
		var (
				urlRegexp = regexp.MustCompile(`^(http:|https:)`)
				urlMather = regexp.MustCompile(`[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(/.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+/.?`)
		)
		if filename == "" {
				return nil, errors.New("empty filename")
		}
		if urlRegexp.MatchString(filename) {
				var resp, err = http.Get(filename)
				if err == nil {
						return resp.Body, nil
				}
				return nil, err
		}
		var state, err = os.Stat(filename)
		if err == nil && !state.IsDir() {
				return os.Open(filename)
		}
		if urlMather.MatchString(filename) {
				var resp, err = http.Get("http://" + filename)
				if err == nil {
						return resp.Body, nil
				}
				return nil, err
		}
		return nil, err
}

func (this *Qrcode) getLevel(options []int) qrcode.RecoveryLevel {
		if len(options) < 0 {
				return qrcode.Medium
		}
		for _, level := range this.GetLevels() {
				if level == options[0] {
						return qrcode.RecoveryLevel(options[0])
				}
		}
		return qrcode.Medium
}

func (this *Qrcode) GetLevels() []int {
		return []int{
				int(qrcode.Low),
				int(qrcode.Medium),
				int(qrcode.High),
				int(qrcode.Highest),
		}
}

func (this *Qrcode) getSize(options []int) int {
		if len(options) < 2 {
				return 256
		}
		for _, level := range this.GetSizes() {
				if level == options[1] {
						return options[1]
				}
		}
		return 256
}

func (this *Qrcode) GetSizes() []int {
		return []int{60, 180, 256, 300, 600}
}

func (this *Options) GetStorage() string {
		if this.FileName != "" {
				return this.FileName
		}
		var saver = os.Getenv("STORAGE_PATH")
		if saver == "" {
				saver, _ = os.Getwd()
				saver = saver + "/static/storage"
		}
		var err = os.MkdirAll(saver, os.ModePerm)
		if err != nil {
				return filepath.Join(os.TempDir(), "/"+fmt.Sprintf("%v.png", time.Now().Unix()))
		}
		return filepath.Join(saver, "/"+fmt.Sprintf("%v.png", time.Now().Unix()))
}

func (this *Options) Init() {
		if this.FileName == "" {
				this.FileName = this.GetStorage()
		}
		if this.Size == 0 {
				this.Size = 256
		}
		if this.Level == 0 {
				this.Level = qrcode.Medium
		}
		if this.ForegroundColor == nil {
				if this.BackgroundColor == nil {
						return
				}
				this.ForegroundColor = color.RGBA{R: 0, G: 0, B: 0, A: 0}
		}
		if this.BackgroundColor == nil {
				if this.ForegroundColor == nil {
						return
				}
				this.BackgroundColor = color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
		}
}

func (this *Options) hasColor() bool {
		if this.BackgroundColor == nil && this.ForegroundColor == nil {
				return false
		}
		return true
}

func (this *Options) GetLevel() qr.ErrorCorrectionLevel {
		switch this.Level {
		case qrcode.Low:
				return qr.L
		case qrcode.Medium:
				return qr.M
		case qrcode.High:
				return qr.Q
		case qrcode.Highest:
				return qr.H
		}
		return qr.M
}
