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
	simple(c, getSimpleByFilename("discordtokens.html"), []message{successMessage{
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
    
	db.QueryRow("SELECT token, userid, role_id, verified, discord_id FROM discord_tokens WHERE userid = ? AND verified = 1 LIMIT 1", ctx.User.ID).Scan(&Token, &Userid, &RoleID, &Verified, &DiscordID)
	if Verified == 1 {
		simple(c, getSimpleByFilename("discordblock.html"), nil, map[string]interface{}{
			"DiscordID": DiscordID,
		})
        } else {
		resp(c, 200, "discordtokens.html", &baseTemplateData{
			TitleBar:  "Discord Link Account",
			KyutGrill: "default.jpg",
		})
	}
}