package util

import (
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"

	"github.com/EDDYCJY/go-gin-example/pkg/setting"
)

// GetPage get page parameters
func GetPage(c *gin.Context) int {
	result := 0
	// 获取 page 参数, StrTo .MustInt 强制转换为 int
	page := com.StrTo(c.Query("page")).MustInt()

	// 如果 page 大于 0
	// 则返回 page * setting.AppSetting.PageSize
	if page > 0 {
		result = (page - 1) * setting.AppSetting.PageSize
	}

	return result
}
