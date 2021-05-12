package main

import (
	"database/sql"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"net/http"
	"io/ioutil"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/osu-datenshi/api/common"
	"github.com/osu-datenshi/lib/schiavolib"
)

func ip_check(c *gin.Context) {
	raw, err := http.Get(config.IP_API + "/" + clientIP(c) + "/country")
	// Redirect kalau ada ipv6
	if strings.Contains(c.ClientIP(), ":") {
		c.Redirect(302, "/blockedipv6")
	}

	if err != nil {
		panic(err.Error())
	}

	data, err := ioutil.ReadAll(raw.Body)
	if err != nil {
		panic(err.Error())
	}

	if db.QueryRow("SELECT alamat_ip FROM simpen_ip WHERE alamat_ip = ?", clientIP(c)).
        Scan(new(int)) != sql.ErrNoRows {
    } else {
		db.Exec("INSERT INTO simpen_ip(alamat_ip, kode_negara) VALUES (?, ?)", clientIP(c), string(data))
	}
	
	respon_ip(c)
}

func respon_ip(c *gin.Context) {
	kn := "ID"
	var (
		ID int
		alamatIp string
		kodeNegara string
	)
	if !regblockHidup() {
		register(c)
	} else {
		err := db.QueryRow("SELECT id, alamat_ip, kode_negara FROM simpen_ip WHERE alamat_ip = ? AND kode_negara = ?", clientIP(c), kn).Scan(&ID, &alamatIp, &kodeNegara)

		if err != nil {
			simple(c, getSimpleByFilename("register/block.html"), nil, map[string]interface{}{
				"Alamat_ip": clientIP(c),
			})
    	} else {
			register(c)
		}
	}
}

func register(c *gin.Context) {
	if getContext(c).User.ID != 0 {
		resp403(c)
		return
	}
	if c.Query("stopsign") != "1" {
		u, _ := tryBotnets(c)
		if u != "" {
			simple(c, getSimpleByFilename("register/elmo.html"), nil, map[string]interface{}{
				"Username": u,
			})
			return
		}
	}
	registerResp(c)
}

