package handler

import (
	"awesomeProject/internal/app/ds"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Get_ships(c *gin.Context) {
	search := c.Query("search")
	ships, err := h.Repository.Select_ships(search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "succses", "data": ships})
	return
}

func (h *Handler) Get_ship(c *gin.Context) {
	id_param := c.Param("id")
	id, err := strconv.Atoi(id_param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}

	ship, err := h.Repository.Select_ship(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "succses", "data": ship})
	return
}

func (h *Handler) Post_ship(c *gin.Context) {
	var newShip ds.Ship
	err := c.BindJSON(&newShip)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}
	if newShip.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Поле id не должно быть заполнено!"})
		return
	}
	if strings.ReplaceAll(newShip.Title, " ", "") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Название - обязательное поле!"})
		return
	}
	if strings.ReplaceAll(newShip.Type, " ", "") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Тип - обязательное поле!"})
		return
	}
	err2 := h.Repository.Insert_ship(&newShip)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err2.Error()})
	}
	c.JSON(http.StatusCreated, gin.H{"status": "succses"})
}

func (h *Handler) Post_application(c *gin.Context) {
	var newApplication ds.Application
	err := c.BindJSON(&newApplication)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}
	if newApplication.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Поле id не должно быть заполнено!"})
		return
	}
	err2 := h.Repository.Insert_application(&newApplication)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err2.Error()})
	}
	c.JSON(http.StatusCreated, gin.H{"status": "succses"})
}

func (h *Handler) Put_ship(c *gin.Context) {
	var updateShip ds.Ship
	err := c.BindJSON(&updateShip)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}
	if updateShip.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "error": "Такой id не найден"})
		return
	}

	err2 := h.Repository.Update_ship(&updateShip)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err2.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "seccess"})
	return
}

func (h *Handler) Delete_ship(c *gin.Context) {
	id_param := c.Param("id")
	id, err := strconv.Atoi(id_param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
	}
	err2 := h.Repository.Delete_ship(id)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err2.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
