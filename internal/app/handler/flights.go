package handler

import (
	"awesomeProject/internal/app/ds"
	"awesomeProject/internal/app/role"
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

// Put_flight_date godoc
// @Summary Установить дату полета
// @Description Установить дату полета
// @Tags Полет
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {string} ds.Flights "Успешно"
// @Failure 400 {object} string "Неверный запрос"
// @Failure 401 {object} string "Неавторизованый пользователь"
// @Failure 403 {object} string "Нет доступа"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /api/flights/date [put]
func (h *Handler) put_flight_date(c *gin.Context) {
	flight := ds.Flights{}
	err := c.BindJSON(&flight)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}

	userRole, exists := c.Get("user_role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не найден"})
		return
	}
	if userRole != role.Buyer {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Изменять информацию о полете может только создатель заявки"})
		return
	}

	res, err2 := h.Repository.Update_flight_date(flight.Id_Application, flight.Id_Ship, flight.Date_Flight)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err2.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "comment": "Дата полета изменена", "flight": res})
	return
}

// Cosmodrom_begin_put godoc
// @Summary Установить космодром вылета
// @Description Установить космодром вылета
// @Tags Полет
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {string} ds.Flights "Успешно"
// @Failure 400 {object} string "Неверный запрос"
// @Failure 401 {object} string "Неавторизованый пользователь"
// @Failure 403 {object} string "Нет доступа"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /api/flights/cosmodrom/begin [put]
func (h *Handler) put_cosmodrom_begin(c *gin.Context) {
	userRole, exists := c.Get("user_role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не найден"})
		return
	}
	if userRole != role.Buyer {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Изменять информацию о полете может только создатель заявки"})
		return
	}

	flight := ds.Flights{}
	err := c.BindJSON(&flight)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}
	res, err2 := h.Repository.Update_cosmodrom_begin(flight.Id_Application, flight.Id_Ship, flight.Id_Cosmodrom_Begin)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err2.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "comment": "Космодром вылета изменен", "flight": res})
	return
}

// Cosmodrom_end_put godoc
// @Summary Установить космодром прилета
// @Description Установить космодром прилета
// @Tags Полет
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {string} ds.Flights "Успешно"
// @Failure 400 {object} string "Неверный запрос"
// @Failure 401 {object} string "Неавторизованый пользователь"
// @Failure 403 {object} string "Нет доступа"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /api/flights/cosmodrom/end [put]
func (h *Handler) put_cosmodrom_end(c *gin.Context) {
	userRole, exists := c.Get("user_role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не найден"})
		return
	}
	if userRole != role.Buyer {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Изменять информацию о полете может только создатель заявки"})
		return
	}

	flight := ds.Flights{}
	err := c.BindJSON(&flight)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}
	res, err2 := h.Repository.Update_cosmodrom_end(flight.Id_Application, flight.Id_Ship, flight.Id_cosmodrom_End)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err2.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "comment": "Космодром прилета изменен", "flight": res})
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
	userRole, exists := c.Get("user_role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не найден"})
		return
	}
	if userRole != role.Buyer {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Удалить полет из заявки может только создатель заявки"})
		return
	}

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
