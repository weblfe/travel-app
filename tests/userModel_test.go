package test

import (
		"fmt"
		"github.com/globalsign/mgo/bson"
		. "github.com/smartystreets/goconvey/convey"
		"github.com/weblfe/travel-app/services"
		"testing"
)

func TestCreateUser(t *testing.T) {
		_ = services.UserServiceOf().Add(map[string]interface{}{
				"username":     "app1",
				"nickname":     "app_nickname",
				"avatar":       bson.NewObjectId().String(),
				"mobile":       "13112260987",
				"passwordHash": "123456",
		})
}

func TestUser(t *testing.T) {
		Convey("Test Register User", t, func() {
				user, token, errs := services.AuthServiceOf().LoginByUserPassword("username", "app1", "123456")
				So(errs, ShouldBeNil)
				So(user, ShouldNotBeNil)
				So(token, ShouldNotBeEmpty)
				user2, err := services.AuthServiceOf().GetByAccessToken(token)
				So(user2, ShouldNotEqual, user)
				So(err, ShouldBeNil)
		})
}

func TestArrayDelete(t *testing.T) {
		var (
				index = 0
				arr   = []string{"1", "2", "3", "4"}
		)

		Convey("Test Array Delete", t, func() {
				fmt.Println(cap(arr))
				fmt.Println(arr[index+1:])
				arr = append(arr[:index], arr[index+1:]...)
				fmt.Println(arr, len(arr))
				So(len(arr) == 3, ShouldBeTrue)
		})
}
