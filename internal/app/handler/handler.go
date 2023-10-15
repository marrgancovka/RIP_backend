package handler

import (
	"awesomeProject/internal/app/repository"

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

// иницилизируем запросы
func (h *Handler) Register(r *gin.Engine) {
	r.GET("/ships", h.Get_ships)
	r.GET("/ships/:id", h.Get_ship)
	r.POST("/ships", h.Post_ship)
	r.POST("/ships/application", h.Post_application)
	r.PUT("/ships", h.Put_ship)
	r.DELETE("/ships/:id", h.Delete_ship)

	r.GET("/applications", h.get_applications)
	r.GET("/applications/:id", h.get_application)
	r.PUT("/application/admin", h.put_application_admin)
	r.PUT("/application/client", h.put_application_client)
	r.DELETE("/application/:id", h.delete_application)

	r.GET("/flights/cosmodroms", h.get_cosmodroms)
	r.PUT("/flights/date", h.put_flight_date)
	r.PUT("flights/cosmodrom/begin", h.put_cosmodrom_begin)
	r.PUT("/flights/cosmodrom/end", h.put_cosmodrom_end)
	r.DELETE("/flights/application:id_application/ship:id_ship", h.delete_flight)

	// r.GET("/home", h.ShipsTMPL)
	// r.GET("/home/:id", h.ShipsTMPL)
	// r.POST("/home/:id", h.ShipDelete)

	// r.LoadHTMLGlob("static/templates/*")
	// r.Static("/styles", "./static/css")
	// r.Static("/image", "./static/image")
}
