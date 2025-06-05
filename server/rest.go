package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed htmx/*.*
var templatesFS embed.FS

//go:embed htmx/static
var staticFS embed.FS

type PageData struct {
	State         AppData
	SelectedEvent string
}

func getRootPage(c *gin.Context) {
	selectedEventId := c.Query("selectedEvent")

	pd := PageData{
		State:         getState(),
		SelectedEvent: selectedEventId,
	}
	fmt.Println(pd.SelectedEvent)

	c.HTML(http.StatusOK, "index.html", pd)
}

func setupHttpEndpoints() {
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
	// Load templates from the embedded file system
	// We need to create a sub-filesystem for the templates
	tmplFS, err := fs.Sub(templatesFS, "htmx")
	if err != nil {
		log.Fatalf("failed to create sub filesystem for templates: %s", err)
	}
	// Create a new html/template instance and parse the templates
	// We use ParseFS to parse templates directly from an fs.FS
	tmpl := template.Must(template.ParseFS(tmplFS, "*.html"))

	// Set the template engine for Gin to use our parsed templates
	r.SetHTMLTemplate(tmpl)

	// Serve static files from the embedded file system
	// We need to create a sub-filesystem for the static files
	staticSubFS, err := fs.Sub(staticFS, "htmx/static")
	if err != nil {
		log.Fatalf("failed to create sub filesystem for static: %s", err)
	}
	r.StaticFS("/static", http.FS(staticSubFS))

	// all of the main pages
	r.GET("/", getRootPage)
	r.GET("/index.html", getRootPage)
	r.GET("/eventList.html", func(c *gin.Context) {
		selectedEventId := c.Query("selectedEvent")

		pd := PageData{
			State:         getState(),
			SelectedEvent: selectedEventId,
		}
		fmt.Println(pd.SelectedEvent)

		c.HTML(http.StatusOK, "eventList.html", pd)
	})

	r.GET("/fullcalendar.html", func(c *gin.Context) {
		selectedEventId := c.Query("selectedEvent")

		pd := PageData{
			State:         getState(),
			SelectedEvent: selectedEventId,
		}
		fmt.Println(pd.SelectedEvent)

		c.HTML(http.StatusOK, "fullcalendar.html", pd)
	})
	r.GET("/api/fullcalendarEvents.json", func(c *gin.Context) {
		c.JSON(http.StatusOK, getEventsAsFullcalendar())
	})

	api := r.Group("/api")
	{
		// Account endpoints
		api.GET("/accounts", func(c *gin.Context) {
			d := getAccounts()

			c.JSON(http.StatusOK, d)
		})

		api.POST("/accounts", func(c *gin.Context) {
			var account Account
			if err := c.BindJSON(&account); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			newAccount, err := createAccount(account)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusCreated, newAccount)
		})

		// Financial events endpoints
		api.GET("/events", func(c *gin.Context) {
			events := getEvents()

			c.JSON(http.StatusOK, events)
		})

		api.POST("/events", func(c *gin.Context) {
			var incomingEvent Event
			if err := c.BindJSON(&incomingEvent); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			newEvent, err := createEvent(incomingEvent)
			if err != nil {
				c.JSON(http.StatusBadRequest, err)
				return
			}

			c.JSON(http.StatusCreated, newEvent)
		})

		// Financial events endpoints
		api.GET("/state", func(c *gin.Context) {
			localData := getState()

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
	}

	// Serve static files directly
	// r.Static("/", filepath.Join(".", "htmx"))

	log.Fatal(r.Run(":8080"))
}
