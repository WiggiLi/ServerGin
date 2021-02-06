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

	/*var (
		corsAllowHeaders     = "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token"
		corsAllowMethods     = "HEAD,GET,POST,PUT,DELETE,OPTIONS"
		corsAllowOrigin      = "*"
		corsAllowCredentials = "true"
	)

	ctx.Response.Header.Set("Access-Control-Allow-Credentials", corsAllowCredentials)
	ctx.Response.Header.Set("Access-Control-Allow-Headers", corsAllowHeaders)
	ctx.Response.Header.Set("Access-Control-Allow-Methods", corsAllowMethods)
	ctx.Response.Header.Set("Access-Control-Allow-Origin", corsAllowOrigin)
	*/

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
	/*var (
		corsAllowHeaders     = "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token"
		corsAllowMethods     = "HEAD,GET,POST,PUT,DELETE,OPTIONS"
		corsAllowOrigin      = "*"
		corsAllowCredentials = "true"
	)

	ctx.Response.Header.Set("Access-Control-Allow-Credentials", corsAllowCredentials)
	ctx.Response.Header.Set("Access-Control-Allow-Headers", corsAllowHeaders)
	ctx.Response.Header.Set("Access-Control-Allow-Methods", corsAllowMethods)
	ctx.Response.Header.Set("Access-Control-Allow-Origin", corsAllowOrigin)
	*/
	newEvent := app.NewEvent()

	if err := c.ShouldBindJSON(&newEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//json.Unmarshal(ctx.PostBody(), &newEvent)
	count := server.application.GiveCount(newEvent)
	//fmt.Fprint(c, count)
	c.JSON(http.StatusOK, count)

}

func (server *WebServer) csv_get(c *gin.Context) {
	//ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	//ctx.Response.Header.Set("Access-Control-Allow-Methods", "POST")

	/*payLoad := string(ctx.PostBody())
	log.Print(payLoad)
	newEvent := app.NewEvent()
	json.Unmarshal(ctx.PostBody(), &newEvent)
	count := server.application.GiveCount(newEvent)
	fmt.Fprint(ctx, count)
	*/

	server.application.GiveCsv()

	Openfile, err := os.Open("/tmp/products.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer Openfile.Close()
	//_ = Openfile
	//Send the headers before sending the file
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET")
	c.Header("Content-Disposition", "attachment; filename="+"result.csv")
	c.Header("Content-Type", "text/comma-separated-values")

	//Send the file
	//io.Copy(ctx, Openfile)
	//buf := new(strings.Builder)
	io.Copy(c.Writer, Openfile)
	//b, _ := ioutil.ReadAll(Openfile)
	//s = string(b)
	//s := strings.Replace(string(b), "\r", "\n", -1)

	//c.JSON(http.StatusOK, Openfile)
}

// Start initializes Web Server, starts application and begins serving
func (server *WebServer) Start(errc chan<- error) {
	router := gin.Default()

	port := ":8081"

	router.POST("/comms", server.getAllEvents)
	router.POST("/count", server.getCount)
	router.GET("/get-csv", server.csv_get)

	log.Print("Server is starting on port ", port)
	errc <- router.Run(":8081")
}

// NewWebServer constructs Web Server
func NewWebServer(application app.IncomeRegistration) *WebServer {
	res := &WebServer{}
	res.application = application

	return res
}
