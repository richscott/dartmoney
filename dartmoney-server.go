package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type securityPosition struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
	Shares int    `json:"shares"`
}

func main() {
	log.Printf("Using Gin framework version %s\n", gin.Version)

	router := gin.Default()

	router.Static("/css", "./webroot/css")
	router.Static("/js", "./webroot/js")
	router.Static("/html", "./webroot/html")

	router.StaticFile("/", "./webroot/html/index.html")
	router.StaticFile("/favicon.ico", "./webroot/img/favicon.ico")

	router.GET("/api/quote/:symbol", singleSymbolQuote)
	router.GET("/api/portfolio/:userId", userPositions)

	router.Run(":8000")
}

func singleSymbolQuote(c *gin.Context) {
	// Get your free key from https://www.alphavantage.co
	apiKey := os.Getenv("ALPHAVANTAGE_API_KEY")
	baseURL := "https://www.alphavantage.co/query"

	fetchTimeout := time.Duration(10 * time.Second)
	symbol := c.Param("symbol")

	quoteURL := fmt.Sprintf("%s?function=TIME_SERIES_DAILY_ADJUSTED&symbol=%s&outputsize=compact&apikey=%s",
		baseURL, symbol, apiKey)

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
	userPositions["aapl"] = securityPosition{Name: "Apple", Symbol: "aapl", Shares: 250}
	userPositions["baba"] = securityPosition{Name: "Alibaba Group", Symbol: "baba", Shares: 300}
	userPositions["brk-a"] = securityPosition{Name: "Berkshire Hathaway", Symbol: "brk-a", Shares: 40}
	userPositions["sbny"] = securityPosition{Name: "Signature Bank of New York", Symbol: "sbny", Shares: 160}

	c.JSON(http.StatusOK, userPositions)
}
