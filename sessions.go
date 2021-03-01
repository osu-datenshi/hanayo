package main

import (
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/osu-datenshi/api/common"
	"github.com/osu-datenshi/lib/rs"
)

func sessionInitializer() func(c *gin.Context) {
	return func(c *gin.Context) {
		sess := sessions.Default(c)

		var ctx context

		var passwordChanged bool
		userid := sess.Get("userid")
		if userid, ok := userid.(int); ok {
			ctx.User.ID = userid
			var (
				pRaw     int64
				password string
			)
			err := db.QueryRow("SELECT username, privileges, flags, password_md5 FROM users WHERE id = ?", userid).
				Scan(&ctx.User.Username, &pRaw, &ctx.User.Flags, &password)
			if err != nil {
				c.Error(err)
			}
			if sess.Get("logout") == nil {
				sess.Set("logout", rs.String(15))
			}
			ctx.User.Privileges = common.UserPrivileges(pRaw)
			db.Exec("UPDATE users SET latest_activity = ? WHERE id = ?", time.Now().Unix(), userid)
			if s, ok := sess.Get("pw").(string); !ok || cmd5(password) != s {
				ctx = context{}
				sess.Clear()
				passwordChanged = true
			}
		}

		if v, _ := sess.Get("2fa_must_validate").(bool); !v && ctx.User.ID != 0 {
			tok := sess.Get("token")
			if tok, ok := tok.(string); ok {
				ctx.Token = tok
			}
			oldToken := ctx.Token
			ctx.Token, _ = checkToken(ctx.Token, ctx.User.ID, c)
			// Set rt cookie in case:
			// - User has not got a token in rt
			// - Token has been updated with checkToken
			// - user still has old token in rt
			if x, _ := c.Cookie("rt"); oldToken != ctx.Token || x != ctx.Token {
				http.SetCookie(c.Writer, &http.Cookie{
					Name:    "rt",
					Value:   ctx.Token,
					Expires: time.Now().Add(time.Hour * 24 * 30 * 1),
				})
				sess.Set("token", ctx.Token)
			}
		}

		var addBannedMessage bool
		if ctx.User.ID != 0 && (ctx.User.Privileges&common.UserPrivilegeNormal == 0) {
			ctx = context{}
			sess.Clear()
			addBannedMessage = true
		}

		ctx.Language = getLanguageFromGin(c)

		c.Set("context", ctx)
		c.Set("session", sess)

		if addBannedMessage {
			addMessage(c, warningMessage{T(c, "Akun kamu sudah otomatis dikeluarkan karena akun kamu sepertinya sudah di banned atau locked. Tolong hubungi staff jika ada masalah tertentu melalui Discord.")})
		}
		if passwordChanged {
			addMessage(c, warningMessage{T(c, "Akun kamu sudah otomatis dikeluarkan dengan alasan keamanan. Tolong <a href='/login?redir=%s'>masuk kembali</a>.", url.QueryEscape(c.Request.URL.Path))})
		}

		c.Next()
	}
}

func addMessage(c *gin.Context, m message) {
	sess := getSession(c)
	var messages []message
	messagesRaw := sess.Get("messages")
	if messagesRaw != nil {
		messages = messagesRaw.([]message)
	}
	messages = append(messages, m)
	sess.Set("messages", messages)
}

func getMessages(c *gin.Context) []message {
	sess := getSession(c)
	messagesRaw := sess.Get("messages")
	if messagesRaw == nil {
		return nil
	}
	sess.Delete("messages")
	return messagesRaw.([]message)
}

func getSession(c *gin.Context) sessions.Session {
	return c.MustGet("session").(sessions.Session)
}
