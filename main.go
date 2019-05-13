package main

import (
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lionsoul2014/ip2region/binding/golang/ip2region"
)

var logger *Logger

func init() {
	rand.Seed(time.Now().Unix())

	SetLevel("info")
	logger = NewLogger(os.Stdout)
	gin.SetMode(gin.ReleaseMode)
}

func mapRoutes() *gin.Engine {
	ret := gin.New()
	ret.Use(gin.Recovery())

	ret.GET("", ipToRegion)
	ret.NoRoute(func(c *gin.Context) {
		c.String(http.StatusOK, "The piper will lead us to reason.\n\n欢迎访问黑客与画家的社区 https://hacpai.com")
	})

	return ret
}

func ipToRegion(c *gin.Context) {
	result := NewResult()
	result.Code = CodeErr

	ip := c.Query("ip")
	ipr, err := region.MemorySearch(ip)
	if nil != err {
		msg := err.Error()
		logger.Errorf("search ip [" + ip + "] failed: " + msg)
		result.Msg = msg
		c.JSON(http.StatusOK, result)

		return
	}

	result.Data = map[string]interface{}{
		"country":  ipr.Country,
		"province": ipr.Province,
		"city":     ipr.City,
	}
	result.Code = CodeOk
	c.JSON(http.StatusOK, result)
}

var region *ip2region.Ip2Region

func main() {
	region, _ = ip2region.New("ip2region.db")
	defer region.Close()

	router := mapRoutes()
	server := &http.Server{
		Addr:    "127.0.0.1:1126",
		Handler: router,
	}

	logger.Infof("ip2region-http is running")
	server.ListenAndServe()
}
