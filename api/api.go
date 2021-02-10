package api

import (
	"ServerGin/app"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// WebServer accepts POST requests with payload of XML docs of Receipts
// Then it parses them with XPath and pushes data to Application
type WebServer struct {
	application app.IncomeRegistration
}

func (server *WebServer) getAllEvents(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST")

	newEvent := app.NewEvent()

	if err := c.ShouldBindJSON(&newEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("newEvent ", newEvent.Count)
	events := app.GetEvents()
	events = server.application.GiveEvents(newEvent)

	c.JSON(http.StatusOK, events)
}

func (server *WebServer) getCount(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST")

	newEvent := app.NewEvent()

	if err := c.ShouldBindJSON(&newEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	count := server.application.GiveCount(newEvent)

	c.JSON(http.StatusOK, count)
}

func (server *WebServer) csv_get(c *gin.Context) {
	server.application.GiveCsv()

	Openfile, err := os.Open("/tmp/products.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer Openfile.Close()

	//Send the headers before sending the file
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET")
	c.Header("Content-Disposition", "attachment; filename="+"result.csv")
	c.Header("Content-Type", "text/comma-separated-values")

	io.Copy(c.Writer, Openfile)
}

// Start initializes Web Server, starts application and begins serving
func (server *WebServer) Start(errc chan<- error) {
	router := gin.Default()

	port := ":8081"

	router.POST("/comms", server.getAllEvents)
	router.POST("/count", server.getCount)
	router.GET("/get-csv", server.csv_get)

	//router.Static("style.css", "web/style.css")
	router.StaticFile("/style.css", "web/style.css")
	router.StaticFile("/", "web/index.html")

	log.Print("Server is starting on port ", port)
	errc <- router.Run(":8081")
}

// NewWebServer constructs Web Server
func NewWebServer(application app.IncomeRegistration) *WebServer {
	res := &WebServer{}
	res.application = application

	return res
}
