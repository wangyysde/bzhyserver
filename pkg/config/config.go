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

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"
	"strings"

	sysadmlog "github.com/wangyysde/bzhyserver/pkg/logger"
	"gopkg.in/yaml.v2"
)

/* define default configs struct of application */
type appSetting struct {
	Progname   string //The software name
	Proversion string //The version of the release
	Prefix     string //The path of the software will be installed
	ConFile    string //The path of the configure file for the software
}

//Struct for parsing server block of server file
type server struct {
	Listen   string `yaml: "listen"`
	Port     int    `yaml: "port"`
	RootPath string `yaml: "root"`
	PidPath  string `yaml: "pid"`
	Indexs   string `yaml: "index"`
}

//Struct for log block of config file
type logger struct {
	Loglevel  string `yaml: "loglevel"`
	AccessLog string `yaml: "accesslog"`
	ErrorLog  string `yaml: "errorlog"`
	Logtype   string `yaml: "logtype"`
}

//Struct for runtime settings
type runtime struct {
	Logger *sysadmlog.SysadmLogger //for Logger
}

//Struct for application configuration
type Configs struct {
	App     appSetting //for application
	Server  server     //for server block
	Logger  logger     //for log block
	Runtime runtime
}

//Initate a variable for application configuration
var Settings = Configs{
	App:     DefaultAppSettings,
	Server:  server{},
	Logger:  logger{},
	Runtime: runtime{},
}

//Create  a instance
func New() *Configs {
	return &Settings
}

//Read the values from confFile and put them into settings
func (settings *Configs) ParseConfig(confFile string) (err error) {
	err = nil
	if len(confFile) <= 0 {
		confFile = settings.App.ConFile
	} else {
		settings.App.ConFile = confFile
	}

	yamlFile, err := ioutil.ReadFile(confFile)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, settings)

	return err
}

func (settings *Configs) CheckConfig() (err error) {
	err = nil
	var fp *os.File

	//Checking whether the listen ip is valid
	if len(settings.Server.Listen) > 0 {
		if net.ParseIP(settings.Server.Listen) == nil {
			err = fmt.Errorf("The Listen address:%s is invalid", settings.Server.Listen)
			return err
		}
	} else {
		settings.Server.Listen = defaultServerSettings.Listen
	}

	//If there isn't port in the Yaml file
	if settings.Server.Port == 0 {
		settings.Server.Port = defaultServerSettings.Port
	}

	if settings.Server.Port <= 1024 || settings.Server.Port > 65535 {
		err = fmt.Errorf("The Port:%d  is invalid. The valid port number is between 1024 and 65535!", settings.Server.Port)
		return err
	}

	if len(settings.Server.RootPath) == 0 {
		settings.Server.RootPath = defaultServerSettings.RootPath
	}
	if !path.IsAbs(settings.Server.RootPath) {
		settings.Server.RootPath = path.Join(DefaultAppSettings.Prefix, settings.Server.RootPath)
	}

	if !path.IsAbs(settings.Server.RootPath) {
		err = fmt.Errorf("The path of static file:%s is invalid ", settings.Server.RootPath)
		return err
	}

	if len(settings.Server.PidPath) == 0 {
		settings.Server.PidPath = defaultServerSettings.PidPath
	}

	if !path.IsAbs(settings.Server.PidPath) {
		settings.Server.PidPath = path.Join(DefaultAppSettings.Prefix, settings.Server.PidPath)
	}

	fp, err = os.OpenFile(settings.Server.PidPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	fp.Close()

	if len(settings.Server.Indexs) == 0 {
		settings.Server.Indexs = defaultServerSettings.Indexs
	}

	if len(settings.Logger.Loglevel) == 0 {
		settings.Logger.Loglevel = defaultLoggerSettings.Loglevel
	}

	found := -1
	for i := 0; i < len(sysadmlog.LevelList); i++ {
		if strings.ToLower(settings.Logger.Loglevel) == sysadmlog.LevelList[i] {
			found = i
			break
		}
	}

	if found < 0 {
		err = fmt.Errorf("The loglevel:%s is invalid", settings.Logger.Loglevel)
		return err
	}

	if len(settings.Logger.AccessLog) == 0 {
		settings.Logger.AccessLog = defaultLoggerSettings.AccessLog
	}

	if !path.IsAbs(settings.Logger.AccessLog) {
		settings.Logger.AccessLog = path.Join(DefaultAppSettings.Prefix, settings.Logger.AccessLog)
	}

	fp, err = os.OpenFile(settings.Logger.AccessLog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	fp.Close()

	if len(settings.Logger.ErrorLog) == 0 {
		settings.Logger.ErrorLog = defaultLoggerSettings.ErrorLog
	}

	if !path.IsAbs(settings.Logger.ErrorLog) {
		settings.Logger.ErrorLog = path.Join(DefaultAppSettings.Prefix, settings.Logger.ErrorLog)
	}

	fp, err = os.OpenFile(settings.Logger.ErrorLog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	fp.Close()

	if len(settings.Logger.Logtype) < 1 {
		settings.Logger.Logtype = defaultLoggerSettings.Logtype
	}

	if strings.ToLower(settings.Logger.Logtype) != "text" && strings.ToLower(settings.Logger.Logtype) != "json" {
		err = fmt.Errorf("The logType:%s is invalid", settings.Logger.Logtype)
		return err
	}

	return err

}
