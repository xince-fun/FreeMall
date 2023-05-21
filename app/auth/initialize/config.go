package initialize

import (
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/spf13/viper"
	"github.com/xince-fun/FreeMall/app/leaf/global"
	"github.com/xince-fun/FreeMall/pkg/consts"
)

// InitConfig init config
func InitConfig() {
	v := viper.New()
	v.SetConfigFile(consts.AuthConfigPath)

	if err := v.ReadInConfig(); err != nil {
		klog.Fatalf("read config file failed, err: %v", err)
	}
	if err := v.Unmarshal(&global.GlobalServerConfig); err != nil {
		klog.Fatalf("unmarshal config file failed, err: %v", err)
	}

	klog.Infof("config: %+v", global.GlobalServerConfig)

}
