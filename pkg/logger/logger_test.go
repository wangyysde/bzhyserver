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
	"os"
	"testing"
)

func Test_logger(t *testing.T) {

	sysadmLogger := New()
	sysadmLogger.LoggerFormat = "text"
	sysadmLogger.InitStdoutLogger()
	defer sysadmLogger.EndLogger("stdout")
	sysadmLogger.LoggingLog("stdout", "info", "This message output to stdout")

	_, err := sysadmLogger.OpenLogfile("access", "/var/log/bzhy_access.log")
	if err != nil {
		sysadmLogger.LoggingLog("stdout", "error", err)
	} else {
		defer sysadmLogger.EndLogger("access")
	}

	sysadmLogger.LoggingLog("access", "info", "This message output to access")
	sysadmLogger.Allstdout = false

	sysadmLogger.LoggingLog("access", "info", "This message output to access at allstdout is false")
	sysadmLogger.LoggingLog("stdout", "info", "This message output to stdout at allstdout is false")

	_, err = sysadmLogger.OpenLogfile("error", "/var/log/bzhy_error.log")
	if err != nil {
		sysadmLogger.LoggingLog("stdout", "error", err)
	} else {
		defer sysadmLogger.EndLogger("error")
	}

	sysadmLogger.Allstdout = true
	sysadmLogger.LoggingLog("error", "error", "This message output to access at allstdout is true")
	sysadmLogger.LoggingLog("error", "error", "This message output to stdout at allstdout is true")

	sysadmLogger.Allstdout = false
	sysadmLogger.LoggingLog("error", "error", "This message output to access at allstdout is false")
	sysadmLogger.LoggingLog("error", "error", "This message output to stdout at allstdout is false")

	os.Exit(0)
}
