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
	r.PUT("/application/admin/:id", h.put_application_admin)
	r.PUT("/application/client/:id", h.put_application_client)
	r.PUT("/application/:id", h.put_application)
	r.DELETE("/application/:id", h.delete_application)

	r.GET("/flights/cosmodroms", h.get_cosmodroms)
	r.PUT("/flights/date/application:id_application/ship:id_ship", h.put_flight_date)
	r.PUT("flights/cosmodrom/begin/application:id_application/ship:id_ship", h.put_cosmodrom_begin)
	r.PUT("/flights/cosmodrom/end/application:id_application/ship:id_ship", h.put_cosmodrom_end)
	r.DELETE("/flights/application:id_application/ship:id_ship", h.delete_flight)

	// r.GET("/home", h.ShipsTMPL)
	// r.GET("/home/:id", h.ShipsTMPL)
	// r.POST("/home/:id", h.ShipDelete)

	// r.LoadHTMLGlob("static/templates/*")
	// r.Static("/styles", "./static/css")
	// r.Static("/image", "./static/image")
}

// // метод получения услуг всех, фильтрованных или по айди
// func (h *Handler) ShipsTMPL(c *gin.Context) {
// 	// подробнее об 1 услуге
// 	id_query := c.Param("id")
// 	if id_query != "" {
// 		id, err := strconv.Atoi(id_query)
// 		if err != nil {
// 			return
// 		}
// 		ship, err := h.Repository.GetShipByID(id)
// 		if err != nil {
// 			return
// 		}
// 		c.HTML(http.StatusOK, ShipOne, ship)
// 		return
// 	}
// 	//фильтр всех услуг (поиск)
// 	search := c.Query("search")
// 	ships, err := h.Repository.GetAllShip(search)
// 	if err != nil {
// 		return
// 	}
// 	c.HTML(http.StatusOK, ShipsAll, gin.H{
// 		"Ships":  ships,
// 		"Search": search,
// 	})
// 	return
// }

// // удалить услугу по айди из пост-запроса
// func (h *Handler) ShipDelete(c *gin.Context) {
// 	id_del := c.Param("id")
// 	id, err := strconv.Atoi(id_del)
// 	if err != nil {
// 		return
// 	}
// 	err1 := h.Repository.DeleteShip(id)
// 	if err1 != nil {
// 		return
// 	}
// 	c.Redirect(http.StatusFound, "/home")
// }
