package handler

import (
	"awesomeProject/internal/app/ds"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// возвращает список всех космических кораблей
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

// возвращает космический корабль по id из запроса
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

// создает новый космический корабль
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

	file, header, err3 := c.Request.FormFile("file")
	if header == nil || header.Size == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Не найдет header"})
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err3.Error()})
		return
	}
	defer func(file multipart.File) {
		errLol := file.Close()
		if errLol != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": errLol.Error()})
			return
		}
	}(file)

	// Upload the image to minio server.
	newImageURL, errMinio := h.ImageInMinio(&file, header)
	if errMinio != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": errMinio.Error()})
		return
	}

	newShip.Image_url = newImageURL
	err2 := h.Repository.Insert_ship(&newShip)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err2.Error()})
	}

	c.JSON(http.StatusCreated, gin.H{"status": "succses"})
}

// добавляем услугу в заявку (или создаем новую заявку и добавляем в нее услугу)
func (h *Handler) Post_application(c *gin.Context) {
	var request struct {
		Id_Ship            uint
		Id_Cosmodrom_Begin uint
		Id_cosmodrom_End   uint
		Id_user            uint
		Date_Flight        time.Time
	}
	request.Id_user = 2
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}

	if request.Id_Ship == 0 || request.Id_Cosmodrom_Begin == 0 || request.Id_cosmodrom_End == 0 || request.Date_Flight.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Космический корабль, космодромы и дата полета не могут быть пустыми"})
		return
	}

	err2 := h.Repository.Insert_application(&request)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err2.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "comment": "Услуга добалена в заявку"})
	return
}

// изменяет данные про космический корабль
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

	file, header, err3 := c.Request.FormFile("file")
	if header == nil || header.Size == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Не найдет header"})
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err3.Error()})
		return
	}
	defer func(file multipart.File) {
		errLol := file.Close()
		if errLol != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": errLol.Error()})
			return
		}
	}(file)

	// Upload the image to minio server.
	newImageURL, errMinio := h.ImageInMinio(&file, header)
	if errMinio != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": errMinio.Error()})
		return
	}

	updateShip.Image_url = newImageURL

	err2 := h.Repository.Update_ship(&updateShip)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err2.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "seccess"})
	return
}

// логически удаляет космический корабль
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
