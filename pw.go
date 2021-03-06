package main

import (
	"database/sql"
	//"fmt"
	"strings"
	"regexp"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
	"github.com/osu-datenshi/api/common"
	"github.com/osu-datenshi/lib/rs"
	"github.com/alexabrahall/goWebhook"
	"encoding/json"
)

func passwordReset(c *gin.Context) {
	ctx := getContext(c)
	if ctx.User.ID != 0 {
		simpleReply(c, errorMessage{T(c, "You're already logged in!")})
		return
	}

	field := "username"
	if strings.Contains(c.PostForm("username"), "@") {
		field = "email"
	}

	var (
		id         int
		username   string
		email      string
		privileges uint64
	)

	err := db.QueryRow("SELECT id, username, email, privileges FROM users WHERE "+field+" = ?",
		c.PostForm("username")).
		Scan(&id, &username, &email, &privileges)

	switch err {
	case nil:
		// ignore
	case sql.ErrNoRows:
		simpleReply(c, errorMessage{T(c, "That user could not be found.")})
		return
	default:
		c.Error(err)
		resp500(c)
		return
	}

	if common.UserPrivileges(privileges)&
		(common.UserPrivilegeNormal|common.UserPrivilegePendingVerification) == 0 {
		simpleReply(c, errorMessage{T(c, "You look pretty banned/locked here.")})
		return
	}

	// generate key
	key := rs.String(50)

	// TODO: WHY THE FUCK DOES THIS USE USERNAME AND NOT ID PLEASE WRITE MIGRATION
	_, err = db.Exec("INSERT INTO password_recovery(k, u) VALUES (?, ?)", key, username)

	if err != nil {
		c.Error(err)
		resp500(c)
		return
	}

	// KIRIM EMAIL MENGGUNAKAN ZOHO
	sendToZoho(c, username, email, key)

	addMessage(c, successMessage{T(c, "Done! You should shortly receive an email from us at the email you used to sign up on Datenshi.")})
	getSession(c).Save()
	c.Redirect(302, "/")
}

func sendToZoho(c *gin.Context, nicknamedaten, emaildaten, kunciNya string) {
	mailer := gomail.NewMessage()
    mailer.SetHeader("From", config.ZohoSenderName)
    mailer.SetHeader("To", emaildaten)
    mailer.SetHeader("Subject", "Datenshi Password Recovery")
    mailer.SetBody("text/html", "Hey "+nicknamedaten+"! Someone, which we really hope was you, requested a password reset for your account. In case it was you, please click the button below to reset your password on Datenshi. Otherwise, silently ignore this email.<br><a href='"+config.BaseURL+"/pwreset/continue?k="+kunciNya+"'>Click Here for reset your password</a><br><br>If you still having a problem, please reply to this email or go to the our <a href='https://link.troke.id/datenshi'>discord server</a>!")

    dialer := gomail.NewDialer(
        config.ZohoSenderHost,
        config.ZohoSenderPort,
        config.ZohoSenderAuth,
        config.ZohoSenderPassword,
    )

    err := dialer.DialAndSend(mailer)
    if err != nil {
        c.Error(err)
    }
}

func passwordResetContinue(c *gin.Context) {
	k := c.Query("k")

	// todo: check logged in
	if k == "" {
		respEmpty(c, T(c, "Password reset"), errorMessage{T(c, "Nope.")})
		return
	}

	var username string
	switch err := db.QueryRow("SELECT u FROM password_recovery WHERE k = ? LIMIT 1", k).
		Scan(&username); err {
	case nil:
		// move on
	case sql.ErrNoRows:
		respEmpty(c, T(c, "Reset password"), errorMessage{T(c, "That key could not be found. Perhaps it expired?")})
		return
	default:
		c.Error(err)
		resp500(c)
		return
	}

	renderResetPassword(c, username, k)
}

