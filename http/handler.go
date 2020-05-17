package http

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// MyContext ...
type MyContext struct {
	*gin.Context
}

// HandlerFunc ...
type HandlerFunc func(*MyContext)

func handler(handler HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		context := new(MyContext)
		context.Context = c
		handler(context)
	}
}

// RespCommon ...
type RespCommon struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// HTTPErrorCode define ...
const (
	HTTPErrorCodeSuccess = 0
	HTTPErrorCodeFail    = -1
)

func pingHandler(c *MyContext) {
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST")
		c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func httpErrorCommon(c *MyContext, httpCode *int, err *error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r, string(debug.Stack()))
		}
	}()
	if *err != nil {
		c.JSON(*httpCode, &RespCommon{
			Code: HTTPErrorCodeFail,
			Msg:  (*err).Error(),
		})
		log.Println(c.Request.URL.Path+" error:", (*err).Error())
	}
}
