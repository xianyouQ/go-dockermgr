package utils

import (
	"crypto/md5"
	"encoding/hex"
	"reflect"
	"strconv"
)

func Strtomd5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	rs := hex.EncodeToString(h.Sum(nil))
	return rs
}

//password hash function
func Pwdhash(str string) string {
	return Strtomd5(str)
}

func StringsToJson(str string) string {
	rs := []rune(str)
	jsons := ""
	for _, r := range rs {
		rint := int(r)
		if rint < 128 {
			jsons += string(r)
		} else {
			jsons += "\\u" + strconv.FormatInt(int64(rint), 16) // json
		}
	}

	return jsons
}

func ObjToMap(obj interface{}, preserves ...string) map[string]interface{} {
	//t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj).Elem()
	t := v.Type()
	var data = make(map[string]interface{})
	for i := 0; i < v.NumField(); i++ {
		if len(preserves) == 0 {
			data[t.Field(i).Name] = v.Field(i).Interface()
		} else {
			for _, preserve := range preserves {
				if preserve == t.Field(i).Name {
					data[t.Field(i).Name] = v.Field(i).Interface()
				}
			}
		}
	}
	return data
}
