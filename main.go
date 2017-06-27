package main

import (
	"net/http"

	"strconv"

	"fmt"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/imshuai/serverchan"
)

const numPerPage = 25

var wnotice = serverchan.NewServerChan("SCU9423Tf5f3aa306a2fa6a1a75394d80e1113c7594b8169596ed")

func main() {
	e := gin.Default()

	e.Static("/static", "./static")
	e.StaticFile("/favicon.ico", "./static/favicon.ico")
	e.LoadHTMLGlob("./tmpls/*")

	e.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	e.POST("/data/pics", func(c *gin.Context) {
		pics := make([]string, 0)
		err := c.Bind(&pics)
		if err != nil {
			errorlog("bind post pics to []string fail with error:", err)
			c.String(200, "fail with error:%v", err)
			return
		}
		conn := rdb.Get()
		defer conn.Close()
		for _, v := range pics {
			conn.Send("LPUSH", listUnchecked, v)
		}
		r, err := conn.Do("")
		if err != nil {
			errorlog("store post pics fail with error:", err)
			c.String(200, "%v", "fail")
			return
		}
		debug("store post pics success, return data:", r)
		wnotice.Send("GirlPic有新图片需要审核了", fmt.Sprintf("\n\n接收到%d张新图片，尽快审核放出", len(pics)))
		c.String(200, "%v", "ok")
	})

	e.GET("/detail/:i", func(c *gin.Context) {
		picID, err := strconv.Atoi(c.Param("i"))
		if err != nil {
			info("pass wrong pic id to route")
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		conn := rdb.Get()
		defer conn.Close()
		str, err := redis.String(conn.Do("LINDEX", listChecked, picID))
		if err != nil {
			errorlog("get", picID, "fail with error:", err)
			info("pass wrong pic id to route")
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		pic := tPic{
			ID:      picID,
			GirlPic: deserialize(str),
		}
		c.HTML(http.StatusOK, "detail.html", pic)
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
	e.GET("/comments", func(c *gin.Context) {
		c.HTML(http.StatusOK, "comment.html", nil)
	})
	e.POST("/comments/new", func(c *gin.Context) {
		if c.Request.Header.Get("X-Requested-With") == "XMLHttpRequest" {
			var com Comments
			if c.Bind(&com) == nil {

			} else {
				c.AbortWithStatus(http.StatusBadRequest)
			}
		} else {
			c.AbortWithStatus(http.StatusMethodNotAllowed)
		}
	})
	manage := e.Group("/review", gin.BasicAuth(gin.Accounts{
		"admin": "shuai6563",
	}))
	manage.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "review.html", nil)
	})
	manage.GET("/next", func(c *gin.Context) {
		if c.Request.Header.Get("X-Requested-With") == "XMLHttpRequest" {
			pics := getPicWaitReview()
			c.JSON(http.StatusOK, pics)
		} else {
			c.AbortWithStatus(http.StatusMethodNotAllowed)
		}
	})
	manage.GET("/delete/:id", func(c *gin.Context) {
		if c.Request.Header.Get("X-Requested-With") == "XMLHttpRequest" {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				errorlog("pass wrong id param to route")
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
			conn := rdb.Get()
			defer conn.Close()
			conn.Do("LSET", listTemp, id, "deleted")
			c.JSON(http.StatusOK, gin.H{
				"effect": 1,
			})
		} else {
			c.AbortWithStatus(http.StatusMethodNotAllowed)
		}
	})
	e.Run(":34533")
}
