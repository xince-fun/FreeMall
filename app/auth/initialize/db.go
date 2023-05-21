package initialize

import (
	"fmt"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/xince-fun/FreeMall/app/leaf/global"
	"github.com/xince-fun/FreeMall/pkg/consts"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/opentelemetry/logging/logrus"
	"gorm.io/plugin/opentelemetry/tracing"
	"time"
)

// InitDB init DB
func InitDB() *gorm.DB {
	c := global.GlobalServerConfig.MysqlConfig
	dsn := fmt.Sprintf(consts.MySqlDSN, c.User, c.Password, c.Host, c.Port, c.DbName)

	newLogger := logger.New(
		logrus.NewWriter(), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // Slow SQL Threshold
			LogLevel:      logger.Silent, // Log level
			Colorful:      true,          // Disable color printing
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})
	if err != nil {
		klog.Fatalf("init db failed, err: %v", err)
	}
	klog.Infof("init db success")

	if err = db.Use(tracing.NewPlugin()); err != nil {
		klog.Fatalf("use tracing plugin failed, err: %v", err)
	}
	return db
}
