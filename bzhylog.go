package bzhyserver

import (
	"fmt"
	"net/http"
	"os"
	"strings"
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
  	ErrorMessage string
 	// isTerm shows whether does gin's output descriptor refers to a terminal.
  	isTerm bool
 	// BodySize is the size of the Response Body
 	BodySize int
 	// Keys are the keys set on the request's context.
 	Keys map[string]interface{}
}

// StatusCodeColor is the ANSI color for appropriately logging http status code to a terminal.
func (p *LogFormatterParams) StatusCodeColor() string {
    code := p.StatusCode

    switch {
        case code >= http.StatusOK && code < http.StatusMultipleChoices:
            return green
        case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
            return white
        case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
            return yellow
        default:
            return red
    }
}

// MethodColor is the ANSI color for appropriately logging http method to a terminal.
func (p *LogFormatterParams) MethodColor() string {
        method := p.Method

        switch method {
        case http.MethodGet:
                return blue
        case http.MethodPost:
                return cyan
        case http.MethodPut:
                return yellow
        case http.MethodDelete:
                return red
        case http.MethodPatch:
                return green
        case http.MethodHead:
                return magenta
        case http.MethodOptions:
                return white
        default:
                return reset
        }
}

// ResetColor resets all escape attributes.
func (p *LogFormatterParams) ResetColor() string {
        return reset
}

// IsOutputColor indicates whether can colors be outputted to the log.
func (p *LogFormatterParams) IsOutputColor() bool {
        return consoleColorMode == forceColor || (consoleColorMode == autoColor && p.isTerm)
}


type bzhyLoggerConfig struct {
	//The path of access log file
	AccLogFile string
	//The path of error log file
	ErrLogFile string
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
 
       // SkipPaths is a url path array which logs are not written.
       // Optional.
        SkipPaths []string
}

