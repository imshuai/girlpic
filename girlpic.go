package main

import (
	"net/http"
	"os"
	"time"

	logger "github.com/orzzz/lightlog"

	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

var debug = logger.Debug
var info = logger.Info
var warn = logger.Warn
var errorlog = logger.Error
var fatal = logger.Fatal

var db *xorm.Engine
var err error
var delimiter string
var dir string

const numPerPage = 12

type dbConfig struct {
	DBip      string `json:"db_ip"`
	DBport    string `json:"db_port"`
	DBname    string `json:"db_name"`
	DBuser    string `json:"db_username"`
	DBpasswd  string `json:"db_password"`
	DBcharset string `json:"db_charset"`
}

func getDbConfig() dbConfig {
	var conf dbConfig
	if os.Getenv("APPENV") == "Daocloud" {
		conf.DBip = os.Getenv("MYSQL_PORT_3306_TCP_ADDR")
		conf.DBport = os.Getenv("MYSQL_PORT_3306_TCP_PORT")
		conf.DBname = os.Getenv("MYSQL_INSTANCE_NAME")
		conf.DBuser = os.Getenv("MYSQL_USERNAME")
		conf.DBpasswd = os.Getenv("MYSQL_PASSWORD")
		conf.DBcharset = "utf8"
	} else if os.Getenv("APPENV") == "Product" {
		conf.DBip = os.Getenv("MYSQL_PORT_3306_TCP_ADDR")
		conf.DBport = os.Getenv("MYSQL_PORT_3306_TCP_PORT")
		conf.DBname = "prissh"
		conf.DBuser = "prissh"
		conf.DBpasswd = "prissh"
		conf.DBcharset = "utf8"
	} else {
		conf.DBip = "192.168.1.2"
		conf.DBport = "3306"
		conf.DBname = "spider"
		conf.DBuser = "spider"
		conf.DBpasswd = "spider"
		conf.DBcharset = "utf8"
	}
	return conf
}

func init() {
	if os.IsPathSeparator('\\') { //前边的判断是否是系统的分隔符
		delimiter = "\\"
	} else {
		delimiter = "/"
	}
	dir, _ = os.Getwd() //当前的目录

	gin.SetMode(gin.ReleaseMode)

	logger.SetLevel(logger.INFO)
	logger.SetPrefix("[GirlPic]")
	logger.SetRollingFile(dir+delimiter+"logs", "server.log", 10, 1, logger.MB)

	conf := getDbConfig()
	db, err = xorm.NewEngine("mysql",
		conf.DBuser+":"+
			conf.DBpasswd+"@tcp("+
			conf.DBip+":"+
			conf.DBport+")/"+
			conf.DBname+"?charset="+conf.DBcharset+"&timeout=5s&parseTime=True&loc=Asia%2FChongqing")
	if err != nil {
		fatal("Got error when connect database, the error:", err)
		os.Exit(1)
	}
	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(100)
	db.SetTableMapper(core.SnakeMapper{})
	db.SetColumnMapper(core.GonicMapper{})

	db.TZLocation, _ = time.LoadLocation("Asia/Chongqing")
	//      db.ShowExecTime(true)
	//db.ShowSQL(true)
	err = db.Sync2(new(GirlPic))
	if err != nil {
		fatal("Database struct sync fail, error:", err)
		os.Exit(1)
	}
	info("Server initialization complete! Current directory is [", dir, "] , PathSeparator is [", delimiter, "]")
}

//GirlPic define database struct for table girl_pic
type GirlPic struct {
	ID         int64
	URL        string    `xorm:"unique"`
	Review     bool      `xorm:"notnull default 0"`
	Like       int       `xorm:"notnull default 0"`
	Unlike     int       `xorm:"notnull default 0"`
	CreateTime time.Time `xorm:"created"`
}

type tPic struct {
	URL    string `json:"url"`
	ID     int64  `json:"id"`
	Like   int    `json:"like"`
	Unlike int    `json:"unlike"`
}

func getPicNum() int {
	p := new(GirlPic)
	length, err := db.Where("`review` = ?", 1).Count(p)
	if err != nil {
		errorlog("Database query fail, error:", err)
		return 0
	}
	return int(length)
}

func getPics(page int) []tPic {
	num := getPicNum()
	if num == 0 {
		info("number of pics is:", num)
		return nil
	}
	var start int
	if start = (page - 1) * numPerPage; start < 0 || start > num {
		start = 0
	}
	pic := new(GirlPic)
	rows, err := db.Where("`review` = ?", 1).Desc("id").Limit(numPerPage, start).Rows(pic)
	if err != nil {
		errorlog("database query fail, error:", err)
		return nil
	}
	defer rows.Close()
	pics := make([]tPic, 0)
	for rows.Next() {
		rows.Scan(pic)
		pics = append(pics, tPic{
			URL:    pic.URL,
			ID:     pic.ID,
			Like:   pic.Like,
			Unlike: pic.Unlike,
		})
	}
	return pics
}

func getPicWaitReview() GirlPic {
	pic := new(GirlPic)
	exist, err := db.Where("`review` = ?", 0).Get(pic)
	if err != nil {
		errorlog("database query fail, error:", err)
		return GirlPic{}
	}
	if !exist {
		info("all pics been reviewd")
	}
	return *pic
}

func main() {
	e := gin.Default()

	e.Static("/static", "./static")
	e.StaticFile("/favicon.ico", "./static/favicon.ico")
	e.LoadHTMLGlob("./tmpls/*")

	e.GET("/", func(c *gin.Context) {
		if c.Request.Header.Get("X-Requested-With") == "XMLHttpRequest" {
			pics := getPics(0)
			c.JSON(http.StatusOK, pics)
			return
		}
		c.HTML(http.StatusOK, "index.html", nil)
	})
	e.GET("/detail/:i", func(c *gin.Context) {
		picID, err := strconv.Atoi(c.Param("i"))
		if err != nil {
			info("pass wrong pic id to route")
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		pic := new(GirlPic)
		exist, err := db.ID(picID).Get(pic)
		if err != nil {
			errorlog("database query fail, error:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if !exist {
			info("pass wrong pic id to database")
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		c.HTML(http.StatusOK, "detail.html", pic)
	})
	e.GET("/detail/:i/:action", func(c *gin.Context) {
		if c.Request.Header.Get("X-Requested-With") == "XMLHttpRequest" {
			picID, err := strconv.Atoi(c.Param("i"))
			if err != nil {
				info("pass wrong pic id to route")
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
			action := c.Param("action")
			switch action {
			case "like":
				db.Exec("update `girl_pic` set `like` = `like` + 1 where `id`=?", picID)
				c.String(http.StatusOK, "%v", "ok")
				break
			case "unlike":
				db.Exec("update `girl_pic` set `unlike` = `unlike` + 1 where `id`=?", picID)
				c.String(http.StatusOK, "%v", "ok")
				break
			default:
				c.AbortWithStatus(http.StatusBadRequest)
			}
		} else {
			c.AbortWithStatus(http.StatusMethodNotAllowed)
		}
	})
	e.GET("/page/:i", func(c *gin.Context) {
		page, err := strconv.Atoi(c.Param("i"))
		if err != nil {
			info("pass wrong page param to route")
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		pics := getPics(page)
		c.JSON(http.StatusOK, pics)
	})
	manage := e.Group("/review", gin.BasicAuth(gin.Accounts{
		"admin": "shuai6563",
	}))
	manage.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "review.html", nil)
	})
	manage.GET("/next", func(c *gin.Context) {
		if c.Request.Header.Get("X-Requested-With") == "XMLHttpRequest" {
			pic := getPicWaitReview()
			c.JSON(http.StatusOK, gin.H{
				"id":  pic.ID,
				"url": pic.URL,
			})
		} else {
			c.AbortWithStatus(http.StatusMethodNotAllowed)
		}
	})
	manage.GET("/save/:id/:action", func(c *gin.Context) {
		if c.Request.Header.Get("X-Requested-With") == "XMLHttpRequest" {
			pic := new(GirlPic)
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				errorlog("pass wrong id param to route")
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
			switch c.Param("action") {
			case "accept":
				pic.Review = true
				effect, err := db.ID(id).Cols("review").Update(pic)
				if err != nil {
					errorlog("database update fail, error:", err)
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}
				c.JSON(http.StatusOK, gin.H{
					"effect": effect,
				})
				break
			case "reject":
				effect, err := db.ID(id).Delete(pic)
				if err != nil {
					errorlog("database delete fail, error:", err)
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}
				c.JSON(http.StatusOK, gin.H{
					"effect": effect,
				})
				break
			default:
				c.AbortWithStatus(http.StatusBadRequest)
			}
		} else {
			c.AbortWithStatus(http.StatusMethodNotAllowed)
		}
	})
	if os.Getenv("APPENV") != "" {
		e.Run(":80")
	} else {
		e.Run(":8080")
	}
}
