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
*       @License GNU Lesser General Public License  https://www.sysadm.cn/lgpl.html
* @Modified Apr 02 2021
**/

package logger

import (
	"github.com/wangyysde/bzhyserver/logger"
	"os"

	log "github.com/sirupsen/logrus"
)

func Test_logger(){
	
	var Logger *log.Logger
	var fp *os.File
	var err error

	logger.InitStdoutLogger()
	defer EndLogger("stdout")
	logger.LoggingLog("stdout","info","This message output to stdout")

	fp, err = logger.ConfigLogfile("access","/var/log/bzhy_access.log")
	if err != nil {
		logger.LoggingLog("stdout","error",err)
	} else {
		defer EndLogger("access")
	}
	
	logger.LoggingLog("access","info","This message output to access")
	logger.Config.allstdout = false

	logger.LoggingLog("access","info","This message output to access at allstdout is false")
	logger.LoggingLog("stdout","info","This message output to stdout at allstdout is false")

	fp, err = logger.ConfigLogfile("access","/var/log/bzhy_error.log")
	if err != nil {
		logger.LoggingLog("stdout","error",err)
	} else {
		defer EndLogger("error")
	}
	
	logger.Config.allstdout = true
	logger.LoggingLog("error","error","This message output to access at allstdout is true")
	logger.LoggingLog("error","error","This message output to stdout at allstdout is true")

	logger.Config.allstdout = false
	logger.LoggingLog("error","error","This message output to access at allstdout is false")
	logger.LoggingLog("error","error","This message output to stdout at allstdout is false")

  os.Exit(0)
}
