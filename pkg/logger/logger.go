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

type ConfigLogger struct {
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
	loggerFormat string
	//set date formate
	dateFormat string

	//If all log message log to stdout ,then allstdout should be set to True
	allstdout bool
}

var ConfigLogLevel = [7]string{"panic", "fatal", "error", "warn", "info", "debug", "trace"}

//Set global variable config and its default value
var Config = ConfigLogger{
	accessLoggerFile: "",
	errorLoggerFile:  "",

	accessFp: nil,
	errorFp:  nil,

	accessLogger: nil,
	errorLogger:  nil,
	stdoutLogger: nil,

	loggerFormat: "Text",
	dateFormat:   time.RFC3339, //Ref: https://studygolang.com/static/pkgdoc/pkg/time.htm#Time.Format
	allstdout:     true,
}

/*
* initated a logger to logging log message to stdout
 */
func InitStdoutLogger() (stdoutLogger *log.Logger, err error) {
	stdoutLogger = log.New()
	stdoutLogger.Out = os.Stdout
	stdoutLogger = SetLogFormat(stdoutLogger)

	Config.stdoutLogger = stdoutLogger
	if Config.accessLogger == nil && Config.errorLogger == nil {
		Config.allstdout = true
	}

	return stdoutLogger, nil

}

/*
* Init logger enity for access or error.
* before call this func, ConfigLogfile(logType, logFile) should be called
 */
func InitLogger(logType string,toStdout bool) (logger *log.Logger, err error) {
	err = nil
	if strings.ToLower(logType) != "access" && strings.ToLower(logType) != "error" {
		err = fmt.Errorf("LogType must be access or error.You input is: %s", logType)
		return nil, err
	}

	if strings.ToLower(logType) == "access" {
		if Config.accessFp == nil {
			err = fmt.Errorf("May be not set access log, you should call ConfigLogfile(%s, logFile) before call InitLogger", logType)
			return nil, err
		}
		logger = log.New()
		logger.Out = Config.accessFp
		logger = SetLogFormat(logger)
		Config.accessLogger = logger
		if toStdout {
			Config.allstdout=true
		} else {
			Config.allstdout = false
		}
		return logger, nil
	}

	if Config.errorFp == nil {
		err = fmt.Errorf("May be not set error log, you should call ConfigLogfile(%s, logFile) before call InitLogger", logType)
		return nil, err
	}

	logger = log.New()
	logger.Out = Config.errorFp
	logger = SetLogFormat(logger)
	Config.errorLogger = logger

	return logger, nil
}


/*
* 
*/
func EndLogger(logType string) (err error){
	err = nil
	var fp *os.file

	switch strings.ToLower(logType) { 
		case "access":
			fp = Config.accessFp
			if fp != nil {
				err = fp.Close()
			} else {
				err = fmt.Errorf("Access logger have closed")
			}
		
			if err == nil {
				Config.accessFp = nil
				Config.accessLoggerFile = ""
				Config.accessLogger = nil
				if Config.accessLogger == nil && Config.errorLogger == nil {
					Config.allstdout = true
				}
			}
			break
		case "error":
			fp = Config.errorFp
			if fp != nil {
				err = fp.Close()
			} else {
				err = fmt.Errorf("Error logger have closed")
			}
			
			if err == nil {
				Config.errorFp = nil
				Config.errorLogger = nil
				Config.errorLoggerFile = ""
				if Config.accessLogger == nil && Config.errorLogger == nil {
					Config.allstdout = true
				}
			}
			break
		case "stdout":
			if Config.accessLogger != nil || Config.errorLogger != nil {
				err = fmt.Errorf("Access logger and Error logger should be end first")
			}

			if err == nil {
				Config.stdoutLogger = nil
			}
			break
		default:
			err = fmt.Errorf("logType: %s is invalid",logType)
			break
	}

	return err
}

/*
* set log formate to text or json and set loggerFormat
* the fields of the struct of loggerFormat refer to :https://pkg.go.dev/github.com/sirupsen/logrus#JSONFormatter
 */
