package test

import (
		. "github.com/smartystreets/goconvey/convey"
		"github.com/weblfe/travel-app/models"
		"testing"
)

func TestUserAddress_String(t *testing.T) {
		var (
			address = models.NewUserAddress()
			text = "广东省广州市天河区五山路108号"
		)
		address.Parse(text)
		Convey("test convey",t, func() {
				So(address.String(),ShouldNotBeEmpty)
		})

}