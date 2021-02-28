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
	addMessage(c, successMessage{T(c, "Your Discord token is <code>%s</code> . Please note that code are for 1 time validation!", s)})
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
		DiscordGenResp(c)
	}
}

func DiscordGenResp(c *gin.Context, messages ...message) {
	resp(c, 200, "discordtokens.html", &baseTemplateData{
		TitleBar:  "Discord Link Account",
		KyutGrill: "default.jpg",
		Messages:  messages,
		FormData:  normaliseURLValues(c.Request.PostForm),
	})
}