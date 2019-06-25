package models

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	//  _ "github.com/jinzhu/gorm/dialects/postgres"
	//  _ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/gocommon/models/errors"
)

// Drivers Drivers
var Drivers map[string]*gorm.DB

// GormService GormService
type GormService struct {
	Enable bool   `dsn:"query.enable"`
	Debug  bool   `dsn:"query.debug"`
	Driver string `dsn:"query.driver"`
	Host   string `dsn:"address"`
	DB     string `dsn:"query.db"`
	User   string `dsn:"username"`
	Passwd string `dsn:"password"`

	Path string `dsn:"query.path"` // for sqlite,tidb

	MaxIdle int    `dsn:"query.maxidle"` // 连接池的空闲数大小
	MaxOpen int    `dsn:"query.maxopen"` // 最大打开连接数
	LogPath string `dsn:"query.logpath"`
}

// Model Model
func Model(name ...string) *gorm.DB {
	k := "default"
	if len(name) > 0 {
		k = name[0]
	}
	if db, ok := Drivers[k]; ok {
		return db
	}

	panic(errors.New("model 不存在 %s", k))

}

// InitModels InitModels
func InitModels(confs map[string]GormService) error {
	Drivers = make(map[string]*gorm.DB)
	for k, v := range confs {
		if !v.Enable {
			continue
		}
		db, err := newGorm(v)
		if err != nil {
			return err
		}

		Drivers[k] = db
	}

	return nil
}

func newGorm(conf GormService) (*gorm.DB, error) {
	dsn := ""
	switch conf.Driver {
	case "mysql":
		//[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
		dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&zeroDateTimeBehavior=convertToNull", conf.User, conf.Passwd, conf.Host, conf.DB)
	// case "postgres":
	// case "sqlite3":
	default:
		return nil, errors.New("未知的 gorm 驱动：%s", conf.Driver)

	}

	db, err := gorm.Open(conf.Driver, dsn)
	if err != nil {
		return nil, errors.Wrap(err, "连接数据库失败")
	}

	db.SingularTable(true)

	if conf.MaxIdle > 0 {
		db.DB().SetMaxIdleConns(conf.MaxIdle)
	}

	if conf.MaxOpen > 0 {
		db.DB().SetMaxOpenConns(conf.MaxOpen)
	}

	// if conf.Debug {
	// 	db.LogMode(true)
	// }

	// logpath := "./log/gorm.log"
	// if len(conf.LogPath) > 0 {
	// 	logpath, _ = filepath.Abs(conf.LogPath)
	// }

	// os.MkdirAll(path.Dir(logpath), os.ModePerm)

	// fd, err := rotatefile.NewWriter(rotatefile.Options{Filename:logpath})

	// if err != nil {
	// 	return nil, err
	// }

	// db.SetLogger(log.New(fd, "\n", 0))

	return db, nil

}

// 检测目录或文件是否存在
func isExist(name string) bool {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return false
	}

	return true
}