func SetLogFormat(Logger *log.Logger) (logger *log.Logger) {
	if strings.ToLower(Config.loggerFormat) == "text" {
		Logger.SetFormatter(&log.TextFormatter{
			ForceColors:               true, //Ref: https://pkg.go.dev/github.com/sirupsen/logrus#pkg-functions
			DisableColors:             false,
			ForceQuote:                true,
			DisableQuote:              false,
			EnvironmentOverrideColors: true,
			DisableTimestamp:          false,
			FullTimestamp:             true,
			TimestampFormat:           Config.dateFormat,
			DisableSorting:            true,
			DisableLevelTruncation:    true,
			PadLevelText:              true,
		})
	} else {
		Logger.SetFormatter(&log.JSONFormatter{
			TimestampFormat:  Config.dateFormat,
			DisableTimestamp: false,
		})
	}

	return Logger

}

/*
* set logger level to Config.loggerLevel
 */
func SetLoglevel(loggerLevel string, Logger *log.Logger) (logger *log.Logger) {

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
* according to logType, ConfigLogfile set logFile to Config.accessLoggerFile or Config.errorLoggerFile
* and set file descriptor to accessFp or errorFp if logFile can be opened.
* to close the openned file on time, a defer function should be called following call this function if this return successful.
 */
func ConfigLogfile(logType string, logFile string) (fp *os.File, err error) {

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
		Config.accessFp = fp
		Config.accessLoggerFile = logFile
	} else {
		Config.errorFp = fp
		Config.errorLoggerFile = logFile

	}

	return fp, nil
}


/**
* Logging a message to Logger 
* if the Config.allstdout ,then logging the log messages to stdout
 */
func LoggingLog(logType string,logLevel string, args ...interface{}) {
	
	var logger *log.Logger 
	var tostdout bool
	var stdLogger *log.Logger

	tostdout = Config.allstdout
  logger = nil
	stdLogger = nil 

	switch strings.ToLower(logType){
		case "access":
			logger = Config.accessLogger
			break
		case "error":
			logger = Config.errorLogger
			break
		case "stdout":
			tostdout = false
			logger = Config.stdoutLogger
		default:
			tostdout = false
			logger = Config.stdoutLogger
	}

	if tostdout{
		stdLogger = Config.stdoutLogger
	}

	found := -1
	for i :=0; i<len(ConfigLogLevel);i++ {
		if strings.ToLow(logLevel) ==  ConfigLogLevel[i] {
			found = i
			break
		}
	}

	if found < 0 {
		logLevel = "debug"
	} 

	if logger != nil {
		logger = SetLoglevel(logLevel,logger)
	}

	if stdLogger != nil {
		stdLogger = SetLoglevel(stdLogger,logger)
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
				stdoutLogger.Warn(args...)
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
				stdoutLogger.Trace(args...)
			}
			break
	}
}

/*
 * Logging a message to Logger 
 * if the Config.allstdout ,then logging the log messages to stdout
 */
