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
* @Modified Mar 29 2021
**/

package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

/*
configLogger struct used to save parameters for logger, such as access log file error log file
*/

type SysadmLogger struct {
	//accessLoggerFile is the path of access log file ,if logger access log to a file
	accessLoggerFile string
	//errorLoggerFile is the path of error log file ,if logger error log to a file
	errorLoggerFile string

	//The file descriptor of the access log file
	accessFp *os.File
	//The file descriptor of the error log file
	errorFp *os.File

	//Logger for access log
	accessLogger *log.Logger
	//Logger for error log
	errorLogger *log.Logger
	//Logger for stdout
	stdoutLogger *log.Logger

	//set log format for output
	LoggerFormat string
	//set date formate
	DateFormat string

	//If all log message log to stdout ,then Allstdout should be set to True
	Allstdout bool
}

var levelList = [7]string{"panic", "fatal", "error", "warn", "info", "debug", "trace"}

//Set global variable config and its default value
var sysadmLogger = SysadmLogger{
	accessLoggerFile: "",
	errorLoggerFile:  "",

	accessFp: nil,
	errorFp:  nil,

	accessLogger: nil,
	errorLogger:  nil,
	stdoutLogger: nil,

	LoggerFormat: "Text",
	DateFormat:   time.RFC3339, //Ref: https://studygolang.com/static/pkgdoc/pkg/time.htm#Time.Format
	Allstdout:    true,
}

func New() *SysadmLogger {
	return &sysadmLogger
}

/*
* initated a logger to logging log message to stdout
 */
func (sysadmLogger *SysadmLogger) InitStdoutLogger() (stdoutLogger *log.Logger, err error) {
	stdoutLogger = log.New()
	stdoutLogger.Out = os.Stdout
	stdoutLogger = sysadmLogger.SetLogFormat(stdoutLogger, "stdout")

	sysadmLogger.stdoutLogger = stdoutLogger
	if sysadmLogger.accessLogger == nil && sysadmLogger.errorLogger == nil {
		sysadmLogger.Allstdout = true
	}

	return stdoutLogger, nil

}

/*
* Init logger enity for access or error.
* before call this func, sysadmLoggerLogfile(logType, logFile) should be called
 */
func (sysadmLogger *SysadmLogger) InitLogger(logType string, toStdout bool) (logger *log.Logger, err error) {
	err = nil
	if strings.ToLower(logType) != "access" && strings.ToLower(logType) != "error" {
		err = fmt.Errorf("LogType must be access or error.You input is: %s", logType)
		return nil, err
	}

	if strings.ToLower(logType) == "access" {
		if sysadmLogger.accessFp == nil {
			err = fmt.Errorf("May be not set access log, you should call sysadmLoggerLogfile(%s, logFile) before call InitLogger", logType)
			return nil, err
		}
		logger = log.New()
		logger.Out = sysadmLogger.accessFp
		logger = sysadmLogger.SetLogFormat(logger, logType)
		sysadmLogger.accessLogger = logger
		if toStdout {
			sysadmLogger.Allstdout = true
		} else {
			sysadmLogger.Allstdout = false
		}
		return logger, nil
	}

	if sysadmLogger.errorFp == nil {
		err = fmt.Errorf("May be not set error log, you should call sysadmLoggerLogfile(%s, logFile) before call InitLogger", logType)
		return nil, err
	}

	logger = log.New()
	logger.Out = sysadmLogger.errorFp
	logger = sysadmLogger.SetLogFormat(logger, logType)
	sysadmLogger.errorLogger = logger

	return logger, nil
}

/*
* EndLogger function will be call by defer
* EndLogger will close file descriptor of access or error
* and reset logger of access or error to nil
 */
func (sysadmLogger *SysadmLogger) EndLogger(logType string) (err error) {
	err = nil
	var fp *os.File

	switch strings.ToLower(logType) {
	case "access":
		fp = sysadmLogger.accessFp
		if fp != nil {
			err = fp.Close()
		} else {
			err = fmt.Errorf("Access logger have closed")
		}

		if err == nil {
			sysadmLogger.accessFp = nil
			sysadmLogger.accessLoggerFile = ""
			sysadmLogger.accessLogger = nil
			if sysadmLogger.accessLogger == nil && sysadmLogger.errorLogger == nil {
				sysadmLogger.Allstdout = true
			}
		}
		break
	case "error":
		fp = sysadmLogger.errorFp
		if fp != nil {
			err = fp.Close()
		} else {
			err = fmt.Errorf("Error logger have closed")
		}

		if err == nil {
			sysadmLogger.errorFp = nil
			sysadmLogger.errorLogger = nil
			sysadmLogger.errorLoggerFile = ""
			if sysadmLogger.accessLogger == nil && sysadmLogger.errorLogger == nil {
				sysadmLogger.Allstdout = true
			}
		}
		break
	case "stdout":
		if sysadmLogger.accessLogger != nil || sysadmLogger.errorLogger != nil {
			err = fmt.Errorf("Access logger and Error logger should be end first")
		}

		if err == nil {
			sysadmLogger.stdoutLogger = nil
		}
		break
	default:
		err = fmt.Errorf("logType: %s is invalid", logType)
		break
	}

	return err
}