var LoggerConf bzhyLoggerConfig = bzhyLoggerConfig{"", "", nil, nil, nil, nil, nil, nil,nil}

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
	return fmt.Sprintf("%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
//		param.TimeStamp.Format("2006/01/02 - 15:04:05"),
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
	StdLog := bzhylog.New()
	StdLog.Out = os.Stdout
	StdLog.SetLevel(bzhylog.TraceLevel)
	StdLog.SetFormatter(&bzhylog.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	LoggerConf.StdLog = StdLog
}

// Create a new instance of the logger for Access Log.
func CreateAccLog(AccLogFile string) (ret int) {
	accFd, err := os.OpenFile(AccLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err == nil {
		accLog := bzhylog.New()
		if LoggerConf.AccFd != nil {
			CloseAccLogFd()
		}
		LoggerConf.AccLog = accLog
		LoggerConf.AccFd = accFd
		defer CloseAccLogFd()
	} else {
		LogError2StdAndFile(fmt.Sprintf("Failed to open the ACCESS log file %s Error message: %s", AccLogFile, err), "fatal")
		return 200001
	}

	return 0
}

//Closing the file descriptors of accesss log.
func CloseAccLogFd() (ret int) {
	if LoggerConf.AccFd != nil {
		err := LoggerConf.AccFd.Close()
		if err != nil {
			LogError2StdAndFile(fmt.Sprintf("Closing Access log err %s", err), "error")
		}

		LoggerConf.AccLog = LoggerConf.StdLog
		LoggerConf.AccFd = nil
	}

	return 0
}

// Create a new instance of the logger for Error Log.
func CreateErrLog(ErrLogFile string) (ret int) {
	ErrFd, err := os.OpenFile(ErrLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err == nil {
		ErrLog := bzhylog.New()
		if LoggerConf.ErrFd != nil {
			CloseErrLogFd()
		}
		LoggerConf.ErrLog = ErrLog
		LoggerConf.ErrFd = ErrFd
		defer CloseErrLogFd()
	} else {
		LogError2StdAndFile(fmt.Sprintf("Failed to open the ERROR log file %s Error message: %s", ErrLogFile, err), "fatal")
		return 200002
	}

	return 0
}

//Closing the file descriptors of error log.
func CloseErrLogFd() (ret int) {
	if LoggerConf.ErrFd != nil {
		err := LoggerConf.ErrFd.Close()
		if err != nil {
			LogError2StdAndFile(fmt.Sprintf("Closing Error log err %s", err), "error")
		}

		LoggerConf.ErrLog = LoggerConf.StdLog
		LoggerConf.ErrFd = nil
	}

	return 0
}

//Write log msg to StdOut
func WriteLog2Stdout(msg string, level string) (ret int) {
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
func WriteLog2Acclog(msg string, level string) (ret int) {
	if LoggerConf.AccLog == nil {
		OpenAccessLogger("")
	}
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
		LoggerConf.AccLog.Debug(msg)
	case "trace":
		LoggerConf.AccLog.Trace(msg)
	default:
		WriteLog2Stdout("We got a log message without UNKNOW log level", "warn")
	}

	return 0

}

func LogError2StdAndFile(msg string, level string) (ret int) {
	WriteLog2Stdout(msg, level)
	if LoggerConf.ErrLog != nil && LoggerConf.ErrLog != LoggerConf.StdLog {
		WriteLog2Errlog(msg, level)
	}

	return 0
}

func LogAccess2StdAndFile(msg string, level string) (ret int) {
	WriteLog2Stdout(msg, level)
	if LoggerConf.AccLog != nil && LoggerConf.AccLog != LoggerConf.StdLog {
		WriteLog2Acclog(msg, level)
	}

	return 0
}

//Write log msg to Error log file
func WriteLog2Errlog(msg string, level string) (ret int) {
	if LoggerConf.ErrLog == nil {
		OpenErrorLogger("")
	}
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

	return 0

}

// Logger instances a Logger middleware that will write the logs to gin.DefaultWriter.
// By default gin.DefaultWriter = os.Stdout.

func Logger() HandlerFunc {
	//Initating StdOut Logger
	if LoggerConf.StdLog == nil {
		CreateStdLog()
	}

        if LoggerConf.AccLog == nil || LoggerConf.AccLog == LoggerConf.StdLog {
             OpenAccessLogger("")
        }


        if LoggerConf.ErrLog == nil || LoggerConf.ErrLog == LoggerConf.StdLog {
              OpenErrorLogger("")
        } 
         

	if LoggerConf.Formatter == nil {
		LoggerConf.Formatter = defaultLogFormatter
	}
       
        isTerm := true
        
        var skip map[string]struct{}

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

			if len(param.ErrorMessage) > 0  {
				LogError2StdAndFile(LoggerConf.Formatter(param), "error")
			} else {
				LogAccess2StdAndFile(LoggerConf.Formatter(param), "info")
			}
		}
	}

}


// Set Access logfile and create Access logger
func OpenAccessLogger(AccessLogFile string) (ret int) {

	if len(AccessLogFile) > 0 {
		ret := CreateAccLog(AccessLogFile)
		if ret == 0 {
			LoggerConf.AccLogFile = AccessLogFile
		}
	}

	if len(LoggerConf.AccLogFile) >0 && (LoggerConf.AccLog == nil || LoggerConf.AccLog == LoggerConf.StdLog) {
		CreateAccLog(LoggerConf.AccLogFile)
	}

	if LoggerConf.AccLog == nil {
		LoggerConf.AccLog = LoggerConf.StdLog
	}

	return 0

}

//Set Error logfile and create Error logger
func OpenErrorLogger(ErrorLogFile string) (ret int) {
	if len(ErrorLogFile) > 0 {
		ret := CreateErrLog(ErrorLogFile)
		if ret == 0 {
			LoggerConf.ErrLogFile = ErrorLogFile
		}
	}

	if len(LoggerConf.ErrLogFile) > 0 && (LoggerConf.ErrLog == nil || LoggerConf.ErrLog == LoggerConf.StdLog) {
		CreateErrLog(LoggerConf.ErrLogFile)
	}

	if LoggerConf.ErrLog == nil {
		LoggerConf.ErrLog = LoggerConf.StdLog
	}

	return 0
}

//Set Access file to LoggerConf
func SetAccessFile(AccessFile string) (ret int) {
	if len(AccessFile) < 1 {
		LogError2StdAndFile(fmt.Sprintf("Access file path is NULL"),"warn")
		return 200003
	}
	if strings.Compare(LoggerConf.AccLogFile,AccessFile) != 0 {
		AccFd, err := os.OpenFile(AccessFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			LogError2StdAndFile(fmt.Sprintf("Failed to open the ACCESS log file %s Error message: %s", AccessFile, err), "fatal")
			return 200004
		}
		AccFd.Close()
		OpenAccessLogger(AccessFile)
		return 0
	} 
	
	return 0
}

//Set Error file to LoggerConf
func SetErrorFile(ErrorFile string) (ret int) {
	if len(ErrorFile) < 1 {
		LogError2StdAndFile(fmt.Sprintf("Error file path is NULL"),"warn")
		return 200005
	}
	if strings.Compare(LoggerConf.ErrLogFile,ErrorFile) != 0 {
		ErrFd, err := os.OpenFile(ErrorFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			LogError2StdAndFile(fmt.Sprintf("Failed to open the ERROR log file %s Error message: %s", ErrorFile, err), "fatal")
			return 200004
		}
		ErrFd.Close()
		OpenErrorLogger(ErrorFile)
		return 0
	} 
	
	return 0
}