package tool

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"os"
	"reflect"
	"strconv"
	"time"
)

// Contain judge item is in array
func Contain(array interface{}, item interface{}) bool {
	switch value := array.(type) {
	case []string:
	case []int:
		for _, eachItem := range value {
			if eachItem == item {
				return true
			}
		}
	}
	return false
}

// GetMapFormString
func GetMapFormString(str string) (m map[string]string) {
	json.Unmarshal([]byte(str), &m)
	return
}

// Typeof
func Typeof(v interface{}) string {
	return reflect.TypeOf(v).String()
}

// var ProfileListMap map[string]string
//

// SHA256Str SHA256 encryption
func SHA256Str(src string) string {
	h := sha256.New()
	h.Write([]byte(src))
	return hex.EncodeToString(h.Sum(nil))
}

// PathFileExists check file exist
func PathFileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// GenerateTimestamp  generate timestamp
func GenerateTimestamp() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

// GetHeaderValue get header value by header key
func GetHeaderValue(c *gin.Context, headerKey string) (headerValue string) {
	for k, v := range c.Request.Header {
		if k == headerKey {
			headerValue = v[0]
			return
		}
	}
	return
}
