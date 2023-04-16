package main

import (
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/xince-fun/FreeMall/app/leaf/dao"
	"github.com/xince-fun/FreeMall/app/leaf/global"
	"github.com/xince-fun/FreeMall/app/leaf/initialize"
	"github.com/xince-fun/FreeMall/app/leaf/service"
	leaf "github.com/xince-fun/FreeMall/kitex_gen/leaf/leafservice"
	"github.com/xince-fun/FreeMall/pkg/consts"
	"github.com/xince-fun/FreeMall/pkg/logger"
	"log"
	"net"
	"strconv"
)

func main() {
	r, err := etcd.NewEtcdRegistry([]string{"127.0.0.1:2379"})
	if err != nil {
		panic(err)
	}
	initialize.InitConfig()
	logger.InitKLogger(consts.KlogFilePath, global.GlobalServerConfig.LogLevel)
	db := initialize.InitDB()
	cli := initialize.InitEtcd()
	segmentIdGenRepo := dao.NewSegmentIdGenRepo(db, global.GlobalServerConfig.MysqlConfig.TableName)
	snowflakeIdGenRepo := dao.NewSnowflakeIdGenRepoImpl(cli)
	addr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(global.GlobalServerConfig.Host, strconv.Itoa(global.GlobalServerConfig.Port)))
	if err != nil {
		panic(err)
	}
	svr := leaf.NewServer(&LeafServiceImpl{
		segmentIdGenUseCase:   service.NewSegmentIdGenUseCase(segmentIdGenRepo),
		snowflakeIdGenUseCase: service.NewSnowflakeIdGenUseCase(snowflakeIdGenRepo),
	},
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "LeafService"}),
	)

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
