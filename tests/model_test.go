package test

import (
		"fmt"
		"github.com/astaxie/beego"
		"github.com/globalsign/mgo/bson"
		. "github.com/smartystreets/goconvey/convey"
		"github.com/weblfe/travel-app/libs"
		"github.com/weblfe/travel-app/models"
		"testing"
		"time"
)

func init() {
		beego.TestBeegoInit("")
		initDatabase()
}

func TestCollection(t *testing.T) {
		var user = models.UserModelOf()
		id :=libs.NewId(user.GetDatabaseName(), user.TableName(), user.GetConn())
		err := user.Collection().Insert(models.User{
				Id:           bson.NewObjectId(),
				UserNumId:    id.GetId(),
				UserName:     "test"+ fmt.Sprintf("%d",id.GetId()),
				PasswordHash: libs.PasswordHash("12323"),
				Mobile:       "13112260971",
				RegisterWay:  "mobile",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
				DeletedAt:    0,
				Status:       1,
		})
		if err == nil {
				_ = id.Commit()
		}else{
			 _ =	id.RollBack()
		}
		Convey("mongodb Test", t, func() {
				So(err, ShouldBeNil)
				So(libs.GetId(user.GetDatabaseName(), user.TableName(), user.GetConn()), ShouldNotBeEmpty)
		})
}

func BenchmarkGetId(b *testing.B) {
		var user = models.UserModelOf()

		for i := 0; i < b.N; i++ {
				//	fmt.Println(i)
				_ = user.Insert(models.User{
						Id:           bson.NewObjectId(),
						UserNumId:    libs.GetId(user.GetDatabaseName(), user.TableName(), user.GetConn()),
						UserName:     "test",
						PasswordHash: libs.PasswordHash("12323"),
						Mobile:       "13112260977",
						RegisterWay:  "mobile",
						CreatedAt:    time.Now(),
						UpdatedAt:    time.Now(),
						DeletedAt:    0,
						Status:       1,
				})
		}
}

func BenchmarkAdd(b *testing.B) {
		var user = models.UserModelOf()
		for i := 0; i < b.N; i++ {
				//	fmt.Println(i)
				libs.GetId(user.GetDatabaseName(), user.TableName(), user.GetConn())
		}
}

func initDatabase() {
		mode := beego.BConfig.RunMode
		if database, err := beego.AppConfig.GetSection(mode + ".database"); err == nil {
				// fmt.Println(database)
				if driver, ok := database["db_driver"]; ok && driver == "mongodb" {
						initMongodb(database)
				}
		}

}

func initMongodb(data map[string]string) {
		for key, v := range data {
				models.SetProfile(key, v)
		}
}
