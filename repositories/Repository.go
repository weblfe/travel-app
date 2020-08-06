package repositories

import (
		"fmt"
		"github.com/astaxie/beego"
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/middlewares"
		"github.com/weblfe/travel-app/models"
		"github.com/weblfe/travel-app/services"
		"github.com/weblfe/travel-app/transforms"
		"math"
		"strconv"
		"strings"
)

const (
		BigNumberK      = 1000
		BigNumberKUnit  = "k"
		BigNumberW      = BigNumberK * 10
		BigNumberWUnit  = "w"
		BigNumberKW     = BigNumberW * BigNumberK
		BigNumberKWUnit = "kw"
		BigNumberWW     = BigNumberKW * BigNumberK
		BigNumberWWUnit = "ww"
)

var (
		BigNumberMapper = []struct {
				Limiter int64
				Unit    string
		}{
				{BigNumberK, BigNumberKUnit},
				{BigNumberW, BigNumberWUnit},
				{BigNumberKW, BigNumberKWUnit},
				{BigNumberWW, BigNumberWWUnit},
		}
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
				m["isFollowed"] = false
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
				if strId == userId {
						m["isSelf"] = true
						return m
				}
				m["isFollowed"] = services.UserBehaviorServiceOf().IsFollowed(userId, strId)
				return m
		}
}

// 大数转换器
func TransBigNumberToText(m beego.M, keys ...string) beego.M {
		if len(keys) == 0 {
				return m
		}
		for _, key := range keys {
				numKey := key + "Text"
				m[numKey] = "0"
				if n, ok := m[key]; ok {
						m[numKey] = DecorateNumberToText(n)
				}
		}
		return m
}

// 大数字装饰器
func DecorateNumberToText(v interface{}) string {
		switch v.(type) {
		case string:
				return v.(string)
		case int:
				var num = v.(int)
				return BigNumberStringer(int64(num))
		case int64:
				var num = v.(int64)
				return BigNumberStringer(num)
		case int32:
				var num = v.(int32)
				return BigNumberStringer(int64(num))
		case int8:
				var num = v.(int8)
				return BigNumberStringer(int64(num))
		case int16:
				var num = v.(int16)
				return BigNumberStringer(int64(num))
		case float64:
				var num = v.(float64)
				return BigNumberStringer(int64(num))
		case float32:
				var num = v.(float32)
				return BigNumberStringer(int64(num))
		}
		return "0"
}

// 大数字转换
func BigNumberStringer(num int64) string {
		var (
				times      int64 = 0
				before     int64 = BigNumberK
				beforeUnit       = BigNumberKUnit
				total            = len(BigNumberMapper)
				n                = int64(math.Abs(float64(num)))
		)
		for i, mapper := range BigNumberMapper {
				times = 0
				if total > i+1 {
						times = before / BigNumberMapper[i+1].Limiter
				}
				if times > 1 && n/times > 1 {
						continue
				}
				if mapper.Limiter >= n && n/mapper.Limiter >= 1 {
						before = mapper.Limiter
						beforeUnit = mapper.Unit
						break
				}

		}
		if before == BigNumberK {
				if BigNumberK/2 > num {
						return fmt.Sprintf("%d", num)
				}
				if BigNumberK == num {
						return fmt.Sprintf("%d%s", int(float64(num)/float64(before)), beforeUnit)
				}
				return fmt.Sprintf("%s%s", DecimalText(float64(num)/float64(before)), beforeUnit)
		}
		return fmt.Sprintf("%d%s", num/before, beforeUnit)
}

func Round2(f float64, n int) float64 {
		floatStr := fmt.Sprintf("%."+strconv.Itoa(n)+"f", f)
		inst, _ := strconv.ParseFloat(floatStr, 64)
		return inst
}

// 保留
func DecimalText(f float64, n ...int) string {
		if len(n) == 0 {
				n = append(n, 2)
		}
		if n[0] < 0 {
				n[0] = 2
		}
		var (
				str  = fmt.Sprintf("%v", Decimal(f*math.Pow10(n[0]), n...))
				nums = strings.SplitN(str, ".", -1)
		)
		if len(nums) >= 1 {
				inst, err := strconv.ParseFloat(nums[0], 64)
				if err != nil {
						return "0." + strings.Repeat("0", n[0])
				}
				inst = inst / math.Pow10(n[0])
				return fmt.Sprintf("%v", inst)
		}
		return "0.000"
}

func Decimal(f float64, n ...int) float64 {
		if len(n) == 0 {
				n = append(n, 2)
		}
		if n[0] < 0 {
				n[0] = 2
		}
		n10 := math.Pow10(n[0])
		return math.Trunc((f+0.5/n10)*n10) / n10
}

// 对应用户的数据转换器
func getUserTransform(userId string) func(m beego.M) beego.M {
		return transforms.FilterWrapper(transforms.FilterUserBase,
				FilterUserAvatarUrl, FilterBackgroundUrl, transformUser,
				appendFollowedLogic(userId), transformUserNumber,
		)
}

// 用户数字装饰器
func transformUserNumber(m beego.M) beego.M {
		return TransBigNumberToText(m, "thumbsUpNum", "followNum", "fansNum")
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