func LoggingLogf(logType string,logLevel string,	format string, args ...interface{}) {

	var logger *log.Logger
	var tostdout bool
	var stdLogger *log.Logger

	tostdout = Config.allstdout
	logger = nil
	stdLogger = nil 
	
	switch strings.ToLower(logType){                                                                                                                                                                                                     
		case "access":
			logger = Config.accessLogger
			break
	  case "error":
			logger = Config.errorLogger
			break
	  case "stdout":
			tostdout = false
			logger = Config.stdoutLogger
		default:
			tostdout = false
			logger = Config.stdoutLogger
  }


 if tostdout{
		stdLogger = Config.stdoutLogger
 }

	found := -1
	for i :=0; i<len(ConfigLogLevel);i++ {
		if strings.ToLow(logLevel) ==  ConfigLogLevel[i] {
			found = i
			break
	  }
	}
	 
	if found < 0 {
		logLevel = "debug"
	} 
	 
	if logger != nil {
			logger = SetLoglevel(logLevel,logger)
			switch strings.ToLower(logLevel) {
				case "panic":
					logger.Panicf(format, args...)
					break
				case "fatal":
					logger.Fatalf(format, args...)
					break;
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
			stdLogger = SetLoglevel(stdLogger,logger)
			switch strings.ToLower(logLevel) {
	      case "panic":
				  stdLogger.Panicf(format, args...)
					break     
				case "fatal":
					stdLogger.Fatalf(format, args...)
					break;    
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
   * if the Config.allstdout ,then logging the log messages to stdout
   */
  func LoggingLogln(logType string,logLevel string,args ...interface{}) {
  
    var logger *log.Logger
    var tostdout bool
    var stdLogger *log.Logger
  
    tostdout = Config.allstdout
    logger = nil
    stdLogger = nil 
    
    switch strings.ToLower(logType){                                                                                                                                                                                                     
      case "access":
        logger = Config.accessLogger
        break
      case "error":
        logger = Config.errorLogger
        break
      case "stdout":
        tostdout = false
        logger = Config.stdoutLogger
      default:
        tostdout = false
        logger = Config.stdoutLogger
    }
  
  
   if tostdout{
    stdLogger = Config.stdoutLogger                                                                                                                                                                                                    
   }
  
    found := -1
    for i :=0; i<len(ConfigLogLevel);i++ {
      if strings.ToLow(logLevel) ==  ConfigLogLevel[i] {
        found = i
        break
       }
    }
 
    if found < 0 {
      logLevel = "debug"
    }
   if logger != nil {
        logger = SetLoglevel(logLevel,logger)
        switch strings.ToLower(logLevel) {
          case "panic":
            logger.Panicln(format, args...)
            break
          case "fatal":
            logger.Fatalln(format, args...)
            break;
          case "error":
            logger.Errorln(format, args...)
            break
          case "warn":
            logger.Warnln(format, args...)
            break
          case "info":
            logger.Infoln(format, args...)
            break
          case "debug":
            logger.Debugln(format, args...)
            break
          case "trace":
            logger.Traceln(format, args...)
            break
        }
    }
  
    if stdLogger != nil {
        stdLogger = SetLoglevel(stdLogger,logger)
        switch strings.ToLower(logLevel) {
   		  case "panic":
          	stdLogger.Panicln(format, args...)
            break     
          case "fatal":
            stdLogger.Fatalln(format, args...)
            break;    
          case "error":
            stdLogger.Errorln(format, args...)
            break     
          case "warn":
            stdLogger.Warnln(format, args...)
            break    
          case "info":
            stdLogger.Infoln(format, args...)
            break    
          case "debug":
            stdLogger.Debugln(format, args...)
            break 
          case "trace":
            stdLogger.Traceln(format, args...)
            break     
        } 
   }                                                                                                                                                                                                                                    
  } 
  /*
   * Logging a message to Logger 
   * if the Config.allstdout ,then logging the log messages to stdout
   */
  func LoggingLogFn(logType string,logLevel string,fn log.LogFunction) {
  
    var logger *log.Logger
    var tostdout bool
    var stdLogger *log.Logger
  
    tostdout = Config.allstdout
    logger = nil
    stdLogger = nil 
    
    switch strings.ToLower(logType){                                                                                                                                                                                                     
      case "access":
        logger = Config.accessLogger
        break
      case "error":
        logger = Config.errorLogger
        break
      case "stdout":
        tostdout = false
        logger = Config.stdoutLogger
      default:
        tostdout = false
        logger = Config.stdoutLogger
    }
  
  
   if tostdout{
    stdLogger = Config.stdoutLogger                                                                                                                                                                                                    
   }
  
    found := -1
    for i :=0; i<len(ConfigLogLevel);i++ {
      if strings.ToLow(logLevel) ==  ConfigLogLevel[i] {
        found = i
        break
      }
    }
 
    if found < 0 {
      logLevel = "debug"
    }
   if logger != nil {
        logger = SetLoglevel(logLevel,logger)
        switch strings.ToLower(logLevel) {
          case "panic":
            logger.PanicFn(format, args...)
            break
          case "fatal":
            logger.FatalFn(format, args...)
            break;
          case "error":
            logger.ErrorFn(format, args...)
            break
          case "warn":
            logger.WarnFn(format, args...)
            break
          case "info":
            logger.InfoFn(format, args...)
            break
          case "debug":
            logger.DebugFn(format, args...)
            break
          case "trace":
            logger.TraceFn(format, args...)
            break
        }
    }
  
    if stdLogger != nil {
        stdLogger = SetLoglevel(stdLogger,logger)
        switch strings.ToLower(logLevel) {
    		case "panic":
          		stdLogger.PanicFn(format, args...)
            	break     
          	case "fatal":
            	stdLogger.FatalFn(format, args...)
            	break    
          	case "error":
            	stdLogger.ErrorFn(format, args...)
            	break     
          	case "warn":
            	stdLogger.WarnFn(format, args...)
            	break    
          	case "info":
            	stdLogger.InfoFn(format, args...)
            	break    
          	case "debug":
            	stdLogger.DebugFn(format, args...)
            	break 
          	case "trace":
            	stdLogger.TraceFn(format, args...)
            	break     
        } 
   }                                                                                                                                                                                                                                    
  } 
