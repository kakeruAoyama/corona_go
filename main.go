package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Root struct {
	Value []Item `json:"itemList"`
}

type Item struct {
	Date      string `json:"date"`
	NameJp    string `json:"name_jp"`
	Npatients string `json:"npatients"`
}

func main() {
	router := gin.Default()

	// 自動的にファイルを返すよう設定
	router.StaticFS("/templates", http.Dir("templates"))

	// ルートなら /templates/index.html にリダイレクト
	router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(302, "/templates/index.html")
	})

	// フォームの内容を受け取ってAPIを取得
	router.GET("/hello", func(ctx *gin.Context) {
		prefectures := ctx.Query("prefectures")
		date := ctx.Query("date")
		number_date := strings.Replace(date, "-", "", -1)
		encoded_prefecture := url.QueryEscape(prefectures)
		ctx.Header("Content-Type", "text/html; charset=UTF-8")

		cumulative_value := fetchApi(number_date, encoded_prefecture)
		back_1day := day_back(-1, date)
		back_1day_cumulative_value := fetchApi(back_1day, encoded_prefecture)
		difference_1days := subtract_strings(cumulative_value[0].Npatients, back_1day_cumulative_value[0].Npatients)

		back_7day := day_back(-7, date)
		back_7day_cumulative_value := fetchApi(back_7day, encoded_prefecture)
		difference_7days := subtract_strings(cumulative_value[0].Npatients, back_7day_cumulative_value[0].Npatients)

		ctx.String(200, "<h1>都道府県： "+prefectures+"</h1><h1>累計感染者数："+cumulative_value[0].Npatients+"</h1><h1>前日比増加数："+strconv.Itoa(difference_1days)+"</h1><h1>1週間前比増加数："+strconv.Itoa(difference_7days)+"</h1>")
	})

	// サーバーを起動
	err := router.Run("127.0.0.1:8888")
	if err != nil {
		log.Fatal("サーバー起動に失敗", err)
	}
}

func fetchApi(number_date string, encoded_prefecture string) []Item {
	resp, err := http.Get("https://opendata.corona.go.jp/api/Covid19JapanAll?date=" + number_date + "&dataName=" + encoded_prefecture)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// JSONを構造体にエンコード
	var Item Root
	json.Unmarshal(body, &Item)

	return Item.Value
}

func day_back(back_count int, string_date string) string {
	t, _ := time.Parse("2006-01-02", string_date)
	t2 := t.AddDate(0, 0, back_count)
	t3 := t2.Format("20060102")
	return t3
}

func subtract_strings(before string, after string) int {
	before_int, _ := strconv.Atoi(before)
	after_int, _ := strconv.Atoi(after)
	difference := before_int - after_int
	return difference
}
