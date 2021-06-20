package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"crypto/tls"

	"github.com/gin-gonic/gin"
	"github.com/osuripple/cheesegull/models"
)

type beatmapPageData struct {
	baseTemplateData

	Found      bool
	Beatmap    models.Beatmap
	Beatmapset models.Set
	SetJSON    string
}

type BeatmapSETV2 struct {
	SetID            int `json:"SetID"`
	ChildrenBeatmaps []struct {
		BeatmapID        int     `json:"BeatmapID"`
		ParentSetID      int     `json:"ParentSetID"`
		DiffName         string  `json:"DiffName"`
		FileMD5          string  `json:"FileMD5"`
		Bpm              int     `json:"BPM"`
		Ar               float64 `json:"AR"`
		Od               int     `json:"OD"`
		Cs               int     `json:"CS"`
		Hp               int     `json:"HP"`
		TotalLength      int     `json:"TotalLength"`
		HitLength        int     `json:"HitLength"`
		Playcount        int     `json:"Playcount"`
		Passcount        int     `json:"Passcount"`
		MaxCombo         int     `json:"MaxCombo"`
		DifficultyRating float64 `json:"DifficultyRating"`
	} `json:"ChildrenBeatmaps"`
	RankedStatus int         `json:"RankedStatus"`
	ApprovedDate interface{} `json:"ApprovedDate"`
	LastUpdate   interface{} `json:"LastUpdate"`
	LastChecked  interface{} `json:"LastChecked"`
	Artist       string      `json:"Artist"`
	Title        string      `json:"Title"`
	Creator      string      `json:"Creator"`
	CreatorID    string      `json:"CreatorID"`
	Source       string      `json:"Source"`
	Tags         string      `json:"Tags"`
	HasVideo     bool        `json:"HasVideo"`
	Genre        int         `json:"Genre"`
	Language     int         `json:"Language"`
	Favourites   int         `json:"Favourites"`
}

func beatmapInfo(c *gin.Context) {
	data := new(beatmapPageData)
	defer resp(c, 200, "beatmap.html", data)

	b := c.Param("bid")
	if _, err := strconv.Atoi(b); err != nil {
		c.Error(err)
	} else {
		data.Beatmap, err = getBeatmapData(b)
		if err != nil {
			c.Error(err)
			return
		}
		data.Beatmapset, err = getBeatmapSetData(data.Beatmap)
		if err != nil {
			c.Error(err)
			return
		}
		sort.Slice(data.Beatmapset.ChildrenBeatmaps, func(i, j int) bool {
			if data.Beatmapset.ChildrenBeatmaps[i].Mode != data.Beatmapset.ChildrenBeatmaps[j].Mode {
				return data.Beatmapset.ChildrenBeatmaps[i].Mode < data.Beatmapset.ChildrenBeatmaps[j].Mode
			}
			return data.Beatmapset.ChildrenBeatmaps[i].DifficultyRating < data.Beatmapset.ChildrenBeatmaps[j].DifficultyRating
		})
	}

	if data.Beatmapset.ID == 0 {
		data.TitleBar = T(c, "Beatmap not found.")
		data.Messages = append(data.Messages, errorMessage{T(c, "Beatmap could not be found.")})
		return
	}

	data.KyutGrill = fmt.Sprintf("https://assets.ppy.sh/beatmaps/%d/covers/cover.jpg?%d", data.Beatmapset.ID, data.Beatmapset.LastUpdate.Unix())
	data.KyutGrillAbsolute = true

	setJSON, err := json.Marshal(data.Beatmapset)
	if err == nil {
		data.SetJSON = string(setJSON)
	} else {
		data.SetJSON = "[]"
	}

	data.TitleBar = T(c, "%s - %s", data.Beatmapset.Artist, data.Beatmapset.Title)
	data.Scripts = append(data.Scripts, "/static/tablesort.js", "/static/beatmap.js")
}

func beatmapSetInfo(c *gin.Context) {
	s := c.Param("bsetid")
	raw, err := http.Get(config.CheesegullAPI + "/s/" + s)
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(raw.Body)
	if err != nil {
        panic(err.Error())
    }
	var data BeatmapSETV2
	json.Unmarshal(body, &data)
	location := fmt.Sprintf("/beatmaps/%d", data.ChildrenBeatmaps.ParentSetID)
    c.Redirect(302, location)
}

func getBeatmapData(b string) (beatmap models.Beatmap, err error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, err := http.Get(config.CheesegullAPI + "/b/" + b)
	if err != nil {
		return beatmap, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return beatmap, err
	}

	err = json.Unmarshal(body, &beatmap)
	if err != nil {
		return beatmap, err
	}

	return beatmap, nil
}

func getBeatmapSetData(beatmap models.Beatmap) (bset models.Set, err error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, err := http.Get(config.CheesegullAPI + "/s/" + strconv.Itoa(beatmap.ParentSetID))
	if err != nil {
		return bset, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return bset, err
	}

	err = json.Unmarshal(body, &bset)
	if err != nil {
		return bset, err
	}

	return bset, nil
}
