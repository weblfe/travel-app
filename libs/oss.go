package libs

import (
		"fmt"
		"os"
)

type OssUrlBuilder interface {
		GetUrl(ty ...OssUrlType) string
		SetSourceUrl(string) OssUrlBuilder
}

type OssUrlType string

const (
		Row           OssUrlType = "row"
		Media         OssUrlType = "medium"
		Big           OssUrlType = "big"
		Small         OssUrlType = "small"
		DefOptFlag               = "@" // 类型参数符号
		OptFlagEnvKey            = "QINNIU_URL_OPT_FLAG"
)

type ossUrlBuilderImpl struct {
		url  string
		flag string
}

func OssUrl(url string) OssUrlBuilder {
		return &ossUrlBuilderImpl{
				url: url,
		}
}

func (this *ossUrlBuilderImpl) GetUrl(tys ...OssUrlType) string {
		if len(tys) == 0 {
				tys = append(tys, Row)
		}
		var ty = tys[0]
		switch ty {
		case Row:
				return this.format(this.url, Row)
		case Media:
				return this.format(this.url, Media)
		case Big:
				return this.format(this.url, Big)
		case Small:
				return this.format(this.url, Small)
		default:
				return this.url
		}
}

func (this *ossUrlBuilderImpl) format(url string, ty OssUrlType) string {
		return fmt.Sprintf("%s%s%s", url, this.getFlag(), ty)
}

func (this *ossUrlBuilderImpl) getFlag() string {
		if this.flag == "" {
				this.flag = os.Getenv(OptFlagEnvKey)
				if this.flag == "" {
						this.flag = DefOptFlag
				}
		}
		return this.flag
}

func (this *ossUrlBuilderImpl) SetSourceUrl(url string) OssUrlBuilder {
		this.url = url
		return this
}

func GetCdnUrl(url string,ty OssUrlType) string  {
		return OssUrl(url).GetUrl(ty)
}