func passwordResetContinueSubmit(c *gin.Context) {
	// todo: check logged in
	var username string
	switch err := db.QueryRow("SELECT u FROM password_recovery WHERE k = ? LIMIT 1", c.PostForm("k")).
		Scan(&username); err {
	case nil:
		// move on
	case sql.ErrNoRows:
		respEmpty(c, T(c, "Reset password"), errorMessage{T(c, "That key could not be found. Perhaps it expired?")})
		return
	default:
		c.Error(err)
		resp500(c)
		return
	}

	p := c.PostForm("password")

	if s := validatePassword(p); s != "" {
		renderResetPassword(c, username, c.PostForm("k"), errorMessage{T(c, s)})
		return
	}

	pass, err := generatePassword(p)
	if err != nil {
		c.Error(err)
		resp500(c)
		return
	}

	_, err = db.Exec("UPDATE users SET password_md5 = ?, salt = '', password_version = '2' WHERE username = ?",
		pass, username)
	if err != nil {
		c.Error(err)
		resp500(c)
		return
	}

	_, err = db.Exec("DELETE FROM password_recovery WHERE k = ? LIMIT 1", c.PostForm("k"))
	if err != nil {
		c.Error(err)
		resp500(c)
		return
	}

	addMessage(c, successMessage{T(c, "All right, we have changed your password and you should now be able to login! Have fun!")})
	getSession(c).Save()
	c.Redirect(302, "/login")
}

func renderResetPassword(c *gin.Context, username, k string, messages ...message) {
	simple(c, getSimpleByFilename("pwreset/continue.html"), messages, map[string]interface{}{
		"Username": username,
		"Key":      k,
	})
}

func generatePassword(p string) (string, error) {
	s, err := bcrypt.GenerateFromPassword([]byte(cmd5(p)), bcrypt.DefaultCost)
	return string(s), err
}

func changePassword(c *gin.Context) {
	ctx := getContext(c)
	if ctx.User.ID == 0 {
		resp403(c)
	}
	s, err := qb.QueryRow("SELECT email FROM users WHERE id = ?", ctx.User.ID)
	if err != nil {
		c.Error(err)
	}

	simple(c, getSimpleByFilename("settings/password.html"), nil, map[string]interface{}{
		"email": s["email"],
	})
}

func gantinamabro(c *gin.Context) {
	ctx := getContext(c)
	if ctx.User.ID == 0 {
    resp403(c)
		return
  }
	if ctx.User.Privileges&common.UserPrivilegeDonor == 0 {
		c.Redirect(302, "/settings")
		simpleReply(c, errorMessage{T(c, "You're not a donor!")})
		return
	}
  s, err := qb.QueryRow("SELECT username FROM users WHERE id = ?", ctx.User.ID)
  if err != nil {
    c.Error(err)
  }

  simple(c, getSimpleByFilename("settings/changename.html"), nil, map[string]interface{}{
    "username": s["username"],
  })
}

type NotifKePubNama struct {
	UserID      int    `json:"userID"`
	NewUsername string `json:"newUsername"`
}

func gantinamabroSubmit(c *gin.Context) {
	ctx := getContext(c)
	if ctx.User.ID == 0 {
		resp403(c)
		return
	}

	if ok, _ := CSRF.Validate(ctx.User.ID, c.PostForm("csrf")); !ok {
    addMessage(c, errorMessage{T(c, "Your session has expired. Please try redoing what you were trying to do.")})
    return
  }

	// check username is valid by our criteria
  gantinama := strings.TrimSpace(c.PostForm("gantinama"))
  if !gantinamaRegex.MatchString(gantinama) {
		c.Redirect(302, "/settings/changename")
  	simpleReply(c, errorMessage{T(c, "Your name must contain alphanumerical characters, spaces, or any of<code>_[]-</code>")})
    return
  }

	if c.PostForm("gantinama") == "" {
		c.Redirect(302, "/settings/changename")
    simpleReply(c, errorMessage{T(c, "Username cannot empty.")})
    return
	}

	if (checkBlacklist(c.PostForm("gantinama")))&&(ctx.User.Privileges&3145727 != 3145727) {
		c.Redirect(302, "/settings/changename")
    simpleReply(c, errorMessage{T(c, "Cannot change into that username.")})
    return
	}

	notifData := NotifKePubNama{
		UserID: ctx.User.ID,
		NewUsername: c.PostForm("gantinama"),
	}

	Notif, err := json.Marshal(&notifData)
	if err != nil {
	    c.Error(err)
	}
	//publish disini bro
	err = rd.Publish("peppy:change_username", string(Notif)).Err()
	if err != nil {
	    c.Error(err)
	}


	var username string
	db.Get(&username, "SELECT username FROM users WHERE id = ?", ctx.User.ID)

	db.Exec("UPDATE users SET username = ?, username_safe = ?  WHERE id = ?", c.PostForm("gantinama"),  safeUsername(c.PostForm("gantinama")), ctx.User.ID)
	// kirim ke discord
	hook := goWebhook.CreateWebhook()
	hook.AddField("Changename","Username : "+username+" has changed their name to "+c.PostForm("gantinama")+" !",true)
	hook.SendWebhook(config.LogDiscord)
	//end of discord

	addMessage(c, successMessage{T(c, "Your name change has been saved!")})
  c.Redirect(302, "/settings/changename")
}


