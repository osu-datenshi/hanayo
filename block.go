package main

import (
	"github.com/gin-gonic/gin"
	"strings"
	"net/http"
)

func BlockIp() func(c *gin.Context) {
	return func(c *gin.Context) {
		if strings.Contains(c.ClientIP(), ":") {
			redirect()
		}
	}
}
//memes that going to destroy your life
func redirect(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "https://raw.githubusercontent.com/osu-datenshi/assets/master/fuck-ipv6.jpg", 302)
	c.Abort()
}