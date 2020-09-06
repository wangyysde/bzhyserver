package bzhyserver
import (
        "strings"
        "github.com/wangyysde/bzhylog"
	"os"
)

//Define loggger for
var (
        StdLog *bzhylog.Logger = nil
)


func CreateStdLog() {
        StdLog =  bzhylog.New()
        StdLog.Out = os.Stdout
        StdLog.SetLevel(bzhylog.TraceLevel)
        StdLog.SetFormatter(&bzhylog.TextFormatter{
                DisableColors: false,
                FullTimestamp: true,
        })
}

//Write log msg to StdOut
func WriteLog2Stdout(msg string, level string){
	if StdLog == nil {
		CreateStdLog()
	}
        switch strings.ToLower(level) {
                case "panic":
                        StdLog.Panic(msg)
                case "fatal":
                        StdLog.Fatal(msg)
                case "error":
                        StdLog.Error(msg)
                case "warn", "warning":
                        StdLog.Warn(msg)
                case "info":
                        StdLog.Info(msg)
                case "debug":
                        StdLog.Debug(msg)
                case "trace":
                        StdLog.Trace(msg)
                default:
                        WriteLog2Stdout("We got a log message without UNKNOW log level", "warn")
        }


}
