package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)


func main() {
	eg := gin.Default()
	eg.GET("/", func (c *gin.Context){
		fmt.Fprint(c.Writer, "HELLO")
	})
	eg.Run(":8080")
}
