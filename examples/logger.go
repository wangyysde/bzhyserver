package main

import "github.com/wangyysde/bzhyserver"

func main() {
	bzhyserver.SetErrorFile("/var/log/testError.log")
    defer bzhyserver.CloseErrLogFd()
	bzhyserver.SetAccessFile("/var/log/testAccess.log")
    defer bzhyserver.CloseAccLogFd()
    r := bzhyserver.Default()
	r.GET("/ping", func(c *bzhyserver.Context) {
		c.JSON(200, bzhyserver.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