/*
* set log formate to text or json and set loggerFormat
* the fields of the struct of loggerFormat refer to :https://pkg.go.dev/github.com/sirupsen/logrus#JSONFormatter
 */
func (sysadmLogger *SysadmLogger) SetLogFormat(Logger *log.Logger, logType string) (logger *log.Logger) {
	if strings.ToLower(logType) == "access" || strings.ToLower(logType) == "error" {
		if strings.ToLower(sysadmLogger.LoggerFormat) == "text" {
			Logger.SetFormatter(&log.TextFormatter{
				ForceColors:               false, //Ref: https://pkg.go.dev/github.com/sirupsen/logrus#pkg-functions
				DisableColors:             true,
				ForceQuote:                false,
				DisableQuote:              true,
				EnvironmentOverrideColors: true,
				DisableTimestamp:          false,
				FullTimestamp:             true,
				TimestampFormat:           sysadmLogger.DateFormat,
				DisableSorting:            true,
				DisableLevelTruncation:    true,
				PadLevelText:              true,
			})
		} else {
			Logger.SetFormatter(&log.JSONFormatter{
				TimestampFormat:  sysadmLogger.DateFormat,
				DisableTimestamp: false,
			})
		}
	} else {
		if strings.ToLower(sysadmLogger.LoggerFormat) == "text" {
			Logger.SetFormatter(&log.TextFormatter{
				ForceColors:               true, //Ref: https://pkg.go.dev/github.com/sirupsen/logrus#pkg-functions
				DisableColors:             false,
				ForceQuote:                true,
				DisableQuote:              false,
				EnvironmentOverrideColors: true,
				DisableTimestamp:          false,
				FullTimestamp:             true,
				TimestampFormat:           sysadmLogger.DateFormat,
				DisableSorting:            true,
				DisableLevelTruncation:    true,
				PadLevelText:              true,
			})
		} else {
			Logger.SetFormatter(&log.JSONFormatter{
				TimestampFormat:  sysadmLogger.DateFormat,
				DisableTimestamp: false,
			})
		}
	}

	return Logger

}

/*
* set logger level to sysadmLogger.loggerLevel
 */
func (sysadmLogger *SysadmLogger) SetLoglevel(loggerLevel string, Logger *log.Logger) (logger *log.Logger) {

	switch strings.ToLower(loggerLevel) {
	case "panic":
		Logger.SetLevel(log.PanicLevel)
		break
	case "fatal":
		Logger.SetLevel(log.FatalLevel)
		break
	case "error":
		Logger.SetLevel(log.ErrorLevel)
		break
	case "warn":
		Logger.SetLevel(log.WarnLevel)
		break
	case "info":
		Logger.SetLevel(log.InfoLevel)
		break
	case "debug":
		Logger.SetLevel(log.DebugLevel)
		break
	case "trace":
		Logger.SetLevel(log.TraceLevel)
		break
	default:
		Logger.SetLevel(log.DebugLevel)
	}

	return Logger
}

/*
* according to logType, sysadmLoggerLogfile set logFile to sysadmLogger.accessLoggerFile or sysadmLogger.errorLoggerFile
* and set file descriptor to accessFp or errorFp if logFile can be opened.
* to close the openned file on time, a defer function should be called following call this function if this return successful.
 */
func (sysadmLogger *SysadmLogger) OpenLogfile(logType string, logFile string) (fp *os.File, err error) {

	err = nil
	if strings.ToLower(logType) != "access" && strings.ToLower(logType) != "error" {
		err = fmt.Errorf("LogType must be access or error.You input is: %s", logType)
		return nil, err
	}

	fp, err = os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		err = fmt.Errorf("Open %s log file %s error: %s", logType, logFile, fmt.Sprintf("%s", err))
		return nil, err
	}

	if strings.ToLower(logType) == "access" {
		sysadmLogger.accessFp = fp
		sysadmLogger.accessLoggerFile = logFile
		_, err = sysadmLogger.InitLogger("access", sysadmLogger.Allstdout)
	} else {
		sysadmLogger.errorFp = fp
		sysadmLogger.errorLoggerFile = logFile
		_, err = sysadmLogger.InitLogger("error", sysadmLogger.Allstdout)
	}

	return fp, err
}

