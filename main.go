package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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

// type Item struct {
// 	Title     string    `json:"title"`
// 	CreatedAt time.Time `json:"created_at"`
// }

func main() {
	router := gin.Default()

	// 自動的にファイルを返すよう設定 --- (*1)
	router.StaticFS("/templates", http.Dir("templates"))

	// ルートなら /templates/index.html にリダイレクト --- (*2)
	router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(302, "/templates/index.html")
	})

	// フォームの内容を受け取って挨拶する --- (*3)
	router.GET("/hello", func(ctx *gin.Context) {
		prefectures := ctx.Query("prefectures")
		ctx.Header("Content-Type", "text/html; charset=UTF-8")
		ctx.String(200, "<h1>都道府県： "+prefectures+"</h1>")

		resp, err := http.Get("https://opendata.corona.go.jp/api/Covid19JapanAll")
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		// JSONを構造体にエンコード
		var Books Root
		json.Unmarshal(body, &Books)

		fmt.Printf("%-v", Books)
	})

	// サーバーを起動
	err := router.Run("127.0.0.1:8888")
	if err != nil {
		log.Fatal("サーバー起動に失敗", err)
	}
}
