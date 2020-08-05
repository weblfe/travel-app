package repositories

import (
		"github.com/astaxie/beego"
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/middlewares"
		"github.com/weblfe/travel-app/models"
		"github.com/weblfe/travel-app/services"
		"github.com/weblfe/travel-app/transforms"
)

// 获取登录用户ID
func getUserId(request common.BaseRequestContext) string {
		var (
				session = request.GetSession()
				userId  = session.Get(middlewares.AuthUserId)
		)
		if userId == nil || userId == "" {
				return ""
		}
		if id, ok := userId.(string); ok {
				return id
		}
		return ""
}

// 获取登录用户
func getUser(request common.BaseRequestContext) *models.User {
		var id = getUserId(request)
		if id == "" {
				return nil
		}
		return services.UserServiceOf().GetById(id)
}

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
		return transforms.FilterWrapper(transforms.FilterUserBase, FilterUserAvatarUrl, FilterBackgroundUrl, transformUser)
}

// 追加用户关注关系
func appendFollowedLogic(userId string) func(m beego.M) beego.M {
		return func(m beego.M) beego.M {
				if userId == "" {
						return m
				}
				id, ok := m["id"]
				if !ok {
						return m
				}
				strId, ok := id.(string)
				if !ok {
						return m
				}
				m["isFollowed"] = false
				if strId == userId {
						m["isSelf"] = true
						return m
				}
				m["isFollowed"] = services.UserBehaviorServiceOf().IsFollowed(userId, strId)
				return m
		}
}

// 对应用户的数据转换器
func getUserTransform(userId string) func(m beego.M) beego.M {
		return transforms.FilterWrapper(transforms.FilterUserBase,
				FilterUserAvatarUrl, FilterBackgroundUrl, transformUser,
				appendFollowedLogic(userId),
		)
}

// 用户数据组装
func transformUser(m beego.M) beego.M {
		if _, ok := m["avatar"]; ok {
				return m
		}
		if avatarId, ok := m["avatarId"]; ok {
				if avatarUrl, ok1 := m["avatarUrl"]; ok1 {
						var avatar = new(Avatar)
						avatar.Id = avatarId.(string)
						avatar.AvatarUrl = avatarUrl.(string)
						m["avatar"] = avatar
				}
				delete(m, "avatarId")
				delete(m, "avatarUrl")
		}
		var userId, ok2 = m["id"]
		if !ok2 || userId == nil || userId == "" {
				return m
		}
		var (
				id      = userId.(string)
				service = services.UserServiceOf()
		)
		if _, ok := m["fansNum"]; !ok {
				m["fansNum"] = service.GetUserFansCount(id)
		}
		if _, ok := m["followNum"]; !ok {
				m["followNum"] = service.GetUserFollowCount(id)
		}
		return m
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

// 获取token
func getToken(request common.BaseRequestContext) string {
		return request.GetHeader().Get(middlewares.AppAccessTokenHeader)
}
