package handler

import (
	"awesomeProject/internal/app/ds"
	"awesomeProject/internal/app/role"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Applications godoc
// @Summary Список заявок
// @Description Получение списка заявок
// @Tags Заявки
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param status query string false "Фильтрация по статусу"
// @Param date query string false "Фильтрация по дате начала"
// @Param date_end query string false "Фильтрация по дате конца"
// @Success 200 {string} ds.Application
// @Failure 400 {object} string "Неверный запрос"
// @Failure 401 {object} string "Неавторизованый пользователь"
// @Failure 403 {object} string "Нет доступа"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /api/applications [get]
func (h *Handler) get_applications(c *gin.Context) {
	userID, existsUser := c.Get("user_id")
	userRole, existsRole := c.Get("user_role")
	if !existsUser || !existsRole {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неавторизованный пользователь"})
		return
	}

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

	if userRole == role.Buyer {
		applications, err := h.Repository.Select_applications_buyer(status_query, date, date_end, userID.(uint))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": applications}) // получить все заявки
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

type ApplicationService struct {
	ID             int
	Title          string
	CosmodromBegin string
	CosmodromEnd   string
	Date           time.Time
}

// Application godoc
// @Summary Одна заявка
// @Description Получение заявки с услугами
// @Tags Заявки
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID заявки" Format(int64) default(1)
// @Success 200 {string} ds.Application
// @Failure 400 {object} string "Неверный запрос"
// @Failure 401 {object} string "Неавторизованый пользователь"
// @Failure 403 {object} string "Нет доступа"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /api/application/{id} [get]
func (h *Handler) get_application(c *gin.Context) {
	_, existsUser := c.Get("user_id")
	if !existsUser {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользлватель не авторизован"})
		return
	}

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

// Applications_put_admin godoc
// @Summary Меняем статус заявки на принят или отклонен
// @Description Изменение статуса заявки на принят или отклонен
// @Tags Заявки
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body ds.AppStatus true "Данные для обновления статуса заявки"
// @Success 200 {string} ds.Application
// @Failure 400 {object} string "Неверный запрос"
// @Failure 401 {object} string "Неавторизованый пользователь"
// @Failure 403 {object} string "Нет доступа"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /api/application/admin [put]
func (h *Handler) put_application_admin(c *gin.Context) {
	userId, existsUser := c.Get("user_id")
	userRole, existsUser := c.Get("user_role")
	if !existsUser {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользлватель не авторизован"})
		return
	}
	if userRole == role.Buyer {
		c.JSON(http.StatusForbidden, gin.H{"error": "нет доступа"})
		return
	}
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
	err2 := h.Repository.Update_application_admin(data.ID, data.Status, userId.(uint))
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err2.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "comment": "Статус изменен"})
	return
}

// Applications_put_admin godoc
// @Summary Меняем статус заявки на сформирован
// @Description Изменение статуса заявки на сформирован
// @Tags Заявки
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body ds.AppStatus true "Данные для обновления статуса заявки"
// @Success 200 {string} ds.Application
// @Failure 400 {object} string "Неверный запрос"
// @Failure 401 {object} string "Неавторизованый пользователь"
// @Failure 403 {object} string "Нет доступа"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /api/application/client [put]
func (h *Handler) put_application_client(c *gin.Context) {
	_, existsUser := c.Get("user_id")
	if !existsUser {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользлватель не авторизован"})
		return
	}
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

// Applications_delete godoc
// @Summary Меняем статус заявки на удален
// @Description Изменение статуса заявки на удален
// @Tags Заявки
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID заявки" Format(int64) default(1)
// @Success 200 {string} ds.Application
// @Failure 400 {object} string "Неверный запрос"
// @Failure 401 {object} string "Неавторизованый пользователь"
// @Failure 403 {object} string "Нет доступа"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /api/application/{id} [delete]
func (h *Handler) delete_application(c *gin.Context) {
	_, existsUser := c.Get("user_id")
	if !existsUser {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользлватель не авторизован"})
		return
	}
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