/**
* Logging a message to Logger
* if the sysadmLogger.Allstdout ,then logging the log messages to stdout
 */
func (sysadmLogger *SysadmLogger) LoggingLog(logType string, logLevel string, args ...interface{}) {

	var logger *log.Logger
	var tostdout bool
	var stdLogger *log.Logger

	tostdout = sysadmLogger.Allstdout
	logger = nil
	stdLogger = nil

	switch strings.ToLower(logType) {
	case "access":
		logger = sysadmLogger.accessLogger
		break
	case "error":
		logger = sysadmLogger.errorLogger
		break
	case "stdout":
		tostdout = false
		logger = sysadmLogger.stdoutLogger
	default:
		tostdout = false
		logger = sysadmLogger.stdoutLogger
	}

	if tostdout {
		stdLogger = sysadmLogger.stdoutLogger
	}

	found := -1
	for i := 0; i < len(logLevel); i++ {
		if strings.ToLower(logLevel) == levelList[i] {
			found = i
			break
		}
	}

	if found < 0 {
		logLevel = "debug"
	}

	if logger != nil {
		logger = sysadmLogger.SetLoglevel(logLevel, logger)
	}

	if stdLogger != nil {
		stdLogger = sysadmLogger.SetLoglevel(logLevel, stdLogger)
	}
	switch strings.ToLower(logLevel) {
	case "panic":
		if logger != nil {
			logger.Panic(args...)
		}
		if stdLogger != nil {
			stdLogger.Panic(args...)
		}
		break
	case "fatal":
		if logger != nil {
			logger.Fatal(args...)
		}
		if stdLogger != nil {
			stdLogger.Fatal(args...)
		}
		break
	case "error":
		if logger != nil {
			logger.Error(args...)
		}
		if stdLogger != nil {
			stdLogger.Error(args...)
		}
		break
	case "warn":
		if logger != nil {
			logger.Warn(args...)
		}
		if stdLogger != nil {
			stdLogger.Warn(args...)
		}
		break
	case "info":
		if logger != nil {
			logger.Info(args...)
		}
		if stdLogger != nil {
			stdLogger.Info(args...)
		}
		break
	case "debug":
		if logger != nil {
			logger.Debug(args...)
		}
		if stdLogger != nil {
			stdLogger.Debug(args...)
		}
		break
	case "trace":
		if logger != nil {
			logger.Trace(args...)
		}
		if stdLogger != nil {
			stdLogger.Trace(args...)
		}
		break
	}
}

/*
 * Logging a message to Logger
 * if the sysadmLogger.Allstdout ,then logging the log messages to stdout
 */
func (sysadmLogger *SysadmLogger) LoggingLogf(logType string, logLevel string, format string, args ...interface{}) {

	var logger *log.Logger
	var tostdout bool
	var stdLogger *log.Logger

	tostdout = sysadmLogger.Allstdout
	logger = nil
	stdLogger = nil

	switch strings.ToLower(logType) {
	case "access":
		logger = sysadmLogger.accessLogger
		break
	case "error":
		logger = sysadmLogger.errorLogger
		break
	case "stdout":
		tostdout = false
		logger = sysadmLogger.stdoutLogger
	default:
		tostdout = false
		logger = sysadmLogger.stdoutLogger
	}

	if tostdout {
		stdLogger = sysadmLogger.stdoutLogger
	}

	found := -1
	for i := 0; i < len(levelList); i++ {
		if strings.ToLower(logLevel) == levelList[i] {
			found = i
			break
		}
	}

	if found < 0 {
		logLevel = "debug"
	}

	if logger != nil {
		logger = sysadmLogger.SetLoglevel(logLevel, logger)
		switch strings.ToLower(logLevel) {
		case "panic":
			logger.Panicf(format, args...)
			break
		case "fatal":
			logger.Fatalf(format, args...)
			break
		case "error":
			logger.Errorf(format, args...)
			break
		case "warn":
			logger.Warnf(format, args...)
			break
		case "info":
			logger.Infof(format, args...)
			break
		case "debug":
			logger.Debugf(format, args...)
			break
		case "trace":
			logger.Tracef(format, args...)
			break
		}
	}

	if stdLogger != nil {
		stdLogger = sysadmLogger.SetLoglevel(logLevel, stdLogger)
		switch strings.ToLower(logLevel) {
		case "panic":
			stdLogger.Panicf(format, args...)
			break
		case "fatal":
			stdLogger.Fatalf(format, args...)
			break
		case "error":
			stdLogger.Errorf(format, args...)
			break
		case "warn":
			stdLogger.Warnf(format, args...)
			break
		case "info":
			stdLogger.Infof(format, args...)
			break
		case "debug":
			stdLogger.Debugf(format, args...)
			break
		case "trace":
			stdLogger.Tracef(format, args...)
			break
		}
	}
}