func registerSubmit(c *gin.Context) {

	if getContext(c).User.ID != 0 {
		resp403(c)
		return
	}
	// check registrations are enabled
	if !registrationsEnabled() {
		registerResp(c, errorMessage{T(c, "Sorry, it's not possible to register at the moment. Please try again later.")})
		return
	}
	// check username is valid by our criteria
	username := strings.TrimSpace(c.PostForm("username"))
	if !usernameRegex.MatchString(username) {
		registerResp(c, errorMessage{T(c, "Your username must contain alphanumerical characters, spaces, or any of <code>_[]-</code>")})
		return
	}

	// check whether an username is e.g. cookiezi, shigetora, peppy, wubwoofwolf, loctav
	if in(strings.ToLower(username), forbiddenUsernames) {
		registerResp(c, errorMessage{T(c, "You're not allowed to register with that username.")})
		return
	}

	// check email is valid
	if !govalidator.IsEmail(c.PostForm("email")) {
		registerResp(c, errorMessage{T(c, "Please pass a valid email address.")})
		return
	}

	// passwords check (too short/too common)
	if x := validatePassword(c.PostForm("password")); x != "" {
		registerResp(c, errorMessage{T(c, x)})
		return
	}

	// usernames with both _ and spaces are not allowed
	if strings.Contains(username, "_") && strings.Contains(username, " ") {
		registerResp(c, errorMessage{T(c, "An username or invite code can't contain both underscores and spaces.")})
		return
	}

	// check whether username already exists
	if db.QueryRow("SELECT 1 FROM users WHERE username_safe = ?", safeUsername(username)).
		Scan(new(int)) != sql.ErrNoRows {
		registerResp(c, errorMessage{T(c, "An user with that username already exists!")})
		return
	}

	// check whether an user with that email already exists
	if db.QueryRow("SELECT 1 FROM users WHERE email = ?", c.PostForm("email")).
		Scan(new(int)) != sql.ErrNoRows {
		registerResp(c, errorMessage{T(c, "An user with that email address already exists!")})
		return
	}

	// recaptcha verify
	if config.RecaptchaPrivate != "" && !recaptchaCheck(c) {
		registerResp(c, errorMessage{T(c, "Captcha is invalid.")})
		return
	}

	uMulti, criteria := tryBotnets(c)
	if criteria != "" {
		schiavo.CMs.Send(
			fmt.Sprintf(
				"User **%s** registered with the same %s as %s (%s/u/%s). **POSSIBLE MULTIACCOUNT!!!**. Waiting for ingame verification...",
				username, criteria, uMulti, config.BaseURL, url.QueryEscape(uMulti),
			),
		)
	}

	inpit := c.PostForm("invitecode")

	var ic struct {
		ID   int
		code string
	}

	err := db.QueryRow(`SELECT * FROM invite_code WHERE code = ? LIMIT 1`, strings.TrimSpace(inpit)).Scan(&ic.ID, &ic.code)

	switch {
	case err == sql.ErrNoRows:
		registerResp(c, errorMessage{T(c, "Oops your code is wrong! please contact administrator!")})
		return
	case err != nil:
		c.Error(err)
		resp500(c)
		return
	}
	
	// The actual registration.
	pass, err := generatePassword(c.PostForm("password"))
	if err != nil {
		resp500(c)
		return
	}

	res, err := db.Exec(`INSERT INTO users(username, username_safe, password_md5, salt, email, register_datetime, privileges, password_version)
							  VALUES (?,        ?,             ?,            '',   ?,     ?,                 ?,          2);`,
		username, safeUsername(username), pass, c.PostForm("email"), time.Now().Unix(), common.UserPrivilegePendingVerification)
	if err != nil {
		registerResp(c, errorMessage{T(c, "Whoops, an error slipped in. You might have been registered, though. I don't know.")})
		return
	}
	lid, _ := res.LastInsertId()

	masterStatValues := make([]string, 12)
	masterStatRankValues := nil
	for i:=0; i<3; i++ {
		for j:=0; j<4; i++ {
			mstStatId := (lid-1) * 12 + i * 4 + j + 1
			masterStatValues[i*4+j] = fmt.Sprintf("(%d,%d,%d,%d)",mstStatId,lid,i,j)
			masterStatRankValues = make([]string, 5)
			for k:=0; k<5; k++ {
				masterStatRankValues[k] = fmt.Sprintf("(%d,%d,0)",mstStatId,8-k)
			}
			db.Exec(fmt.Sprintf("insert into `master_stat_ranks` (mst_stat_id, grade_level, grade_count) values %s;",strings.Join(masterStatRankValues,",")))
		}
	}
	masterStatStmt := fmt.Sprintf("insert into `master_stats` (id, user_id, special_mode, game_mode) values %s;",strings.Join(masterStatValues,","))
	res, err = db.Exec(masterStatStmt)
	if err != nil {
		fmt.Println(err)
	}
	schiavo.CMs.Send(fmt.Sprintf("User (**%s** | %s) registered from %s", username, c.PostForm("email"), clientIP(c)))

	setYCookie(int(lid), c)
	logIP(c, int(lid))

	rd.Incr("ripple:registered_users")
	addMessage(c, successMessage{T(c, "You have been successfully registered on Datenshi! You now need to verify your account.")})
	// hapus ip
	db.Exec("DELETE FROM simpen_ip WHERE alamat_ip = ?", clientIP(c))
	getSession(c).Save()
	c.Redirect(302, "/register/verify?u="+strconv.Itoa(int(lid)))
}

func registerResp(c *gin.Context, messages ...message) {
	resp(c, 200, "register/register.html", &baseTemplateData{
		TitleBar:  "Register",
		KyutGrill: "register.jpg",
		Scripts:   []string{"/static/googleapi.js"},
		Messages:  messages,
		FormData:  normaliseURLValues(c.Request.PostForm),
	})
}

func registrationsEnabled() bool {
	var enabled bool
	db.QueryRow("SELECT value_int FROM system_settings WHERE name = 'registrations_enabled'").Scan(&enabled)
	return enabled
}
func regblockHidup() bool {
	var hidup bool
    db.QueryRow("SELECT value_int FROM system_settings WHERE name = 'regblock'").Scan(&hidup)
	return hidup
}

