package utils

import "encoding/json"

// PrintJsonString 转化成json字符串
func PrintJsonString(data interface{}) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(jsonData)
}
