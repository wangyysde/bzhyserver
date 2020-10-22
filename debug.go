// Copyright 2020 Wayne Wang<net_use@bzhy.com>.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bzhyserver

import (
	"fmt"
	"html/template"
	"runtime"
	"strconv"
	"strings"
)

const ginSupportMinGoVer = 10

// IsDebugging returns true if the framework is running in debug mode.
// Use SetMode(gin.ReleaseMode) to disable debug mode.
func IsDebugging() bool {
	return ginMode == debugCode
}

// DebugPrintRouteFunc indicates debug log output format.
var DebugPrintRouteFunc func(httpMethod, absolutePath, handlerName string, nuHandlers int)

func debugPrintRoute(httpMethod, absolutePath string, handlers HandlersChain) {
	if IsDebugging() {
		nuHandlers := len(handlers)
		handlerName := nameOfFunction(handlers.Last())
		if DebugPrintRouteFunc == nil {
			LogError2StdAndFile(fmt.Sprintf("%-6s %-25s --> %s (%d handlers)\n", httpMethod, absolutePath, handlerName, nuHandlers),"debug")
		} else {
			LogError2StdAndFile(fmt.Sprintf("%-6s %-25s %s %d handlers",httpMethod, absolutePath, handlerName, nuHandlers),"debug")
		}
	}
}

func debugPrintLoadTemplate(tmpl *template.Template) {
	if IsDebugging() {
		var buf strings.Builder
		for _, tmpl := range tmpl.Templates() {
			buf.WriteString("\t- ")
			buf.WriteString(tmpl.Name())
			buf.WriteString("\n")
		}
		LogAccess2StdAndFile(fmt.Sprintf("Loaded HTML Templates (%d): \n%s\n", len(tmpl.Templates()), buf.String()),"info")
	}
}

func debugPrint(format string, values ...interface{}) {
	if IsDebugging() {
		if !strings.HasSuffix(format, "\n") {
			format += "\n"
		}
		LogError2StdAndFile(fmt.Sprintf(format, values...),"debug")
	}
}

func getMinVer(v string) (uint64, error) {
	first := strings.IndexByte(v, '.')
	last := strings.LastIndexByte(v, '.')
	if first == last {
		return strconv.ParseUint(v[first+1:], 10, 64)
	}
	return strconv.ParseUint(v[first+1:last], 10, 64)
}

func debugPrintWARNINGDefault() {
	if v, e := getMinVer(runtime.Version()); e == nil && v <= ginSupportMinGoVer {
		LogError2StdAndFile("Now bzhyserver requires Go 1.11 or later and Go 1.12 will be required soon.","warn")
	}
	LogError2StdAndFile("Creating an Engine instance with the Logger and Recovery middleware already attached.","warn")
}

func debugPrintWARNINGNew() {
	LogError2StdAndFile(`Running in "debug" mode. Switch to "release" mode in production.`,"debug")
	LogError2StdAndFile("- using env:	export SERVER_MODE=release","debug")
	LogError2StdAndFile("- using code:	bzhyserver.SetMode(bzhyserver.ReleaseMode)","debug")
}

func debugPrintWARNINGSetHTMLTemplate() {
	LogError2StdAndFile(`Since SetHTMLTemplate() is NOT thread-safe. It should only be called
at initialization. ie. before any route is registered or the router is listening in a socket:

	router := bzhyserver.Default()
	router.SetHTMLTemplate(template) // << good place

`,"debug")
}

func debugPrintError(err error) {
	if err != nil {
		if IsDebugging() {
			LogError2StdAndFile(fmt.Sprintf("%v\n",err),"error")
		}
	}
}
