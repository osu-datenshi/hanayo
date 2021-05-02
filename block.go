package main

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func BlockIp() func(c *gin.Context) {
	return func(c *gin.Context) {
		if strings.Contains(c.ClientIP(), ":") {
			c.Redirect(302, "/blockedipv6")
		}
	}
}