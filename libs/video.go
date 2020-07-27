package libs

import (
		"bytes"
		"encoding/binary"
		"errors"
		"io"
		"os"
		"path/filepath"
		"strings"
		"time"
)

const (
		Mov = "moov"
		Mat = "mdat"
)

// BoxHeader 信息头
type BoxHeader struct {
		Size     uint32
		FourType [4]byte
		Size64   uint64
}

// 获取mp4 时长
func GetMp4FileDuration(fileName string) (time.Duration, error) {
		var (
				fd, err = os.Open(fileName)
		)
		// 文件类型判断
		if !strings.Contains(filepath.Ext(fileName), "mp4") {
				return 0, errors.New("file ext type error")
		}
		if err != nil {
				return 0, err
		}
		defer func(fd *os.File) {
				_ = fd.Close()
		}(fd)
		d, err := GetMP4Duration(fd)
		if err == nil {
				return time.Duration(int64(d)) * time.Second, nil
		}
		return 0, err
}

// GetMP4Duration 获取视频时长，以秒计
func GetMP4Duration(reader io.ReaderAt) (lengthOfTime uint32, err error) {
		var (
				boxHeader BoxHeader
				offset    int64 = 0
				info            = make([]byte, 0x10)
		)
		// 获取moov结构偏移
		for {
				_, err = reader.ReadAt(info, offset)
				if err != nil {
						return
				}
				boxHeader = getHeaderBoxInfo(info)
				fourType := getFourType(boxHeader)
				if fourType == Mov {
						break
				}
				// 有一部分mp4 mdat尺寸过大需要特殊处理
				if fourType == Mat {
						if boxHeader.Size == 1 {
								offset += int64(boxHeader.Size64)
								continue
						}
				}
				offset += int64(boxHeader.Size)
		}
		// 获取moov结构开头一部分
		moovStartBytes := make([]byte, 0x100)
		_, err = reader.ReadAt(moovStartBytes, offset)
		if err != nil {
				return
		}
		// 定义timeScale与Duration偏移
		timeScaleOffset := 0x1C
		durationOffset := 0x20
		timeScale := binary.BigEndian.Uint32(moovStartBytes[timeScaleOffset : timeScaleOffset+4])
		Duration := binary.BigEndian.Uint32(moovStartBytes[durationOffset : durationOffset+4])
		lengthOfTime = Duration / timeScale
		return
}

// getHeaderBoxInfo 获取头信息
func getHeaderBoxInfo(data []byte) (boxHeader BoxHeader) {
		buf := bytes.NewBuffer(data)
		_ = binary.Read(buf, binary.BigEndian, &boxHeader)
		return
}

// getFourType 获取信息头类型
func getFourType(boxHeader BoxHeader) (fourType string) {
		fourType = string(boxHeader.FourType[:])
		return
}
