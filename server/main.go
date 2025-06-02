// main.go
package main

import (
	"crypto/rand"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/teambition/rrule-go"
)

type EventOccurance struct {
	At        time.Time `json:"date"`
	Amount    int       `json:"amount"`
	EventId   string    `json:"eventId"`
	AccountId string    `json:"accountId"`
	EventType string    `json:"eventType"`
	EventName string    `json:"eventName"`
}

type AccountBalance struct {
	At        time.Time `json:"date"`
	Balance   int       `json:"balance"`
	AccountId string    `json:"accountId"`
	EventId   string    `json:"eventId"`
}

type Account struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Exception struct {
	// Date   time.Time `json:"date"`
	Type   string `json:"type"` //single, //forever
	Amount int    `json:"amount"`
}

type Event struct {
	ID         string               `json:"id"`
	Name       string               `json:"name"`
	Category   string               `json:"category"`
	Account    string               `json:"account"`
	Amount     int                  `json:"amount"`
	Start      time.Time            `json:"start"`
	RRule      string               `json:"rrule"` // RRULE string format
	Type       string               `json:"type"`  // "income" or "expense"
	Exceptions map[string]Exception `json:"exceptions"`
}

type AppData struct {
	Accounts        []Account        `json:"accounts"`
	AccountBalances []AccountBalance `json:"-"`
	Events          []Event          `json:"events"`
	EventOccurances []EventOccurance `json:"-"`
}

var (
	// data     AppData
	dataFile = "financial_data.json"
	mu       sync.RWMutex
)

//go:embed htmx/*
var staticFiles embed.FS

// Filter function using generics
func Filter[T any](slice []T, predicate func(T) bool) []T {
	var result []T
	for _, item := range slice {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

func pseudo_uuid() (uuid string) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	uuid = fmt.Sprintf("%X", b[0:4])

	return
}

func loadData() (data AppData, err error) {
	opendFile, err := os.ReadFile(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			data = AppData{
				Accounts: []Account{},
				Events:   []Event{},
			}
			return data, nil
		}
		return AppData{}, nil
	}
	var localData AppData
	err = json.Unmarshal(opendFile, &localData)
	return localData, err
}

func saveData(data AppData) error {
	file, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(dataFile, file, 0644)
}

func printTimeSlice(ts []time.Time) {
	for _, t := range ts {
		fmt.Println(t)
	}
}

// Assume mainData is just a config and no existing state
func calculateState(mainData AppData) AppData {
	// Generate the list of occurances for each event
	// and store them in each event.
	endDate := time.Now().AddDate(2, 0, 0)
	for _, e := range mainData.Events {
		// separate single events from recurring events
		if e.RRule == "" {
			mainData.EventOccurances = append(mainData.EventOccurances, EventOccurance{
				At:        e.Start,
				Amount:    e.Amount,
				EventId:   e.ID,
				AccountId: e.Account,
				EventType: e.Type,
				EventName: e.Name,
			})
		} else {
			// Generate the rrule for + 2 years
			r, _ := rrule.StrToRRule(e.RRule)
			if r == nil {
				continue
			}
			r.DTStart(e.Start)
			ocurrs := r.Between(
				e.Start,
				endDate, true)

			amount := e.Amount
			for _, o := range ocurrs {
				// Extract the occurance date as that's what's important
				checkDate := o.Format("2006-01-02")

				// If the occurance is an exception,
				// adjust the value here
				exception, ok := e.Exceptions[checkDate]
				if ok {
					fmt.Println("Got exception")
					fmt.Println(exception)

					// Based on the exeption type, we need
					// to process it
					if exception.Type == "forever" {
						amount = exception.Amount
						mainData.EventOccurances = append(mainData.EventOccurances, EventOccurance{
							At:        o,
							Amount:    amount,
							EventId:   e.ID,
							AccountId: e.Account,
							EventType: e.Type,
							EventName: e.Name,
						})
					} else if exception.Type == "single" {
						// Only have the amount for this specific occurance
						mainData.EventOccurances = append(mainData.EventOccurances, EventOccurance{
							At:        o,
							Amount:    exception.Amount,
							EventId:   e.ID,
							AccountId: e.Account,
							EventType: e.Type,
							EventName: e.Name,
						})
					} else if exception.Type == "skip" {
						// If we shouldn't have this occurrence, just don't add it
					}
				} else {
					mainData.EventOccurances = append(mainData.EventOccurances, EventOccurance{
						At:        o,
						Amount:    amount,
						EventId:   e.ID,
						AccountId: e.Account,
						EventType: e.Type,
						EventName: e.Name,
					})
				}
			}
		}
	}

	// Sort the EventOccurances by their date
	// to make looping later easier
	slices.SortFunc(mainData.EventOccurances, func(a, b EventOccurance) int {
		if a.At.Before(b.At) {
			return -1
		} else {
			return 1
		}
	})

	//Now that we have a list of the changes on the accounts,
	// make the running balance list
	for _, a := range mainData.Accounts {
		// Get all of the events associated with this account
		filteredEvents := Filter(mainData.EventOccurances, func(eo EventOccurance) bool { return eo.AccountId == a.ID })

		// Calculate the running totals on the account
		running := 0
		for _, eo := range filteredEvents {
			if eo.EventType == "income" {
				running = running + eo.Amount
			} else if eo.EventType == "expense" {
				running = running - eo.Amount
			}
			mainData.AccountBalances = append(mainData.AccountBalances, AccountBalance{
				At:        eo.At,
				AccountId: a.ID,
				Balance:   running,
				EventId:   eo.EventId,
			})
		}

	}

	return mainData
}