/*
 * Logging a message to Logger
 * if the sysadmLogger.Allstdout ,then logging the log messages to stdout
 */
func (sysadmLogger *SysadmLogger) LoggingLogln(logType string, logLevel string, args ...interface{}) {

	var logger *log.Logger
	var tostdout bool
	var stdLogger *log.Logger

	tostdout = sysadmLogger.Allstdout
	logger = nil
	stdLogger = nil

	switch strings.ToLower(logType) {
	case "access":
		logger = sysadmLogger.accessLogger
		break
	case "error":
		logger = sysadmLogger.errorLogger
		break
	case "stdout":
		tostdout = false
		logger = sysadmLogger.stdoutLogger
	default:
		tostdout = false
		logger = sysadmLogger.stdoutLogger
	}

	if tostdout {
		stdLogger = sysadmLogger.stdoutLogger
	}

	found := -1
	for i := 0; i < len(levelList); i++ {
		if strings.ToLower(logLevel) == levelList[i] {
			found = i
			break
		}
	}

	if found < 0 {
		logLevel = "debug"
	}
	if logger != nil {
		logger = sysadmLogger.SetLoglevel(logLevel, logger)
		switch strings.ToLower(logLevel) {
		case "panic":
			logger.Panicln(args...)
			break
		case "fatal":
			logger.Fatalln(args...)
			break
		case "error":
			logger.Errorln(args...)
			break
		case "warn":
			logger.Warnln(args...)
			break
		case "info":
			logger.Infoln(args...)
			break
		case "debug":
			logger.Debugln(args...)
			break
		case "trace":
			logger.Traceln(args...)
			break
		}
	}

	if stdLogger != nil {
		stdLogger = sysadmLogger.SetLoglevel(logLevel, stdLogger)
		switch strings.ToLower(logLevel) {
		case "panic":
			stdLogger.Panicln(args...)
			break
		case "fatal":
			stdLogger.Fatalln(args...)
			break
		case "error":
			stdLogger.Errorln(args...)
			break
		case "warn":
			stdLogger.Warnln(args...)
			break
		case "info":
			stdLogger.Infoln(args...)
			break
		case "debug":
			stdLogger.Debugln(args...)
			break
		case "trace":
			stdLogger.Traceln(args...)
			break
		}
	}
}

/*
 * Logging a message to Logger
 * if the sysadmLogger.Allstdout ,then logging the log messages to stdout
 */
func (sysadmLogger *SysadmLogger) LoggingLogFn(logType string, logLevel string, fn log.LogFunction) {

	var logger *log.Logger
	var tostdout bool
	var stdLogger *log.Logger

	tostdout = sysadmLogger.Allstdout
	logger = nil
	stdLogger = nil

	switch strings.ToLower(logType) {
	case "access":
		logger = sysadmLogger.accessLogger
		break
	case "error":
		logger = sysadmLogger.errorLogger
		break
	case "stdout":
		tostdout = false
		logger = sysadmLogger.stdoutLogger
	default:
		tostdout = false
		logger = sysadmLogger.stdoutLogger
	}

	if tostdout {
		stdLogger = sysadmLogger.stdoutLogger
	}

	found := -1
	for i := 0; i < len(levelList); i++ {
		if strings.ToLower(logLevel) == levelList[i] {
			found = i
			break
		}
	}

	if found < 0 {
		logLevel = "debug"
	}
	if logger != nil {
		logger = sysadmLogger.SetLoglevel(logLevel, logger)
		switch strings.ToLower(logLevel) {
		case "panic":
			logger.PanicFn(fn)
			break
		case "fatal":
			logger.FatalFn(fn)
			break
		case "error":
			logger.ErrorFn(fn)
			break
		case "warn":
			logger.WarnFn(fn)
			break
		case "info":
			logger.InfoFn(fn)
			break
		case "debug":
			logger.DebugFn(fn)
			break
		case "trace":
			logger.TraceFn(fn)
			break
		}
	}

	if stdLogger != nil {
		stdLogger = sysadmLogger.SetLoglevel(logLevel, stdLogger)
		switch strings.ToLower(logLevel) {
		case "panic":
			stdLogger.PanicFn(fn)
			break
		case "fatal":
			stdLogger.FatalFn(fn)
			break
		case "error":
			stdLogger.ErrorFn(fn)
			break
		case "warn":
			stdLogger.WarnFn(fn)
			break
		case "info":
			stdLogger.InfoFn(fn)
			break
		case "debug":
			stdLogger.DebugFn(fn)
			break
		case "trace":
			stdLogger.TraceFn(fn)
			break
		}
	}
}
