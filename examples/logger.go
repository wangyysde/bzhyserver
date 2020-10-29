package main

import "github.com/wangyysde/bzhyserver"

func main() {
	bzhyserver.SetAccessFile("/var/log/testAccess.log")
	bzhyserver.SetErrorFile("/var/log/testError.log")
	r := bzhyserver.Default()
	r.GET("/ping", func(c *bzhyserver.Context) {
		c.JSON(200, bzhyserver.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
