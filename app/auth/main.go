package main

import (
	"log"

	auth "github.com/xince-fun/FreeMall/kitex_gen/auth/authservice"
)

func main() {
	svr := auth.NewServer(new(AuthServiceImpl))

	err := svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