func verifyAccount(c *gin.Context) {
	if getContext(c).User.ID != 0 {
		resp403(c)
		return
	}

	i, ret := checkUInQS(c)
	if ret {
		return
	}

	sess := getSession(c)
	var rPrivileges uint64
	db.Get(&rPrivileges, "SELECT privileges FROM users WHERE id = ?", i)
	if common.UserPrivileges(rPrivileges)&common.UserPrivilegePendingVerification == 0 {
		addMessage(c, warningMessage{T(c, "Nope.")})
		sess.Save()
		c.Redirect(302, "/")
		return
	}

	resp(c, 200, "register/verify.html", &baseTemplateData{
		TitleBar:       "Verify account",
		HeadingOnRight: true,
		KyutGrill:      "welcome.jpg",
	})
}

func BlockerIPV6(c *gin.Context) {
	if getContext(c).User.ID != 0 {
		resp403(c)
		return
	}
	simple(c, getSimpleByFilename("ipv6.html"), nil, map[string]interface{}{
		"IPV6address": c.ClientIP(),
	})
}

func welcome(c *gin.Context) {
	if getContext(c).User.ID != 0 {
		resp403(c)
		return
	}

	i, ret := checkUInQS(c)
	if ret {
		return
	}

	var rPrivileges uint64
	db.Get(&rPrivileges, "SELECT privileges FROM users WHERE id = ?", i)
	if common.UserPrivileges(rPrivileges)&common.UserPrivilegePendingVerification > 0 {
		c.Redirect(302, "/register/verify?u="+c.Query("u"))
		return
	}

	t := T(c, "Welcome!")
	if common.UserPrivileges(rPrivileges)&common.UserPrivilegeNormal == 0 {
		// if the user has no UserNormal, it means they're banned = they multiaccounted
		t = T(c, "Welcome back!")
	}

	resp(c, 200, "register/welcome.html", &baseTemplateData{
		TitleBar:       t,
		HeadingOnRight: true,
		KyutGrill:      "welcome.jpg",
	})
}

// Check User In Query Is Same As User In Y Cookie
func checkUInQS(c *gin.Context) (int, bool) {
	sess := getSession(c)

	i, _ := strconv.Atoi(c.Query("u"))
	y, _ := c.Cookie("y")
	err := db.QueryRow("SELECT 1 FROM identity_tokens WHERE token = ? AND userid = ?", y, i).Scan(new(int))
	if err == sql.ErrNoRows {
		addMessage(c, warningMessage{T(c, "Nope.")})
		sess.Save()
		c.Redirect(302, "/")
		return 0, true
	}
	return i, false
}

func tryBotnets(c *gin.Context) (string, string) {
	var username string

	err := db.QueryRow("SELECT u.username FROM ip_user i INNER JOIN users u ON u.id = i.userid WHERE i.ip = ?", clientIP(c)).Scan(&username)
	if err != nil {
		if err != sql.ErrNoRows {
			c.Error(err)
		}
		return "", ""
	}
	if username != "" {
		return username, "IP"
	}

	cook, _ := c.Cookie("y")
	err = db.QueryRow("SELECT u.username FROM identity_tokens i INNER JOIN users u ON u.id = i.userid WHERE i.token = ?",
		cook).Scan(&username)
	if err != nil {
		if err != sql.ErrNoRows {
			c.Error(err)
		}
		return "", ""
	}
	if username != "" {
		return username, "username"
	}

	return "", ""
}

func in(s string, ss []string) bool {
	for _, x := range ss {
		if x == s {
			return true
		}
	}
	return false
}

