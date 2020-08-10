package main

import (
	"github.com/gin-gonic/gin"
	"zxq.co/x/rs"
)

func DiscordGenToken(c *gin.Context) {
	ctx := getContext(c)
	if ctx.User.ID == 0 {
		resp403(c)
		return
	}

	db.Exec("DELETE FROM discord_tokens WHERE userid = ?", ctx.User.ID)

	key := rs.String(32)

	db.Exec("INSERT INTO discord_tokens(userid, token) VALUES (?, ?)", ctx.User.ID, key)

	simple(c, getSimple("/discordtokens"), []message{successMessage{
		T(c, "Your new Discord token is <code>%s</code>. Do not use this if you already verify it.", key),
	}}, nil)
}
