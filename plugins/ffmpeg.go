package plugins

import (
		"context"
		"fmt"
		"os"
		"os/exec"
		"path/filepath"
		"runtime"
		"strings"
		"sync"
		"time"
		"unicode"
)

type ffmpeg struct {
		Binary       string        `json:"binary"`
		Version      string        `json:"version"`
		Storage      string        `json:"storage"`
		ScreenShotAt time.Duration `json:"screenShotSt"`
}

// 截图
type FFmpegPluginInterface interface {
		GetScreenShotAt() string
		SetScreenShotAt(at time.Duration)
		SetBinary(bin string) FFmpegPluginInterface
		SaveScreenShot(file string, storage ...string) string
}

const (
		_Which               = "which"
		_FFmpegPluginName    = "ffmpeg"
		_FFmpegBinaryName    = "ffmpeg"
		_FFmpegBinaryPath    = "/bin/ffmpeg"
		_DefaultStoragePath  = "/tmp/ffmpeg"
		_FFmpegEnvKey        = "FFMPEG_HOME"
		_FFmpegStorageEnvKey = "FFMPEG_STORAGE_DIR"
		_DefaultBinary       = "/usr/local/bin/ffmpeg"
		_DefaultShotTimeAt   = 0.01
		_ScreenShotTmpl      = " -t %s -i %s -f image2 -y %s "
)

var (
		_syncOne   = sync.Once{}
		_ffmpegIns *ffmpeg
)

// 截图
func ScreenShot(filename string, storage string, services ...FFmpegPluginInterface, ) bool {
		if services == nil || len(services) == 0 {
				services = append(services, GetFfmpeg())
		}
		var (
				service = services[0]
				saver   = service.SaveScreenShot(filename, storage)
		)
		if saver == "" {
				return false
		}
		return true
}

func isExists(filename string) bool {
		if filename == "" {
				return false
		}
		var state, err = os.Stat(filename)
		if err != nil {
				return false
		}
		if state.IsDir() {
				return false
		}
		return true
}

func GetFfmpeg() FFmpegPluginInterface {
		if _ffmpegIns == nil {
				_syncOne.Do(newFfmpeg)
		}
		return _ffmpegIns
}

func newFfmpeg() {
		_ffmpegIns = new(ffmpeg)
		_ffmpegIns.Init()
}

func (this *ffmpeg) PluginName() string {
		return _FFmpegPluginName
}

func (this *ffmpeg) SaveScreenShot(file string, storage ...string) string {
		if this.Binary == "" || !isExists(file) {
				return ""
		}
		if len(storage) == 0 {
				storage = append(storage, this.getStorage())
		}
		var (
				image    = this.getStorageImage(storage[0])
				cmder, _ = this.commander(this.getExec(), this.getOption(), this.Binary+fmt.Sprintf(_ScreenShotTmpl, this.GetScreenShotAt(), file, image))
		)
		var info, err = cmder.CombinedOutput()
		if err != nil {
				state, err := os.Stat(image)
				if err == nil && state.Size() > 0 {
						return image
				}
				fmt.Println(string(info), err)
				return ""
		}
		if string(info) != "" {
				return image
		}
		return ""
}

func (this *ffmpeg) getStorageImage(file string) string {
		if file == "" {
				return filepath.Join(this.getStorage(), time.Now().String()+".jpg")
		}
		if !isExists(file) && !strings.Contains(file, ".jpg") {
				return filepath.Join(file, time.Now().String()+".jpg")
		}
		return file
}

func (this *ffmpeg) getExec() string {
		switch runtime.GOOS {
		case "linux":
				return "bash"
		case "windows":
				return "cmd"
		case "darwin":
				return "bash"
		}
		return ""
}

func (this *ffmpeg) getOption() string {
		switch runtime.GOOS {
		case "linux":
				return "-c"
		case "windows":
				return "/c"
		case "darwin":
				return "-c"
		}
		return ""
}

