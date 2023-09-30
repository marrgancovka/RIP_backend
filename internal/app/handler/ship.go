package handler

import (
	"awesomeProject/internal/app/ds"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// func (h *Handler) Ships(c *gin.Context) {
// 	ship, err := h.Repository.GetAllShip()
// 	if err != nil {
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{
// 		"ships": ship,
// 	})
// }

func (h *Handler) ShipById(c *gin.Context) {
	var data []ds.Ship
	idShip := c.Param("id")
	for _, a := range data {
		id, err := strconv.Atoi(idShip)
		if err != nil {
			return
		}
		if a.ID == uint(id) {
			c.HTML(http.StatusOK, "second.tmpl", a)
		}
	}
}

func (h *Handler) AllShips(c *gin.Context) {
	var data []ds.Ship
	search := c.Query("search")
	if search != "" {
		var filterData []ds.Ship
		for _, a := range data {
			if strings.Contains(strings.ToLower(a.Title), strings.ToLower(search)) {
				filterData = append(filterData, a)
			}
		}
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"Data":   filterData,
			"search": search,
		})
	}
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"Data":   data,
		"Search": search,
	})
}

// func (h *Handler) DeleteShip(c *gin.Context){
// 	var request struct{
// 		ID int `json: "id"`
// 	}
// 	if c.BindJSON(&request) != nil {
// 		return
// 	}
// 	if request.ID == 0{
// 		return
// 	}
// 	if h.Repository.DeleteShip(request.ID) != nil{
// 		return
// 	}
// 	c.Redirect()
// }
