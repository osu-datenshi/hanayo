package main

import (
	"fmt"
	"database/sql"
	"strconv"
	"strings"
	"regexp"

	"github.com/gin-gonic/gin"
)

func ccreate(c *gin.Context) {
	ccreateResp(c)
}

func ccreateSubmit(c *gin.Context) {
	if getContext(c).User.ID == 0 {
		resp403(c)
		return
	}
	// check registrations are enabled
	if !ccreationEnabled() {
		ccreateResp(c, errorMessage{T(c, "Ow, sorry the clan is not available to create right now ;p")})
		return
	}

	// check name is valid by our criteria
	name := strings.TrimSpace(c.PostForm("name"))
	if !cnameRegex.MatchString(name) {
		ccreateResp(c, errorMessage{T(c, "Your clan can have alphabets, number and these symbols <code>_[]-</code>.")})
		return
	}


	// check whether name already exists
	if db.QueryRow("SELECT 1 FROM clans WHERE name = ?", c.PostForm("name")).
		Scan(new(int)) != sql.ErrNoRows {
		ccreateResp(c, errorMessage{T(c, "Someone already took that clan name... oof.")})
		return
	}

	// check whether tag already exists
	if db.QueryRow("SELECT 1 FROM clans WHERE tag = ?", c.PostForm("tag")).
		Scan(new(int)) != sql.ErrNoRows {
		ccreateResp(c, errorMessage{T(c, "Someone already took that TAG!")})
		return
	}
	
	
	// recaptcha verify

	tag := "0"
		if c.PostForm("tag") != "" {
			tag = c.PostForm("tag")
		}
	
	// The actual registration.

	res, err := db.Exec(`INSERT INTO clans(name, description, icon, tag)
							  VALUES (?, ?, ?, ?);`,
		name, c.PostForm("description"), c.PostForm("icon"), tag)
	if err != nil {
		ccreateResp(c, errorMessage{T(c, "Uh oh... Unexpected error! Clan might be created... I'm not sure though.")})
		fmt.Println(err)
		return
	}
	lid, _ := res.LastInsertId()

	db.Exec("INSERT INTO `user_clans`(user, clan, perms) VALUES (?, ?, 8);", getContext(c).User.ID, lid)



	addMessage(c, successMessage{T(c, "Clan created.")})
	getSession(c).Save()
	c.Redirect(302, "/c/"+strconv.Itoa(int(lid)))
}

func ccreateResp(c *gin.Context, messages ...message) {
	resp(c, 200, "clans/create.html", &baseTemplateData{
		TitleBar:  "Create your clan",
		KyutGrill: "clans.jpg",
		Scripts:   []string{"https://www.google.com/recaptcha/api.js"},
		Messages:  messages,
		FormData:  normaliseURLValues(c.Request.PostForm),
	})
}

func ccreationEnabled() bool {
	var enabled bool
	db.QueryRow("SELECT value_int FROM system_settings WHERE name = 'ccreation_enabled'").Scan(&enabled)
	return enabled
}

// Check User In Query Is Same As User In Y Cookie


func ccin(s string, ss []string) bool {
	for _, x := range ss {
		if x == s {
			return true
		}
	}
	return false
}

var cnameRegex = regexp.MustCompile(`^[A-Za-z0-9 '_\[\]-]{2,15}$`)