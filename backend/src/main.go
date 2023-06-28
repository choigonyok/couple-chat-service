package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	eg := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost"} 
	// 허용할 오리진 설정, 원래 리액트의 port가 아니라 리액트가 있는 container의 port 번호를 origin allow 해줘야함
	// localhost:3000로 origin allow 하면 통신 안됨

	config.AllowMethods= []string{"GET"}
	config.AllowHeaders = []string{"Content-type"}
	config.AllowCredentials = true
	eg.Use(cors.New(config)) 
	// origin 설정하고 설정한 config를 gin engine에서 사용하겠다는 이 부분이 있어야 적용이 됨!

	eg.GET("/api/test", func (c *gin.Context){
		c.Writer.WriteHeader(200)
		c.Writer.Write([]byte("HELLO!"))
	})
	eg.Run(":8080")
}
