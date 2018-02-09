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

	"github.com/richscott/dartmoney/db"
)

func main() {
	if os.Getenv("ALPHAVANTAGE_API_KEY") == "" {
		log.Fatal("Error: environment variable ALPHAVANTAGE_API_KEY is not set")
	}

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

func buyEquity(c *gin.Context) {

}

func userPositions(c *gin.Context) {
	//userId := c.Param("userId")
	userPositions := db.UserPositions("investor@somewhere")

	c.JSON(http.StatusOK, userPositions)
}
