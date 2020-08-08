package test

import (
		. "github.com/smartystreets/goconvey/convey"
		"github.com/weblfe/travel-app/plugins"
		"testing"
)

func TestFFmpeg(t *testing.T) {
		var files = []map[string]string{
				{"/Users/hewei/Desktop/12.mp4": "/Users/hewei/Desktop/12.jpg"},
		}
		// var ctx = context.TODO()
	//	var cmder = exec.CommandContext(ctx, "which ", "ffmpeg")
  //  exec.Command("which","ffmpeg").Run()
	//	var result, err = cmder.CombinedOutput()
		// fmt.Println(string(result), err)
		Convey("Test ffmpeg screenShot ", t, func() {
				for _, data := range files {
						for file, saver := range data {
								var f = plugins.GetFfmpeg().SaveScreenShot(file, saver)
								So(f == saver, ShouldBeTrue)
						}

				}
		})
}
