package main

import (
	"github.com/gin-gonic/gin"
	"strings"
	"net/http"
)

func BlockIp() func(c *gin.Context) {
	return func(c *gin.Context) {
	if strings.Contains(c.ClientIP(), ":") {
		c.String(http.StatusForbidden, "IPV6 not allowed!")
		c.Abort()
	}
}