// 设置ffmpeg bin 安装位置
func (this *ffmpeg) SetBinary(bin string) FFmpegPluginInterface {
		if bin == "" || !isExists(bin) {
				return this
		}
		this.Binary = bin
		return this
}

// 设置截图时间
func (this *ffmpeg) SetScreenShotAt(at time.Duration) {
		if at <= 0 {
				return
		}
		this.ScreenShotAt = at
}

func (this *ffmpeg) isNumber(str string) bool {
		if str == "" {
				return false
		}
		var chars = []rune(str)
		for _, ch := range chars {
				if !unicode.IsNumber(ch) {
						return false
				}
		}
		return true
}

func (this *ffmpeg) GetScreenShotAt() string {
		if this.ScreenShotAt == 0 {
				return fmt.Sprintf("%v", _DefaultShotTimeAt)
		}
		return fmt.Sprintf("%v", this.ScreenShotAt/time.Second)
}

func (this *ffmpeg) Register() {
		Plugin(this.PluginName(), this)
}

func (this *ffmpeg) Init() *ffmpeg {
		if this.Version == "" {
				this.Version = "4.2.2"
		}
		if this.Binary == "" {
				this.Binary = this.getBinaryAuto()
		}
		if this.Storage == "" {
				this.Storage = this.getStorage()
		}
		this.Register()
		return this
}

func (this *ffmpeg) getStorage() string {
		if this.Storage != "" {
				return this.Storage
		}
		var dir = os.Getenv(_FFmpegStorageEnvKey)
		if isDir(dir) {
				return dir
		}
		// 创建
		if !isDir(_DefaultStoragePath) {
				_ = os.MkdirAll(_DefaultStoragePath, os.ModePerm)
		}
		return _DefaultStoragePath
}

func (this *ffmpeg) getBinaryAuto() string {
		var (
				cmder, _ = this.commander(_Which, _FFmpegBinaryName)
		)
		// which
		result, err := cmder.CombinedOutput()
		if err == nil && string(result) != "" {
				var bin = string(result)
				if strings.Contains(bin, _FFmpegBinaryName) {
						return strings.TrimSpace(bin)
				}
		}
		var bin = os.Getenv(_FFmpegEnvKey)
		if bin != "" {
				bin = filepath.Join(bin, _FFmpegBinaryPath)
				if isExists(bin) {
						return bin
				}
		}
		var _path = os.Getenv("PATH")
		// win
		if strings.Contains(_path, ";") {
				bin = this.find(_path, ";")
				if bin != "" {
						return bin
				}
		}
		// linux
		if strings.Contains(_path, ":") {
				bin = this.find(_path, ":")
				if bin != "" {
						return bin
				}
		}
		return _DefaultBinary
}

// 查找
func (this *ffmpeg) find(_path string, split string) string {
		var (
				bin   string
				paths = strings.SplitN(_path, split, -1)
		)
		for _, p := range paths {
				if !strings.Contains(p, _FFmpegBinaryName) {
						continue
				}
				bin = filepath.Join(p, _FFmpegBinaryPath)
				if isExists(bin) {
						return bin
				}
		}
		return ""
}

func (this *ffmpeg) Clear() {
		if this.Storage == _DefaultStoragePath {
				_ = os.RemoveAll(this.Storage)
		}
}

func (this *ffmpeg) commander(bin string, args ...string) (*exec.Cmd, *CContext) {
		// 穿件上下文,和取消函数
		var (
				ctxParent, timeOutCancelFn = context.WithTimeout(context.TODO(), time.Minute*3)
				ctx, cancelFn              = context.WithCancel(ctxParent)
				cmd                        = exec.CommandContext(ctx, bin, args...)
				ctrl                       = new(CContext)
		)
		ctrl.Cancel = cancelFn
		ctrl.TimeoutFn = timeOutCancelFn
		return cmd, ctrl
}

type CContext struct {
		TimeoutFn context.CancelFunc
		Cancel    context.CancelFunc
}

func isDir(dir string) bool {
		if state, err := os.Stat(dir); err == nil {
				if state.IsDir() {
						return true
				}
		}
		return false
}
