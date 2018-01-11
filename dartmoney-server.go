package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"dartmoney/db"
)

type securityPosition struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
	Shares int    `json:"shares"`
}

func main() {
	if os.Getenv("ALPHAVANTAGE_API_KEY") == "" {
		log.Fatal("Error: environment variable ALPHAVANTAGE_API_KEY is not set")
	}

	log.Printf("Using Gin framework version %s\n", gin.Version)

	db.CreateSchema()

	router := gin.Default()

	router.Static("/css", "./webroot/css")
	router.Static("/js", "./webroot/js")
	router.Static("/html", "./webroot/html")

	router.StaticFile("/", "./webroot/html/index.html")
	router.StaticFile("/favicon.ico", "./webroot/img/favicon.ico")

	router.GET("/api/quotes/:symbols", symbolsQuote)
	router.GET("/api/portfolio/:userId", userPositions)

	router.Run(":8000")
}

func symbolsQuote(c *gin.Context) {
	// Get your free key from https://www.alphavantage.co
	apiKey := os.Getenv("ALPHAVANTAGE_API_KEY")
	const baseURL = "https://www.alphavantage.co/query"
	const fetchTimeout = time.Duration(10 * time.Second)
	symbols := strings.Split(c.Param("symbols"), ",")

	quoteURL := fmt.Sprintf("%s?function=BATCH_STOCK_QUOTES&symbols=%s&apikey=%s",
		baseURL, strings.Join(symbols, ","), apiKey)

	httpClient := http.Client{Timeout: fetchTimeout}
	resp, err := httpClient.Get(quoteURL)
	if err != nil {
		log.Print(err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(fmt.Printf("could not get body when fetching quote: %v\n", err))
		return
	}

	c.Data(http.StatusOK, "application/json", body)
}

func userPositions(c *gin.Context) {
	//userId := c.Param("userId")

	userPositions := make(map[string]securityPosition)
	userPositions["aapl"] = securityPosition{Name: "Apple", Symbol: "aapl", Shares: 255}
	userPositions["baba"] = securityPosition{Name: "Alibaba Group", Symbol: "baba", Shares: 125}
	userPositions["brk.a"] = securityPosition{Name: "Berkshire Hathaway", Symbol: "brk.a", Shares: 18}
	userPositions["sbny"] = securityPosition{Name: "Signature Bank of New York", Symbol: "sbny", Shares: 160}

	c.JSON(http.StatusOK, userPositions)
}
