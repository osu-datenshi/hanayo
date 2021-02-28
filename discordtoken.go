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
	s := rs.String(32)
	db.Exec("INSERT INTO discord_tokens(userid, token) VALUES (?, ?)", ctx.User.ID, s)
	simple(c, getSimple("/discordtokens"), []message{successMessage{
		T(c, "Your new Discord token is <code>%s</code>. Do not use this if you already verify it.", s),
	}}, nil)
}

func CheckDCToken(c *gin.Context) {
	ctx := getContext(c)
	if ctx.User.ID == 0 {
		resp403(c)
		return
	}
	var (
		Token     string
		Userid    int
		RoleID    interface{}
		Verified  int
		DiscordID interface{}
	)

    d, err := db.QueryRow("SELECT token, userid, role_id, verified, discord_id FROM discord_tokens WHERE verified = 1 AND userid = ? LIMIT 1", ctx.User.ID).Scan(&Token, &Userid, &RoleID, &Verified, &DiscordID)

	if err != nil {
		simple(c, getSimpleByFilename("discordblock.html"), nil, map[string]interface{}{
			"DiscordID": d["discord_id"],
		})
        } else {
		DiscordGenToken(c)
	}
}