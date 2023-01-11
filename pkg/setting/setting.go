package setting

import (
	"log"
	"time"

	"github.com/go-ini/ini"
)

type App struct {

	// jwt 加密密钥
	JwtSecret string

	// 默认分页数量
	PageSize int

	// 服务器地址
	PrefixUrl string

	// 工作目录
	RuntimeRootPath string

	// 图像上传存储目录
	ImageSavePath string

	// 图像大小限制 (M)
	ImageMaxSize int

	// 图像后缀列表
	ImageAllowExts []string

	// 导出文件存储目录
	ExportSavePath string

	// 二维码存储目录
	QrCodeSavePath string

	// 字体文件存储目录
	FontSavePath string

	// 日志文件存储目录
	LogSavePath string

	// 日志文件名
	LogSaveName string

	// 日志文件后缀
	LogFileExt string

	// 时间格式
	TimeFormat string
}

var AppSetting = &App{}

type Server struct {

	// 运行模式, debug, release
	RunMode string

	// HTTP端口
	HttpPort int

	// 读取超时时间
	ReadTimeout time.Duration

	// 写入超时时间
	WriteTimeout time.Duration
}

var ServerSetting = &Server{}

type Database struct {

	// 数据库类型
	Type string

	// 数据用户
	User string

	// 数据库密码
	Password string

	// 数据库地址
	Host string

	// 数据库名称
	Name string

	// 数据库表前缀
	TablePrefix string
}

var DatabaseSetting = &Database{}

type Redis struct {

	// 地址
	Host string

	// redis 密码
	Password string

	// 连接池最大空闲连接数
	MaxIdle int

	// 连接池最大连接数
	MaxActive int

	// 连接池最大空闲时间
	IdleTimeout time.Duration
}

var RedisSetting = &Redis{}

var cfg *ini.File

// Setup initialize the configuration instance
func Setup() {
	var err error
	cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'conf/app.ini': %v", err)
	}

	mapTo("app", AppSetting)
	mapTo("server", ServerSetting)
	mapTo("database", DatabaseSetting)
	mapTo("redis", RedisSetting)

	// MB -> B
	AppSetting.ImageMaxSize = AppSetting.ImageMaxSize * 1024 * 1024

	// second -> time.Second
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second
	RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second
}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}
