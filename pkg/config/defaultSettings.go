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
**/

//NOTE: This file be generated by building scripts. DON'T EDIT THIS FILE.
//Generated at: 20210412 12:05:51

package config

/* DefaultConfigs for default configs of application */
var DefaultAppSettings = appSetting{
    Progname:   "sysadm",
    Proversion: "0.21.3",
    Prefix:     "/usr/local/sysadm",
    ConFile:    "/usr/local/sysadm",
}

//Define default value for server settings
var defaultServerSettings = server{
    Listen:   "0.0.0.0",
    Port:     8080,
    RootPath: "html",
    PidPath:  "/var/run/sysadm.pid",
    Indexs:   "index.html index.htm",
}

var defaultLoggerSettings = logger{
    Loglevel:  "debug",
    AccessLog: "logs/sysadm-access.log",
    ErrorLog:  "logs/sysadm-error.log",
    Logtype:   "text",
}

