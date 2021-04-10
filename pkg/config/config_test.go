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
* @Modified Apr 09 2021
**/

package config

import (
	"os"
	"testing"

	sysadmLogger "github.com/wangyysde/bzhyserver/pkg/logger"
)

func Test_config(t *testing.T) {

	testConfig := New()
	err := testConfig.ParseConfig("./sysadm.yaml")
	testLogger := sysadmLogger.New()
	testLogger.LoggerFormat = "text"
	testLogger.InitStdoutLogger()
	defer testLogger.EndLogger("stdout")

	if err != nil {
		testLogger.LoggingLog("stdout", "error", err)
	}
  
  err = testConfig.CheckConfig()
  if err != nil {
    testLogger.LoggingLog("stdout", "error", err)  
  }

	os.Exit(0)

}
