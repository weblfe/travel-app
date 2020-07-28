package repositories

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/weblfe/travel-app/models"
	"github.com/weblfe/travel-app/services"
	"github.com/weblfe/travel-app/transforms"
)

// 头像图
func FilterUserAvatarUrl(m beego.M) beego.M {
	if url, ok := m["avatarUrl"]; ok && url != nil && url != "" {
		return m
	}
	if id, ok := m["avatarId"]; ok && id != nil && id != "" {
		m["avatarUrl"] = services.AvatarServerOf().GetAvatarUrlById(id.(string))
	} else {
		gender := m["gender"]
		if gender == nil {
			gender = 0
		}
		m["avatarUrl"] = services.AvatarServerOf().GetAvatarUrlDefault(gender.(int))
	}
	return m
}

// 背景墙图
func FilterBackgroundUrl(m beego.M) beego.M {
	if url, ok := m["backgroundCoverUrl"]; ok && url != nil && url != "" {
		return m
	}
	m["backgroundCoverUrl"] = ""
	if id, ok := m["backgroundCoverId"]; ok && id != nil && id != "" {
		m["backgroundCoverUrl"] = services.AttachmentServiceOf().GetAccessUrl(id.(string))
	}
	return m
}

// 获取用户基础数据转换器
func getBaseUserInfoTransform() func(m beego.M) beego.M {
	return transforms.FilterWrapper(transforms.FilterUserBase, FilterUserAvatarUrl, FilterBackgroundUrl)
}

// 获取media Transform
func getMediaInfoTransform() func(m beego.M) beego.M {
	return func(m beego.M) beego.M {
		var (
			key     = PostImagesInfoKey
			arr, ok = m[key]
			service = services.AttachmentServiceOf()
		)
		if !ok {
			key = PostVideoInfoKey
			arr, ok = m[key]
			if !ok {
				return m
			}
		}
		if key == PostImagesInfoKey {
			var items = arr.([]*models.Image)
			if items != nil && len(items) > 0 {
				for i, it := range items {
					it.Url = service.GetAccessUrl(it.MediaId)
					items[i] = it
				}
				m[key] = items
			}
			key = PostVideoInfoKey
			arr, ok = m[key]
		}
		if key == PostVideoInfoKey && ok {
			logs.Info("debug: ", arr)
			var items = arr.([]*models.Video)
			if items != nil && len(items) > 0 {
				for i, it := range items {
					it.Url = service.GetAccessUrl(it.MediaId)

					if it.CoverId != "" {
						it.CoverUrl = service.GetAccessUrl(it.CoverId)
					}
					items[i] = it
				}
				m[key] = items
			}
		}
		return m
	}
}
