package utils

import (
    "encoding/json"
)


func Json2Map(jsonstr []byte) (s map[string]interface{}, err error) {
    var result map[string]interface{}
    if err := json.Unmarshal(jsonstr, &result); err != nil {
        return nil,err
    }
    return result,nil
}