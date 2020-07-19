package repositories

import (
		"github.com/astaxie/beego"
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
