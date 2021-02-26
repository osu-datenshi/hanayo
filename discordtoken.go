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

	//db.Exec("DELETE FROM discord_tokens WHERE userid = ?", ctx.User.ID)

	var (
		Token     string      `json:"token"`
		Userid    int         `json:"userid"`
		RoleID    interface{} `json:"role_id"`
		Verified  int         `json:"verified"`
		DiscordID interface{} `json:"discord_id"`
	)

	err := db.QueryRow("SELECT token, userid, role_id, verified, discord_id FROM discord_tokens WHERE verified = 1 AND userid = ?", ctx.User.ID).Scan(&Token, &Userid, &RoleID, &Verified, &DiscordID)
	if err != nill {
		key := rs.String(32)
		db.Exec("INSERT INTO discord_tokens(userid, token) VALUES (?, ?)", ctx.User.ID, key)
		simple(c, getSimple("/discordtokens"), []message{successMessage{
			T(c, "Your new Discord token is <code>%s</code>. Do not use this if you already verify it.", key),
		}}, nil)
	} else {
		simple(c, getSimple("/discordtokens"), []message{errorMessage{
			T(c, "You are already verified!"),
		}}, nil)
	}
}