package handler

import (
	"awesomeProject/internal/app/ds"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Get_cosmodroms godoc
// @Summary Получить список космодромов
// @Description Получить список космодромов
// @Tags Полет
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {string} ds.Cosmodroms "Успешно"
// @Failure 400 {object} string "Неверный запрос"
// @Failure 401 {object} string "Неавторизованый пользователь"
// @Failure 403 {object} string "Нет доступа"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /api/flights/cosmodroms [get]
func (h *Handler) get_cosmodroms(c *gin.Context) {
	cosmodroms, err := h.Repository.Select_cosmodroms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": cosmodroms})
}

// put_data_flights godoc
// @Summary Изменить данные о полете
// @Description Изменить данные о полете
// @Tags Полет
// @Accept json
// @Produce json
// @Param request body ds.Flights true "Данные для изменения полета"
// @Security ApiKeyAuth
// @Success 200 {string} ds.Flights "Успешно"
// @Failure 400 {object} string "Неверный запрос"
// @Failure 401 {object} string "Неавторизованый пользователь"
// @Failure 403 {object} string "Нет доступа"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /api/flights [put]
func (h *Handler) put_data_flights(c *gin.Context) {
	flight := ds.Flights{}
	err := c.BindJSON(&flight)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error1", "error": err.Error()})
		return
	}

	_, exists := c.Get("user_role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не найден"})
		return
	}

	res, err2 := h.Repository.Update_flight_date(flight.Id_Application, flight.Id_Ship, flight.Date_Flight, flight.Id_Cosmodrom_Begin, flight.Id_cosmodrom_End)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err2.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "comment": "Дата полета изменена", "flight": res})
	return
}

// Flight_delete godoc
// @Summary Удалить космолет из заявки
// @Description Удалить космолет из заявки
// @Tags Полет
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body ds.delete_flight true "Данные для удаления полета из заявки"
// @Success 200 {string} string "Успешно удалено"
// @Failure 400 {object} string "Неверный запрос"
// @Failure 401 {object} string "Неавторизованый пользователь"
// @Failure 403 {object} string "Нет доступа"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /api/application [delete]
func (h *Handler) delete_flight(c *gin.Context) {
	_, exists := c.Get("user_role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не найден"})
		return
	}

	ship_param := c.Param("id_ship")
	app_param := c.Param("id_application")
	fmt.Println(ship_param, app_param)
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