func main() {
	r := gin.Default()

	// Enable CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Prepare the static file system
	// staticFS, err := fs.Sub(staticFiles, ".")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Serve static files directly
	r.StaticFS("/htmx", http.FS(staticFiles))
	// Define routes for static assets
	// r.GET("/", func(c *gin.Context) {
	// 	c.Redirect(http.StatusTemporaryRedirect, "/index.html")
	// })

	// Handle all other static assets using the file server
	// r.GET("/htmx/index.html", func(c *gin.Context) {
	// 	c.FileFromFS("index.html", http.FS(staticFS))
	// })

	// // Handle all other static assets using the file server
	// r.GET("/index.css", func(c *gin.Context) {
	// 	c.FileFromFS("index.css", http.FS(staticFS))
	// })

	// Account endpoints
	r.GET("/accounts", func(c *gin.Context) {
		mu.RLock()
		defer mu.RUnlock()

		data, err := loadData()
		if err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, data.Accounts)
	})

	r.POST("/accounts", func(c *gin.Context) {
		var account Account
		if err := c.BindJSON(&account); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		mu.Lock()
		defer mu.Unlock()
		data, err := loadData()
		if err != nil {
			log.Fatal(err)
		}

		// Ensure an account with the name doesn't exist already
		for _, a := range data.Accounts {
			if a.Name == account.Name {
				c.JSON(http.StatusConflict, gin.H{"conflict": "name in use"})
				return
			}
		}

		account.ID = pseudo_uuid()
		data.Accounts = append(data.Accounts, account)
		if err := saveData(data); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, account)
	})

	// Financial events endpoints
	r.GET("/events", func(c *gin.Context) {
		mu.RLock()
		defer mu.RUnlock()

		data, err := loadData()
		if err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, data.Events)
	})

	r.POST("/events", func(c *gin.Context) {
		var incomingEvent Event
		if err := c.BindJSON(&incomingEvent); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		mu.Lock()
		defer mu.Unlock()
		data, err := loadData()
		if err != nil {
			log.Fatal(err)
		}

		// Make sure the account exists
		for _, a := range data.Accounts {
			if incomingEvent.Account == a.ID {
				incomingEvent.ID = pseudo_uuid()
				data.Events = append(data.Events, incomingEvent)
				if err := saveData(data); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusCreated, incomingEvent)
				return
			}
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": "account not found"})

	})

	// Financial events endpoints
	r.GET("/state", func(c *gin.Context) {
		mu.RLock()
		defer mu.RUnlock()

		data, err := loadData()
		if err != nil {
			log.Fatal(err)
		}

		localData := calculateState(data)

		toReturn := struct {
			EventOccurances []EventOccurance `json:"eventOccurances"`
			AccountBalances []AccountBalance `json:"accountBalances"`
			Events          []Event          `json:"events"`
			Accounts        []Account        `json:"accounts"`
		}{
			EventOccurances: localData.EventOccurances,
			AccountBalances: localData.AccountBalances,
			Events:          localData.Events,
			Accounts:        localData.Accounts,
		}

		c.JSON(http.StatusOK, toReturn)
	})

	log.Fatal(r.Run(":8080"))
}
