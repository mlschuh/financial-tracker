package main

// The _ "embed" seems to be needed to make gopls happy
import (
	"embed"
	_ "embed"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

//go:embed docs/*.yaml
var swaggerYAMLEmbed []byte // Embed the swagger.yaml file

//go:embed ui_dist/*
var staticFiles embed.FS

type PageData struct {
	State         AppData
	SelectedEvent string
}

// Define constants for CDN URLs
const (
	swaggerUICSS              = "https://cdn.jsdelivr.net/npm/swagger-ui-dist@5.17.1/swagger-ui.css"
	swaggerUIJS               = "https://cdn.jsdelivr.net/npm/swagger-ui-dist@5.17.1/swagger-ui-bundle.js"
	swaggerUIStandalonePreset = "https://cdn.jsdelivr.net/npm/swagger-ui-dist@5.17.1/swagger-ui-standalone-preset.js"
)

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

	// Get the embedded filesystem for the dist directory
	distFS, err := fs.Sub(staticFiles, "ui_dist")
	if err != nil {
		panic(err)
	}

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

		api.DELETE("/events/:id", func(c *gin.Context) {
			eventID := c.Param("id") // Get the event ID from the URL path

			err := deleteEvent(eventID)
			if err != nil {
				// Differentiate between "not found" and other internal errors
				if err.Error() == fmt.Sprintf("event with ID %s not found", eventID) {
					c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				} else {
					log.Printf("Error deleting event %s: %v", eventID, err) // Log the error on the server side
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
				}
				return
			}

			c.Status(http.StatusNoContent) // 204 No Content typically for successful DELETE
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

	// 1. Endpoint to serve the swagger.yaml file
	r.GET("/swagger/doc.yaml", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/x-yaml", swaggerYAMLEmbed)
	})

	// 2. Endpoint to render the Swagger UI page
	r.GET("/swagger", func(c *gin.Context) {
		htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="%s" />
    <link rel="icon" type="image/png" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5.17.1/favicon-32x32.png" sizes="32x32" />
    <link rel="icon" type="image/png" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5.17.1/favicon-16x16.png" sizes="16x16" />
    <style>
        html
        {
            box-sizing: border-box;
            overflow: -moz-scrollbars-vertical;
            overflow-y: scroll;
        }
        *,
        *:before,
        *:after
        {
            box-sizing: inherit;
        }
        body
        {
            margin:0;
            background: #fafafa;
        }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>

    <script src="%s" crossorigin></script>
    <script src="%s" crossorigin></script>
    <script>
        window.onload = function() {
            // Begin Swagger UI call region
            const ui = SwaggerUIBundle({
                url: "/swagger/doc.yaml", // This is the path to your embedded YAML
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout"
            });
            // End Swagger UI call region

            window.ui = ui;
        };
    </script>
</body>
</html>`, swaggerUICSS, swaggerUIJS, swaggerUIStandalonePreset)

		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(htmlContent))
	})

	// SPA handler - serves index.html for all non-API routes
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// Skip API routes (adjust prefix as needed)
		if strings.HasPrefix(path, "/api/") {
			c.JSON(404, gin.H{"error": "Not found"})
			return
		}

		// Check if file exists in embedded FS
		if _, err := distFS.Open(strings.TrimPrefix(path, "/")); err == nil {
			// File exists, serve it
			c.FileFromFS(path, http.FS(distFS))
			return
		}

		// File doesn't exist, serve index.html (SPA fallback)
		c.Header("Content-Type", "text/html")
		indexFile, err := distFS.Open("index.html")
		if err != nil {
			c.String(500, "Could not open index.html")
			return
		}
		defer indexFile.Close()

		http.ServeContent(c.Writer, c.Request, "index.html",
			time.Time{}, indexFile.(io.ReadSeeker))
	})

	log.Fatal(r.Run(":8080"))
}
