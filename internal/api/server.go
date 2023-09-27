package api

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Services struct {
	Name        string
	Modificate  string
	Description string
	Photo       string
}

var data = [...]Services{
	{"STARSHIP MK1", "ldlld", "sre", "/image/rocket.jpg"},
	{"STARSHIP MK2", "ldlld", "sre", "/image/rocket.jpg"},
	{"STARSHIP MK3", "ldlld", "sre", "/image/rocket.jpg"},
	{"STARSHIP MK4", "ldlld", "sre", "/image/rocket.jpg"},
	{"STARSHIP MK5", "ldlld", "sre", "/image/rocket.jpg"},
	{"STARSHIP MK6", "ldlld", "sre", "/image/rocket.jpg"},
	{"STARSHIP MK7", "ldlld", "sre", "/image/rocket.jpg"},
	{"STARSHIP MK8", "ldlld", "sre", "/image/rocket.jpg"},
	{"STARSHIP MK9", "ldlld", "sre", "/image/rocket.jpg"},
	{"STARSHIP MK10", "ldlld", "sre", "/image/rocket.jpg"},
	{"STARSHIP MK11", "ldlld", "sre", "/image/rocket.jpg"},
}

func getItemById(c *gin.Context) {
	name := c.Param("name")
	for _, a := range data {
		if a.Name == name {
			c.HTML(http.StatusOK, "second.tmpl", a)
		}
	}
}

func getItems(c *gin.Context) {
	search := c.Query("search")
	var filterData []Services

	if search != "" {
		for _, a := range data {
			if strings.Contains(strings.ToLower(a.Name), strings.ToLower(search)) {
				filterData = append(filterData, a)
			}
		}
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"Data":   filterData,
			"Search": search,
		})
		return
	}
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"Data":   data,
		"Search": search,
	})
}

func StartServer() {

	log.Println("Server start up")

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.Static("/image", "./static/image")
	r.Static("/styles", "./static/css")

	r.GET("/home", getItems)
	r.GET("/home/:name", getItemById)

	r.Run()
	log.Println("Server down")
}
