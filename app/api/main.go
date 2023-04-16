// Code generated by hertz generator.

package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/xince-fun/FreeMall/app/api/rpc"
)

func main() {
	rpc.Init()
	h := server.Default()

	register(h)
	h.Spin()
}