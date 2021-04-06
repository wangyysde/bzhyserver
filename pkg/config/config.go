/*
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
* @License GNU Lesser General Public License  https://www.sysadm.cn/lgpl.html
* @Modified Apr 06 2021
**/

package config

type defaultConfigs struct {
	Progname       string //The software name
	Proversion     string //The version of the release
	DefaultPrefix  string //The path of the software will be installed
	DefaultConFile string //The path of the configure file for the software
}

var DefaultConfigs = defaultConfigs{
	Progname:       "sysadm",
	Proversion:     "0.21.03",
	DefaultPrefix:  "/usr/local/sysadm",
	DefaultConFile: "/usr/local/sysadm/config/sysadm.yaml",
}
