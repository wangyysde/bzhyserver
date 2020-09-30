// Copyright 2020 Wayne wang<net_use@bzhy.com>.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bzhyserver

import (
        "fmt"
        "io"
        "net/http"
        "os"
        "time"

        "github.com/wangyysde/bzhylog"
)

// LogFormatter gives the signature of the formatter function passed to LoggerWithFormatter
type LogFormatter func(params LogFormatterParams) string 

// LogFormatterParams is the structure any formatter will be handed when time to log comes
type LogFormatterParams struct {
   Request *http.Request

  // TimeStamp shows the time after the server returns a response.
  TimeStamp time.Time
  // StatusCode is HTTP response code.
  StatusCode int
  // Latency is how much time the server cost to process a certain request.
  Latency time.Duration
  // ClientIP equals Context's ClientIP method.
  ClientIP string
  // Method is the HTTP method given to the request.
  Method string
  // Path is a path the client requests.
  Path string
  // ErrorMessage is set if error has occurred in processing the request.
  ErrorMessage string = nil 
 // isTerm shows whether does gin's output descriptor refers to a terminal.
  isTerm bool
 // BodySize is the size of the Response Body
 BodySize int
 // Keys are the keys set on the request's context.
 Keys map[string]interface{}
}

type LoggerConfig struct {
   //The path of access log file 
   AccLogFile  string 
   //The path of error log file
   ErrLogFile  string 
   //Logger for Stdout
   StdLog *bzhylog.Logger 
   //Logger for Access Log
   AccLog *bzhylog.Logger 
   //Logger for Stdout
   ErrLog *bzhylog.Logger 
   //The access log file descriptor
   AccFd *os.File 
   //The error log file descriptor
   ErrFd *os.File 
   
   Formatter LogFormatter  

}   

var LoggerConf LoggerConfig = {nil, nil, nil , nil,nil, nil,nil,nil}

// defaultLogFormatter is the default log format function Logger middleware uses.
var defaultLogFormatter = func(param LogFormatterParams) string {
        var statusColor, methodColor, resetColor string
        if param.IsOutputColor() {
                statusColor = param.StatusCodeColor()
                methodColor = param.MethodColor()
                resetColor = param.ResetColor()
        }

        if param.Latency > time.Minute {
                // Truncate in a golang < 1.8 safe way
                param.Latency = param.Latency - param.Latency%time.Second
        }
        return fmt.Sprintf(" %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
                param.TimeStamp.Format("2006/01/02 - 15:04:05"),
                statusColor, param.StatusCode, resetColor,
                param.Latency,
                param.ClientIP,
                methodColor, param.Method, resetColor,
                param.Path,
                param.ErrorMessage,
        )
}

// Create a new instance of the logger for StdOut.
func CreateStdLog() {
     StdLog =  bzhylog.New()
     StdLog.Out = os.Stdout
     StdLog.SetLevel(bzhylog.TraceLevel)
     StdLog.SetFormatter(&bzhylog.TextFormatter{
        DisableColors: false,
        FullTimestamp: true,
    })

    LoggerConf.StdLog = StdLog
}