func changePasswordSubmit(c *gin.Context) {
	var messages []message
	ctx := getContext(c)
	if ctx.User.ID == 0 {
		resp403(c)
	}
	defer func() {
		s, err := qb.QueryRow("SELECT email FROM users WHERE id = ?", ctx.User.ID)
		if err != nil {
			c.Error(err)
		}
		simple(c, getSimpleByFilename("settings/password.html"), messages, map[string]interface{}{
			"email": s["email"],
		})
	}()

	if ok, _ := CSRF.Validate(ctx.User.ID, c.PostForm("csrf")); !ok {
		addMessage(c, errorMessage{T(c, "Your session has expired. Please try redoing what you were trying to do.")})
		return
	}

	var password string
	db.Get(&password, "SELECT password_md5 FROM users WHERE id = ? LIMIT 1", ctx.User.ID)

	if err := bcrypt.CompareHashAndPassword(
		[]byte(password),
		[]byte(cmd5(c.PostForm("currentpassword"))),
	); err != nil {
		messages = append(messages, errorMessage{T(c, "Wrong password.")})
		return
	}

	uq := new(common.UpdateQuery)
	uq.Add("email", c.PostForm("email"))
	if c.PostForm("newpassword") != "" {
		if s := validatePassword(c.PostForm("newpassword")); s != "" {
			messages = append(messages, errorMessage{T(c, s)})
			return
		}
		pw, err := generatePassword(c.PostForm("newpassword"))
		if err == nil {
			uq.Add("password_md5", pw)
		}
		sess := getSession(c)
		sess.Set("pw", cmd5(pw))
		sess.Save()
	}
	_, err := db.Exec("UPDATE users SET "+uq.Fields()+" WHERE id = ? LIMIT 1", append(uq.Parameters, ctx.User.ID)...)
	if err != nil {
		c.Error(err)
	}

	db.Exec("UPDATE users SET flags = flags & ~3 WHERE id = ? LIMIT 1", ctx.User.ID)
	messages = append(messages, successMessage{T(c, "Your settings have been saved.")})
}

/* error they said, comment i do
func checkBlacklist(s string) bool {
	bad := false
	rows, err := db.Query("select name, type from name_blacklist where active = 1")
	if err != nil {
	}
	defer rows.Close()
	for rows.Next() {
		var (
			bliName string
			bliType string
		)
		rows.Scan(&bliName, &bliType)
		switch bliType {
		case "split":
			reg := regexp.MustCompile(fmt.Sprintf(`\b%s\b`,bliName))
			bad = bad || reg.MatchString(s)
		case "advanced":
			reg := regexp.MustCompile(bliName)
			bad = bad || reg.MatchString(s)
		case "partial":
			bad = bad || (strings.Contains(s, bliName))
		default:
			bad = bad || (bliName == s)
		}
	}
	return bad
}*/

var gantinamaRegex = regexp.MustCompile(`^[A-Za-z0-9 _\.[\]-]{2,20}$`)
