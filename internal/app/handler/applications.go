package handler

import (
	"awesomeProject/internal/app/ds"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// возвращает список всех заявок
func (h *Handler) get_applications(c *gin.Context) {
	status_query := c.Query("status")
	date_query := c.Query("date")
	date_end_query := c.Query("date_end")
	if date_query == "" {
		date_query = "0001-01-01"
	}
	if date_end_query == "" {
		date_end_query = "9999-12-31"
	}
	format := "2006-01-02"
	date, err := time.Parse(format, date_query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	date_end, err := time.Parse(format, date_end_query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	applications, err2 := h.Repository.Select_applications(status_query, date, date_end)
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err2.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": applications}) // получить все заявки
	return
}

// возвращает заявку по ID из запроса c услугами
func (h *Handler) get_application(c *gin.Context) {
	id_param := c.Param("id")
	id, err := strconv.Atoi(id_param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}
	application, flights, err2 := h.Repository.Select_application(id)
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "application": application, "flights": flights})
	return
}

// изменяет статус администратора в заявке
func (h *Handler) put_application_admin(c *gin.Context) {
	data := ds.Application{}
	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}
	if data.Status != "accepted" && data.Status != "cancel" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Поменять статус можно только на 'принят' и 'отменен'"})
		return
	}
	err2 := h.Repository.Update_application_admin(data.ID, data.Status)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err2.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "comment": "Статус изменен"})
	return
}

// изменяет статус клиента в заявке
func (h *Handler) put_application_client(c *gin.Context) {
	data := ds.Application{}
	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}
	if data.Status != "formated" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Поменять статус можно только на 'сформирован'"})
		return
	}
	err2 := h.Repository.Update_application_client(data.ID, data.Status)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err2.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "comment": "Статус изменен"})
	return
}

// логически удаляет заявку
func (h *Handler) delete_application(c *gin.Context) {
	id_param := c.Param("id")
	id, err := strconv.Atoi(id_param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
	}
	err2 := h.Repository.Delete_application(id)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err2.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
