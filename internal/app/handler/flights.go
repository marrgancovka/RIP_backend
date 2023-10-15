package handler

import (
	"awesomeProject/internal/app/ds"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) get_cosmodroms(c *gin.Context) {
	cosmodroms, err := h.Repository.Select_cosmodroms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": cosmodroms})
}

func (h *Handler) put_flight_date(c *gin.Context) {
	flight := ds.Flights{}
	err := c.BindJSON(&flight)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}
	err2 := h.Repository.Update_flight_date(flight.Id_Application, flight.Id_Ship, flight.Date_Flight)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err2.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "comment": "Дата полета изменена"})
	return
}

func (h *Handler) put_cosmodrom_begin(c *gin.Context) {
	flight := ds.Flights{}
	err := c.BindJSON(&flight)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}
	err2 := h.Repository.Update_cosmodrom_begin(flight.Id_Application, flight.Id_Ship, flight.Id_Cosmodrom_Begin)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err2.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "comment": "Космодром вылета изменен"})
	return
}

func (h *Handler) put_cosmodrom_end(c *gin.Context) {
	flight := ds.Flights{}
	err := c.BindJSON(&flight)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}
	err2 := h.Repository.Update_cosmodrom_end(flight.Id_Application, flight.Id_Ship, flight.Id_cosmodrom_End)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err2.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "comment": "Космодром прилета изменен"})
	return
}

func (h *Handler) delete_flight(c *gin.Context) {
	ship_param := c.Param("id_ship")
	app_param := c.Param("id_application")
	ship, err := strconv.Atoi(ship_param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
	}
	app, err3 := strconv.Atoi(app_param)
	if err3 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
	}
	err2 := h.Repository.Delete_flight(app, ship)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err2.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "comment": "Полет удален из заявки"})
}
