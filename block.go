package main

import (
	"github.com/gin-gonic/gin"
	"strings"
	"net/http"
)

func BlockIp() func(c *gin.Context) {
	return func(c *gin.Context) {
		if strings.Contains(c.ClientIP(), ":") {
			c.Redirect(http.StatusMovedPermanently, "https://raw.githubusercontent.com/osu-datenshi/assets/master/fuck-ipv6.jpg")
		}
	}
}