var usernameRegex = regexp.MustCompile(`^[A-Za-z0-9 _\[\]-]{2,15}$`)
var forbiddenUsernames = []string{
	"peppy",
	"rrtyui",
	"cookiezi",
	"azer",
	"loctav",
	"banchobot",
	"happystick",
	"doomsday",
	"sharingan33",
	"andrea",
	"cptnxn",
	"reimu-desu",
	"hvick225",
	"_index",
	"my aim sucks",
	"kynan",
	"rafis",
	"sayonara-bye",
	"thelewa",
	"wubwoofwolf",
	"millhioref",
	"tom94",
	"tillerino",
	"clsw",
	"spectator",
	"exgon",
	"axarious",
	"angelsim",
	"recia",
	"nara",
	"emperorpenguin83",
	"bikko",
	"xilver",
	"vettel",
	"kuu01",
	"_yu68",
	"tasuke912",
	"dusk",
	"ttobas",
	"velperk",
	"jakads",
	"jhlee0133",
	"abcdullah",
	"yuko-",
	"entozer",
	"hdhr",
	"ekoro",
	"snowwhite",
	"osuplayer111",
	"musty",
	"nero",
	"elysion",
	"ztrot",
	"koreapenguin",
	"fort",
	"asphyxia",
	"niko",
	"shigetora",
	//update list june 7 2020
	"whitecat",
	"vaxei",
	"alumetri",
	"freddie",
	"benson",
	"mathi",
	"flyingtuna",
	"micca",
	"idke",
	"badeu",
	"abtio",
	"fgsky",
	"karthy",
	"ryuk",
	"spare",
	"asecretbox",
	"o_blue",
	"varvalian",
	"mrekk",
	"dereban",
	"azr8",
	"vahmi",
	"fieryrage",
	"firebat92",
	"jpeg",
	"umbre",
	"andros",
	"bartek22830",
	"mouseeasy",
	"uyghti",
	"chocomint",
	"danyl",
	"dustice",
	"weabole",
	"gasha",
	"maxim",
	"bogdan",
	"morgn",
	"ralidea",
	"merami",
	"okinamo",
	"aricin",
	"nicki minaj",
	"filsdelama",
	"cxu",
	"abyssal",
	"-hebel-",
	"-duckleader-",
	"ephemeral",
	"nameless_ll",
	"1roha_",
	//UPDATE AGUSTUS 4 2020
	"adult",
	"akouka",
	"alkohol",
	"anak haram",
	"anak yatim",
	"analex",
	"anjing",
	"anjink",
	"anjir",
	"arsundal",
	"asu",
	"autis",
	"azizay",
	"babi",
	"babi lu",
	"bacot",
	"bajingan",
	"bajingan tengik",
	"bakka",
	"banci",
	"bandar",
	"bangke",
	"bangsat",
	"bawel",
	"bebon",
	"bedebah",
	"bedon",
	"beer",
	"bego",
	"begok",
	"bencong",
	"Bencong",
	"berak",
	"bercinta",
	"berengsek",
	"bersetubuh",
	"bestiality",
	"betting",
	"biadab",
	"bispak",
	"bitch",
	"blo'on",
	"blowjob",
	"bo’ol",
	"bodo",
	"bodoh",
	"bodooohhh",
	"bokep",
	"boker",
	"bokong",
	"borok",
	"bot",
	"breast",
	"brengsek",
	"brengsex",
	"brengsexxx",
	"buah dada",
	"buah zakar",
	"buaya",
	"buaya darat",
	"budeg",
	"bugil",
	"bujang inam",
	"Burung",
	"cabe lo",
	"cabo",
	"cacat",
	"cacat lu",
	"cangcut",
	"caper",
	"ccookkk",
	"cebong",
	"celaka",
	"celeng",
	"celsit",
	"centil",
	"chealeng",
	"chealshit",
	"cheleng",
	"cheleng tikus",
	"chellenk",
	"chelshit",
	"chleng",
	"cilaka",
	"cina",
	"citikus",
	"citykus",
	"cium pantat ku",
	"cock",
	"coitus",
	"congek",
	"cumi",
	"curut",
	"dafuq",
	"dajal",
	"davit monyet",
	"dekil",
	"dengkulmu",
	"detik",
	"dodol",
	"dongok",
	"downline",
	"drug",
	"drunken",
	"dungu",
	"ectasy",
	"edan",
	"eek",
	"ejakulasi",
	"eneg",
	"enti",
	"erotic",
	"erotis",
	"ewe",
	"fak",
	"fakk",
	"fergay",
	"fetish",
	"foto hot",
	"foto panas",
	"fuck",
	"fuck you",
	"galau-galau munyuk",
	"ganjen",
	"ganyang",
	"gay",
	"geblek",
	"gelandangan",
	"gelo",
	"gembel",
	"gigolo",
	"gila",
	"goblog",
	"goblok",
	"goblokk",
	"gobloookkkkk",
	"grogol grogol",
	"grogol grogol munyak",
	"guk-guk mu",
	"handjob",
	"haram jadah",
	"heunceut",
	"homo",
	"homoseksual",
	"hot girl",
	"hot model",
	"hot picture",
	"hot video",
	"hubungan intim",
	"iblis",
	"idiot",
	"inex",
	"intercouse",
	"intim",
	"itil",
	"jablay",
	"jahanam",
	"jalang",
	"jam8an",
	"jamban",
	"jancok",
	"Jancuk",
	"jancuk",
	"Jangkrik",
	"jangkrik",
	"jaran",
	"Jayus",
	"Jayus lo",
	"jeber bau tai",
	"jelek lu",
	"jembel",
	"jembut",
	"jiancux",
	"jijik gua",
	"judi",
	"kadal",
	"kafir",
	"kalempong",
	"kamasutra",
	"kambing",
	"kampret",
	"kampungan",
	"kancut",
	"kanjut",
	"kaparat",
	"kasian",
	"kasiann",
	"kasihan",
	"katrok",
	"kecepirit",
	"keledai",
	"kentut",
	"keparat",
	"keperawanan",
	"kere",
	"klik",
	"kolor",
	"kondom",
	"kontol",
	"kontol babi",
	"kontol mamakmu",
	"kunyuk",
	"labia",
	"laknat",
	"leak",
	"lebay",
	"lebok",
	"lesbian",
	"live sex",
	"lonte",
	"loserfool",
	"maaampooos",
	"mabuk",
	"madura",
	"maho",
	"makan tai",
	"make love",
	"maling",
	"mamakmu berkontol",
	"mampoooos",
	"mampos",
	"mampus",
	"mampuz",
	"manasturbo",
	"mandul",
	"martole",
	"masturbasi",
	"masturbate",
	"matamu",
	"mati",
	"meki",
	"memek",
	"memek busuk",
	"memex",
	"mesum",
	"miyabi",
	"modar",
	"moddar",
	"monyet",
	"monyong",
	"mother fucker",
	"muka boros",
	"muka plastik",
	"muka sange",
	"munyuk",
	"munyukers",
	"munyuuk",
	"murahan",
	"najis lu",
	"najong",
	"naked",
	"narkoba",
	"narkotika",
	"nekrofilia",
	"ngaceng",
	"ngenti",
	"ngentot",
	"ngepet",
	"ngeseks",
	"ngewe",
	"ngocok",
	"njing",
	"nono",
	"nonok",
	"norak",
	"nude",
	"nyemot",
	"nyolot",
	"nyuk",
	"onani",
	"onta",
	"oral sex",
	"orang kampung",
	"orang udik",
	"orgasm",
	"orgasme",
	"otak melayu",
	"otak miring",
	"otak udang",
	"otakmu",
	"ozawa",
	"paedofilia",
	"palaji",
	"palkon",
	"pantat",
	"pantek",
	"parafilia",
	"payudara",
	"pe’a",
	"pecun",
	"peju",
	"pelacur",
	"pelakor",
	"pele",
	"peler",
	"pemerkosa",
	"pengemis",
	"penis",
	"pepet",
	"perawan",
	"perek",
	"perkosa",
	"playboy",
	"porn",
	"porno",
	"pukimai",
	"pukimak",
	"pulsa",
	"pussy",
	"rape",
	"raping",
	"robot",
	"sakit jiwa",
	"sampah",
	"sarap",
	"seks",
	"selangkangan",
	"self service",
	"sempak",
	"sendal",
	"senggama",
	"setan",
	"settubuh",
	"sex",
	"sex toy",
	"sexual",
	"shit",
	"sial",
	"sialan",
	"sialan lu",
	"silit",
	"sinting",
	"sperma",
	"sundal",
	"taek",
	"taex",
	"tahi",
	"tahi kucing",
	"tai",
	"Tai lo",
	"taik",
	"tampang tua",
	"tampol",
	"telanjang",
	"testicles",
	"testis",
	"tetek",
	"titit",
	"tokai",
	"tokek belang",
	"toket",
	"tolo",
	"tolol",
	"topless",
	"tukang bohong",
	"ublag",
	"udik",
	"united",
	"vagina",
	"video hot",
	"video panas",
	"virginity",
	"vodka",
	"vulva",
	"xxx",
	"zakar",
	"zakar kecil",
	"zoofilia",
}
