package main

import (
	client "gchat/gclient"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func onSetRouter(r *gin.Engine,ws *client.WsServer){
	r.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer , c.Request ,"index.html")
	})
	r.GET("/ws", func(c *gin.Context) {
		client.OnWsServer(c.Writer , c.Request,ws)
	})
}

func main()  {
	//
	router := gin.Default()
	ws := client.OnNewWsServer()
	go ws.OnRun()//loop check event
	//
	onSetRouter(router,ws)
	log.Fatal(router.Run(":3000"))
}