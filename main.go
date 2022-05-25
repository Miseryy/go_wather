package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type W_JSON struct {
	PublicTime          time.Time `json:"publicTime"`
	PublicTimeFormatted string    `json:"publicTimeFormatted"`
	PublishingOffice    string    `json:"publishingOffice"`
	Title               string    `json:"title"`
	Link                string    `json:"link"`
	Description         struct {
		PublicTime          time.Time `json:"publicTime"`
		PublicTimeFormatted string    `json:"publicTimeFormatted"`
		HeadlineText        string    `json:"headlineText"`
		BodyText            string    `json:"bodyText"`
		Text                string    `json:"text"`
	} `json:"description"`
	Forecasts []struct {
		Date      string `json:"date"`
		DateLabel string `json:"dateLabel"`
		Telop     string `json:"telop"`
		Detail    struct {
			Weather string `json:"weather"`
			Wind    string `json:"wind"`
			Wave    string `json:"wave"`
		} `json:"detail"`
		Temperature struct {
			Min struct {
				Celsius    interface{} `json:"celsius"`
				Fahrenheit interface{} `json:"fahrenheit"`
			} `json:"min"`
			Max struct {
				Celsius    string `json:"celsius"`
				Fahrenheit string `json:"fahrenheit"`
			} `json:"max"`
		} `json:"temperature"`
		ChanceOfRain struct {
			T0006 string `json:"T00_06"`
			T0612 string `json:"T06_12"`
			T1218 string `json:"T12_18"`
			T1824 string `json:"T18_24"`
		} `json:"chanceOfRain"`
		Image struct {
			Title  string `json:"title"`
			URL    string `json:"url"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"image"`
	} `json:"forecasts"`
	Location struct {
		Area       string `json:"area"`
		Prefecture string `json:"prefecture"`
		District   string `json:"district"`
		City       string `json:"city"`
	} `json:"location"`
	Copyright struct {
		Title string `json:"title"`
		Link  string `json:"link"`
		Image struct {
			Title  string `json:"title"`
			Link   string `json:"link"`
			URL    string `json:"url"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"image"`
		Provider []struct {
			Link string `json:"link"`
			Name string `json:"name"`
			Note string `json:"note"`
		} `json:"provider"`
	} `json:"copyright"`
}

func GetID(name string) (string, error) {
	s := string([]byte{104, 116, 116, 112, 115, 58, 47, 47, 119, 101, 97, 116, 104, 101, 114, 46, 116, 115, 117, 107, 117, 109, 105, 106, 105, 109, 97, 46, 110, 101, 116, 47, 112, 114, 105, 109, 97, 114, 121, 95, 97, 114, 101, 97, 46, 120, 109, 108})
	res, err := http.Get(s)
	if err != nil {
		return "", err

	}

	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	str_body := string(body)
	array_body := strings.Split(str_body, "\n")

	var id string
	for _, b := range array_body {
		n := strings.Index(b, name)
		if n != -1 {
			id_num := strings.Index(b, "id")
			id = b[id_num+4 : id_num+4+6]
		}
	}

	return id, err

}

func GetW(u string, day int) (string, error) {
	res, err := http.Get(u)
	if err != nil {
		fmt.Println("Get Error")
		return "", err
	}

	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	str_body := string(body)
	// array_body := strings.Split(str_body, "\n")
	var w_j W_JSON
	if err := json.Unmarshal([]byte(str_body), &w_j); err != nil {
		fmt.Println("Json Unmarshal Error")
		fmt.Println(err)
		return "", err

	}

	if len(w_j.Forecasts) < day+1 {
		return "day error", errors.New("day error")

	}

	date := w_j.Forecasts[day].DateLabel
	location := w_j.Location.City
	telop := w_j.Forecasts[day].Telop
	min_celsius := w_j.Forecasts[day].Temperature.Min.Celsius
	max_celsius := w_j.Forecasts[day].Temperature.Max.Celsius

	if min_celsius == nil {
		min_celsius = "--"
	}

	str := fmt.Sprintf("%s %s %s Min%s Max%s\n", date, location, telop, min_celsius, max_celsius)

	return str, err
}

func main() {
	// args := os.Args
	s := string([]byte{230, 157, 177, 228, 186, 172})

	location := flag.String("l", s, "location")
	day_num := flag.Int("d", 0, "get day 0, 1, 2")
	flag.Parse()

	u := string([]byte{104, 116, 116, 112, 115, 58, 47, 47, 119, 101, 97, 116, 104, 101, 114, 46, 116, 115, 117, 107, 117, 109, 105, 106, 105, 109, 97, 46, 110, 101, 116, 47, 97, 112, 105, 47, 102, 111, 114, 101, 99, 97, 115, 116, 47, 99, 105, 116, 121, 47})
	id, err := GetID(*location)
	if err != nil {
		return
	}

	u = u + id
	re, _ := GetW(u, *day_num)

	fmt.Print(re)

}
