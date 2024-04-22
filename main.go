package main

import (
	"fmt"
	"main/scraper/scraper"
	"time"
)

func main() {
	go scraper.StartAPI()
	fmt.Println("Server Running On Port 8080")
	fmt.Println("Starting Scraper...")
	for {
		scraper.SystemEngine()
		time.Sleep(time.Second * 60)

	}
	// session, _ := scraper.CreateSession()
	// scraper.FetchInteractions("/interactions/bisoprolol/", session)
}
