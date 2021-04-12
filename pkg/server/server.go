/**
* SYSADM Server
* @Author  Wayne Wang <net_use@bzhy.com>
* @Copyright Bzhy Network
* @HomePage http://www.sysadm.cn
* @Version 0.21.03
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
* http://www.apache.org/licenses/LICENSE-2.0
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*	@License GNU Lesser General Public License  https://www.sysadm.cn/lgpl.html
* @Modified Mar 25 2021
**/

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	kingpin "gopkg.in/alecthomas/kingpin.v2" //https://github.com/alecthomas/kingpin    https://gopkg.in/alecthomas/kingpin.v2

	"github.com/wangyysde/bzhyserver/pkg/config"
	"github.com/wangyysde/bzhyserver/pkg/logger"
)

type Server struct {
	icontext   context.Context
	shutdownFn context.CancelFunc

	shutdownReason     string
	shutdownInProgress bool

	rootPath string

	nindexs []string

	index   string
	pidFile string

	r *gin.Engine
}

var (
	a          = kingpin.New(filepath.Base(os.Args[0]), "A command-line "+config.DefaultAppSettings.Progname+" application.")
	configFile = a.Flag("config", "Configuration file path").Default(config.DefaultAppSettings.ConFile).String()
	version    = a.Flag("version", "Show the version information for "+config.DefaultAppSettings.Progname).Bool()
)

var Svr = new(Server)

func init_serer() (ret int) {
	r := gin.Default()
	//	r.SetAccLogHandler(WriteLog2Acclog)
	//	r.SetErrLogHandler(WriteLog2Errlog)

	/*
		r.GET("/", GetHandler)
		r.POST("/somePost", posting)
		r.PUT("/somePut", putting)
		r.DELETE("/someDelete", deleting)
		r.PATCH("/somePatch", patching)
		r.HEAD("/someHead", head)
		r.OPTIONS("/someOptions", options)

		gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
			logmsg := fmt.Sprintf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
		}

	*/

	r.GET("/", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "Welcome Gin Server")
		c.JSON(200, gin.H{
			"Blog":   "www.flysnow.org",
			"wechat": "flysnow_org",
		})
	})

	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(8081),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			//WriteStartLog(fmt.Sprintf("Listen %s:%s %s", config.Cfg.Host, config.Cfg.Port, err), "fatal")
		}
	}()

	return 0

}

func main() {
	sysadmLogger := logger.New()
	sysadmLogger.LoggerFormat = "text"
	sysadmLogger.InitStdoutLogger()
	defer sysadmLogger.EndLogger("stdout")

	a.HelpFlag.Short('h')
	_, err := a.Parse(os.Args[1:])
	if err != nil {
		sysadmLogger.LoggingLogf("stdout", "info", "Unkown  commandline arguments:%s", err)
		os.Exit(10001) //Error no: AABBB. AA: file seq,main is 1; BBB: error no
	}

	if ret := init_serer(); ret > 0 {
		fmt.Printf("Starting the server ERROR")
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Printf("Shutting down server...")
	//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	//defer cancel()

	os.Exit(0)
}
