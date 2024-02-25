package libs

import jsonApi "github.com/json-iterator/go"
import "encoding/json"

// import jsonApi "github.com/json-iterator/go"

var (
		_JsonErr error
)

func Json() jsonApi.API {
		return jsonApi.ConfigCompatibleWithStandardLibrary
}

// JsonEncode 编码
func JsonEncode(v interface{}) []byte {
		var (
				data, err = json.Marshal(v)
		)
		if err != nil {
				_JsonErr = err
		}
		return data
}

// JsonDecode json 解码
func JsonDecode(data []byte, bindTo ...interface{}) interface{} {
		if len(bindTo) == 0 {
				var d = map[string]interface{}{}
				bindTo = append(bindTo, &d)
		}
		var (
				err = json.Unmarshal(data, bindTo[0])
		)
		if err != nil {
				_JsonErr = err
		}
		return bindTo[0]
}

// JsonDecodeBy 解码
func JsonDecodeBy(str string, bindTo ...interface{}) interface{} {
		return JsonDecode([]byte(str), bindTo...)
}

// GetLastJsonErr 获取 异常
func GetLastJsonErr() error  {
		defer func() {
				_JsonErr = nil
		}()
		return _JsonErr
}