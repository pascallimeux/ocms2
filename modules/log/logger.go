/*
Copyright Pascal Limeux. 2016 All Rights Reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
		 http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package log

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

var (
	trace   *log.Logger
	info    *log.Logger
	warning *log.Logger
	error   *log.Logger
	fatal   *log.Logger
	mode    int
	modes   = map[string]int{"Trace": 0, "Info": 1, "Warning": 2, "Error": 3, "Fatal": 4}
)

// Initialize logger, with mode, filename etc.
func Init_log(filename, modestr string) *os.File {
	var ok bool
	if mode, ok = modes[modestr]; !ok {
		mode = modes["Warning"]
	}
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("error opening file: %v", err)
	} else {
		log.SetOutput(file)
	}
	trace = log.New(file, "TRACE: ", log.Ldate|log.Ltime)
	info = log.New(file, "INFO: ", log.Ldate|log.Ltime)
	warning = log.New(file, "WARNING: ", log.Ldate|log.Ltime)
	error = log.New(file, "ERROR: ", log.Ldate|log.Ltime)
	fatal = log.New(file, "FATAL: ", log.Ldate|log.Ltime)
	return file
}

func Trace(messages ...string) {
	if modes["Trace"] >= mode {
		trace.Println(build_message(messages))
	}
}

func Info(messages ...string) {
	if modes["Info"] >= mode {
		info.Println(build_message(messages))
	}
}

func Warning(messages ...string) {
	if modes["Warning"] >= mode {
		warning.Println(build_message(messages))
	}
}

func Error(messages ...string) {
	if modes["Error"] >= mode {
		error.Println(build_message(messages))
	}
}

func Fatal(messages ...string) {
	fatal.Println(build_message(messages))
	os.Exit(-1)

}

func build_message(messages []string) string {
	var message string
	for _, m := range messages {
		message += m
	}
	return message
}

// return a string containing the file name, function name
// and the line number of a specified entry on the call stack
func Here(depthList ...int) string {
	var depth int
	if depthList == nil {
		depth = 1
	} else {
		depth = depthList[0]
	}
	_, file, line, _ := runtime.Caller(depth)
	return fmt.Sprintf(" %s :%d: ", chopPath(file), line)
}

// return the source filename after the last slash
func chopPath(original string) string {
	i := strings.LastIndex(original, "/")
	if i == -1 {
		return original
	} else {
		return original[i+1:]
	}
}
