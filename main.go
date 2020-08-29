package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type res struct {
	Count   int      `json:"count"`
	Results []detail `json:"results"`
	Next    string   `json:"next"`
}

type detail struct {
	IsAdult              bool      `json:"is_adult" gorm:""`
	IsUncensored         bool      `json:"is_uncensored" gorm:""`
	Img                  string    `json:"img" gorm:""`
	DistributedAirTime   string    `json:"distributed_air_time" gorm:""`
	Name                 string    `json:"name" gorm:""`
	ID                   int       `json:"id" gorm:"unique_index;not null"`
	CroppedImg           string    `json:"cropped_img" gorm:""`
	IsViewing            bool      `json:"is_viewing" gorm:""`
	IsLaftelOnly         bool      `json:"is_laftel_only" gorm:""`
	IsDubbed             bool      `json:"is_dubbed" gorm:""`
	Genres               string    `gorm:""`
	GenresList           []string  `json:"genres" gorm:"-"`
	LatestEpisodeCreated time.Time `json:"latest_episode_created" gorm:""`
}

var (
	db   *gorm.DB
	link string
)

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&detail{})

	// years := url.QueryEscape("2020년 1분기")
	// link := fmt.Sprintf("https://laftel.net/api/search/v1/discover/?sort=rank&years=%s&offset=0", years)
	link := fmt.Sprintf("https://laftel.net/api/search/v1/discover/?sort=rank&offset=0")

	for {
		if link != "" {
			link = crawl(link)
		} else {
			break
		}
	}
	fmt.Println("end!")
}

func crawl(link string) string {
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("laftel", "TeJava")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)
	data := &res{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		return ""
	}
	for i, _ := range data.Results {
		data.Results[i].Genres = strings.Join(data.Results[i].GenresList, " ")
		insertItem(&data.Results[i])
	}
	return data.Next
}

func insertItem(d *detail) {
	db.Create(d)
}
