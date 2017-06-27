package main

import (
	"os"

	"github.com/gin-gonic/gin"
	logger "github.com/orzzz/lightlog"
)

var err error
var delimiter string
var dir string

var debug = logger.Debug
var info = logger.Info
var warn = logger.Warn
var errorlog = logger.Error
var fatal = logger.Fatal

const (
	listChecked   = "list:checked"
	listUnchecked = "list:unchecked"
	listTemp      = "list:temp"
)

func init() {
	if os.IsPathSeparator('\\') { //前边的判断是否是系统的分隔符
		delimiter = "\\"
	} else {
		delimiter = "/"
	}
	dir, _ = os.Getwd() //当前的目录

	gin.SetMode(gin.ReleaseMode)

	logger.SetLevel(logger.ERROR)
	logger.SetPrefix("[GirlPic]")
	logger.SetRollingFile(dir+delimiter+"logs", "server.log", 5, 1, logger.MB)

	redisInit()

	info("Server initialization complete! Current directory is [", dir, "] , PathSeparator is [", delimiter, "]")
}
