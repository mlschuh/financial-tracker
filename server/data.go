package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slices"
	"time"

	"github.com/teambition/rrule-go"
)

type EventOccurance struct {
	ID        string    `json:"id"`
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
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
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

func saveData(data AppData) error {
	file, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(dataFile, file, 0644)
}

// Assume mainData is just a config and no existing state
func calculateState(mainData AppData) AppData {
	// Generate the list of occurances for each event
	// and store them in each event.
	endDate := time.Now().AddDate(4, 0, 0)
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
					// fmt.Println("Got exception")
					// fmt.Println(exception)

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
						ID:        fmt.Sprintf("%s-%s", e.ID, o),
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

func getAccounts() []Account {
	mu.RLock()
	defer mu.RUnlock()

	data, err := loadData()
	if err != nil {
		log.Fatal(err)
	}

	return data.Accounts
}
func createAccount(newAccount Account) (*Account, error) {
	mu.Lock()
	defer mu.Unlock()
	data, err := loadData()
	if err != nil {
		return nil, err
	}

	// Ensure an account with the name doesn't exist already
	for _, a := range data.Accounts {
		if a.Name == newAccount.Name {
			// c.JSON(http.StatusConflict, gin.H{"conflict": "name in use"})
			return nil, fmt.Errorf("name already in use")
		}
	}

	newAccount.ID = pseudo_uuid()
	data.Accounts = append(data.Accounts, newAccount)
	err = saveData(data)
	if err != nil {
		return nil, err
	}

	return &newAccount, nil
}

func getEvents() []Event {
	mu.RLock()
	defer mu.RUnlock()

	data, err := loadData()
	if err != nil {
		log.Fatal(err)
	}

	return data.Events
}

func createEvent(newEvent Event) (*Event, error) {
	mu.Lock()
	defer mu.Unlock()
	data, err := loadData()
	if err != nil {
		log.Fatal(err)
	}

	// Make sure the account exists
	for _, a := range data.Accounts {
		if newEvent.Account == a.ID {
			newEvent.ID = pseudo_uuid()
			data.Events = append(data.Events, newEvent)
			if err := saveData(data); err != nil {
				return nil, err
			}
			return &newEvent, nil
		}
	}

	return nil, fmt.Errorf("account not found")
}

func getState() AppData {
	mu.RLock()
	defer mu.RUnlock()

	data, err := loadData()
	if err != nil {
		log.Fatal(err)
	}

	localData := calculateState(data)

	return localData
}
