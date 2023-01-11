package logging

import (
	"fmt"
	"time"

	"github.com/EDDYCJY/go-gin-example/pkg/setting"
)

// getLogFilePath get the log file save path
// 获取日志文件保存路径
func getLogFilePath() string {
	return fmt.Sprintf("%s%s", setting.AppSetting.RuntimeRootPath, setting.AppSetting.LogSavePath)
}

// getLogFileName get the save name of the log file
// 获取日志文件的保存名称
func getLogFileName() string {
	return fmt.Sprintf("%s%s.%s",
		setting.AppSetting.LogSaveName,
		time.Now().Format(setting.AppSetting.TimeFormat),
		setting.AppSetting.LogFileExt,
	)
}
