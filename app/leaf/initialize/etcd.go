package initialize

import (
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/xince-fun/FreeMall/app/leaf/global"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func InitEtcd() (cli *clientv3.Client) {

	if global.GlobalServerConfig.EtcdConfig.SnowflakeEnable || global.GlobalServerConfig.EtcdConfig.DiscoveryEnable {
		var err error

		cli, err = clientv3.New(clientv3.Config{
			Endpoints:   global.GlobalServerConfig.EtcdConfig.Endpoints,
			DialTimeout: time.Duration(global.GlobalServerConfig.EtcdConfig.DialTimeout) * time.Second,
		})
		if err != nil {
			klog.Fatal("init etcd client failed, err: %v", err)
		}
	}
	klog.Infof("init etcd client success")
	return cli
}