// Create a new instance of the logger for Access Log.
func CreateAccLog(AccLogFile string)(ret int)  {
        accLog = bzhylog.New()
        accFd, err := os.OpenFile(AccLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
        if err == nil {
                LoggerConf.AccLog = accLog
                LoggerConf.AccFd = accFd
                defer CloseAcceLogFd()
        } else {
                WriteLog2Stdout(fmt.Sprintf("Failed to open the ACCESS log file %s Error message: %s",AccLogFile,err), "fatal")
                return 200001
        }

        return 0
}

//Closing the file descriptors of accesss log.
func CloseAccLogFd()(ret int){
        if LoggerConf.AccFd != nil {
                err := LoggerConf.AccFd.Close()
                if err != nil {
                        WriteLog2Errlog(fmt.Sprintf("Closing Access log err %s", err),"error")
                }
                LoggerConf.AccFd = nil
        }

        return 0
}

// Create a new instance of the logger for Error Log.
func CreateErrLog(ErrLogFile string )(ret int)  {
        ErrLog = bzhylog.New()
        ErrFd, err := os.OpenFile(ErrLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
        if err == nil {
                LoggerConf.ErrLog = ErrLog
                LoggerConf.ErrFd = ErrFd
                defer CloseErrLogFd()
        } else {
                WriteLog2Stdout(fmt.Sprintf("Failed to open the ERROR log file %s Error message: %s",ErrLogFile,err), "fatal")
                return 200002
        }


        return 0
}

//Closing the file descriptors of error log.
func CloseErrLogFd()(ret int){
        if LoggerConf.ErrFd != nil {
                LoggerConf.ErrFd.Close()
        }
        return 0
}



//Write log msg to StdOut
func WriteLog2Stdout(msg string, level string) (ret int){
        if LoggerConf.StdLog == nil {
           CreateStdLog()
        }
        switch strings.ToLower(level) {
                case "panic":
                        LoggerConf.StdLog.Panic(msg)
                case "fatal":
                        LoggerConf.StdLog.Fatal(msg)
                case "error":
                        LoggerConf.StdLog.Error(msg)
                case "warn", "warning":
                        LoggerConf.StdLog.Warn(msg)
                case "info":
                        LoggerConf.StdLog.Info(msg)
                case "debug":
                        LoggerConf.StdLog.Debug(msg)
                case "trace":
                        LoggerConf.StdLog.Trace(msg)
                default:
                        WriteLog2Stdout("We got a log message without UNKNOW log level", "warn")
        }

        return 0

}

//Write log msg to Access log file
func WriteLog2Acclog(msg string, level string) (ret int){
        switch strings.ToLower(level) {
                case "panic":
                        LoggerConf.AccLog.Panic(msg)
                case "fatal":
                        LoggerConf.AccLog.Fatal(msg)
                case "error":
                        LoggerConf.AccLog.Error(msg)
                case "warn", "warning":
                        LoggerConf.AccLog.Warn(msg)
                case "info":
                        LoggerConf.AccLog.Info(msg)
                case "debug":
                        AccLog.Debug(msg)
                case "trace":
                        LoggerConf.AccLog.Trace(msg)
                default:
                        WriteLog2Stdout("We got a log message without UNKNOW log level", "warn")
        }

        return 0

}

//Write log msg to Error log file
func WriteLog2Errlog(msg string, level string) (ret int){
        switch strings.ToLower(level) {
          case "panic":
                LoggerConf.ErrLog.Panic(msg)
          case "fatal":
                LoggerConf.ErrLog.Fatal(msg)
          case "error":
                LoggerConf.ErrLog.Error(msg)
          case "warn", "warning":
                LoggerConf.ErrLog.Warn(msg)
          case "info":
                LoggerConf.ErrLog.Info(msg)
          case "debug":
                LoggerConf.ErrLog.Debug(msg)
          case "trace":
                LoggerConf.ErrLog.Trace(msg)
          default:
              WriteLog2Stdout("We got a log message without UNKNOW log level", "warn")
        }
      

// Logger instances a Logger middleware that will write the logs to gin.DefaultWriter.
// By default gin.DefaultWriter = os.Stdout.
func Logger() HandlerFunc {
   //Initating StdOut Logger
   if LoggerConf.StdLog == nil {
        CreateStdLog()    
   }

   if LoggerConf.AccLog == nil && LoggerConf.AccLogFile != nil {
        var ret := 0
        ret = CreateAccLog(LoggerConf.AccLogFile)
        if ret > 0 {
                LoggerConf.AccLog = LoggerConf.StdLog
        }
   }

   if LoggerConf.ErrLog == nil && LoggerConf.ErrLogFile != nil {
        var ret := 0
        ret = CreateErrLog(LoggerConf.ErrLogFile)
        if ret > 0 {
                LoggerConf.ErrLog = LoggerConf.StdLog
        }
   }

   if LoggerConf.Formatter == nil {
        LoggerConf.Formatter =  defaultLogFormatter   
   }
   return func(c *Context) {
      // Start timer
      start := time.Now()
      path := c.Request.URL.Path
      raw := c.Request.URL.RawQuery

      // Process request
      c.Next()

      // Log only when path is not being skipped
      if _, ok := skip[path]; !ok {
        param := LogFormatterParams{
          Request: c.Request,
          isTerm:  isTerm,
          Keys:    c.Keys,
        }

        // Stop timer
        param.TimeStamp = time.Now()
        param.Latency = param.TimeStamp.Sub(start)

        param.ClientIP = c.ClientIP()
        param.Method = c.Request.Method
        param.StatusCode = c.Writer.Status()
        param.ErrorMessage = c.Errors.ByType(ErrorTypePrivate).String()

        param.BodySize = c.Writer.Size()

        if raw != "" {
           path = path + "?" + raw
        }

        param.Path = path

        if param.ErrorMessage != nil {
                WriteLog2Errlog(LoggerConf.Formatter(param),"error")     
        }
        else{
                WriteLog2Acclog(LoggerConf.Formatter(param),"info")
        }
     }
  }

}

//Set Access logfile and create Access logger
func OpenAccessLogger(AccessLogFile string) int {
 
 //Checking if AccessLogFile is empty  
 if len(AccessLogFile) {
    WriteLog2Stdout("OpenAccessLogger has been called, but the AccessLogFile is empty","warn")
    return 2003
  }

  //Retruning  if  AccLog has been created  and it is not StdLogger 
  if LoggerConf.AccLog != nil && LoggerConf.AccLog != LoggerConf.StdLog {
    WriteLog2Stdout("OpenAccessLogger has been called, but access logger has been setted","warn")
    return 2004
 }
 
 //Returning if create accesslog error
 ret := CreateAccLog(AccessLogFile)
 if ret > 0 {
    return 2005
 }
 
 LoggerConf.AccLogFile = AccessLogFile
 
 return 0
}


//Set Error logfile and create Error logger
func OpenErrorLogger(ErrorLogFile string) int {

 //Checking if ErrorLogFile is empty
 if len(ErrorLogFile) {
    WriteLog2Stdout("OpenErrorLogger has been called, but the ErrorLogFile is empty","warn")
    return 2006
  }

  //Retruning  if  ErrLog has been created  and it is not StdLogger
  if LoggerConf.ErrLog != nil && LoggerConf.ErrLog != LoggerConf.StdLog {
    WriteLog2Stdout("OpenErrorLogger has been called, but error logger has been setted","warn")
    return 2007
 }

 //Returning if create accesslog error
 ret := CreateErrLog(ErrorLogFile)
 if ret > 0 {
    return 2008
 }

 LoggerConf.ErrLogFile = ErrorLogFile

 return 0
}
 
