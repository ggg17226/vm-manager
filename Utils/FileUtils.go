package Utils

import (
	"math/rand"
	"os"
	"strconv"
	"time"
)

/**
检查文件是否存在
*/
func FileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

/**
生成文件名
*/
func GenFileName() string {
	timeStr := time.Now().Format("2006-01-02T15-04-05Z0700")
	return timeStr + "-" + strconv.Itoa(rand.Intn(89999)+10000)
}
