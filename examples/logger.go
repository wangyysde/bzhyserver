package main

import "github.com/wangyysde/bzhyserver"

func main() {
	r := bzhyserver.Default()
	r.GET("/ping", func(c *bzhyserver.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
