package handler

import (
	"awesomeProject/internal/app/ds"
	"awesomeProject/internal/app/repository"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	ShipsAll = "index.tmpl"
	ShipOne  = "second.tmpl"
)

type Handler struct {
	Logger     *logrus.Logger
	Repository *repository.Repository
}

func New(l *logrus.Logger, r *repository.Repository) *Handler {
	return &Handler{
		Logger:     l,
		Repository: r,
	}
}

func (h *Handler) Register(r *gin.Engine) {
	r.GET("/home", h.ShipsTMPL)
	r.GET("/home/:id", h.ShipsTMPL)
	r.POST("/home/:id", h.ShipDelete)

	r.LoadHTMLGlob("static/templates/*")
	r.Static("/styles", "./static/css")
	r.Static("/image", "./static/image")
}

// метод получения услуг всех, фильтрованных или по айди
func (h *Handler) ShipsTMPL(c *gin.Context) {
	// подробнее об 1 услуге
	id_query := c.Param("id")
	if id_query != "" {
		id, err := strconv.Atoi(id_query)
		if err != nil {
			return
		}
		ship, err := h.Repository.GetShipByID(id)
		if err != nil {
			return
		}
		c.HTML(http.StatusOK, ShipOne, ship)
		return
	}
	//фильтр всех услуг (поиск)
	search := c.Query("search")
	if search != "" {
		var filterData []ds.Ship
		ships, err := h.Repository.GetAllShip()
		if err != nil {
			return
		}
		for _, a := range *ships {
			if strings.Contains(strings.ToLower(search), strings.ToLower(a.Title)) {
				filterData = append(filterData, a)
			}
		}
		c.HTML(http.StatusOK, ShipsAll, gin.H{
			"Ships":  ships,
			"Search": search,
		})
	}

	//вывод всех услуг

	ships, err := h.Repository.GetAllShip()
	if err != nil {
		return
	}
	c.HTML(http.StatusOK, ShipsAll, gin.H{
		"Ships": ships,
	})
}

func (h *Handler) ShipDelete(c *gin.Context) {
	id_del := c.Param("id")
	id, err := strconv.Atoi(id_del)
	if err != nil {
		return
	}
	err1 := h.Repository.DeleteShip(id)
	if err1 != nil {
		return
	}
	c.Redirect(http.StatusFound, "/home")
}
