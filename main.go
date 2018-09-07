package main

import (
	"flag"
	"fmt"

	"github.com/JacobDenver007/UserSystemL0/user"
	gin "gopkg.in/gin-gonic/gin.v1"
)

func main() {
	listenport := flag.Int("listenport", 9000, "api listen port")

	router := gin.Default()
	router.Use(user.Checklogged())
	user.RegisterAPI(router)
	if err := router.Run(fmt.Sprintf(":%d", *listenport)); err != nil {
		panic(err)
	}
	user.Init()

}
