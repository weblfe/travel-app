package plugins

type FFmpeg struct {
		Binary  string `json:"binary"`
		Version string `json:"version"`
		Storage string `json:"storage"`
}

// ffmpeg -ss 00:02:06 -i test1.flv -f image2 -y test1.jpg

func ScreenShot(filename string,fmpeg FFmpeg,storage...string) bool{

		return false
}