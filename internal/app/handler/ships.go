package handler

import (
	"awesomeProject/internal/app/ds"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// возвращает список всех космических кораблей
func (h *Handler) Get_ships(c *gin.Context) {
	search := c.Query("search")
	ships, err := h.Repository.Select_ships(search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": ships})
	return
}

// возвращает космический корабль по id из запроса
func (h *Handler) Get_ship(c *gin.Context) {
	id_param := c.Param("id")
	id, err := strconv.Atoi(id_param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ship, err := h.Repository.Select_ship(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "succses", "data": ship})
	return
}

// создает новый космический корабль
func (h *Handler) Post_ship(c *gin.Context) {
	newShip := ds.Ship{}
	err := c.BindJSON(&newShip)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if newShip.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Поле id не должно быть заполнено!"})
		return
	}
	if strings.ReplaceAll(newShip.Title, " ", "") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Название - обязательное поле!"})
		return
	}
	if strings.ReplaceAll(newShip.Type, " ", "") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Тип - обязательное поле!"})
		return
	}
	errRep := h.Repository.Insert_ship(&newShip)
	if errRep != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errRep.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "Создан новый космический корабль"})
}

// добавляем услугу в заявку (или создаем новую заявку и добавляем в нее услугу)
func (h *Handler) Post_application(c *gin.Context) {
	var request struct {
		Id_Ship uint
		Id_user uint
	}
	request.Id_user = 2
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.Id_Ship == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Выберите космический корабль"})
		return
	}

	err2 := h.Repository.Insert_application(&request)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"comment": "Услуга добалена в заявку"})
	return
}

// изменяет данные про космический корабль
func (h *Handler) Put_ship(c *gin.Context) {
	var updateShip ds.Ship
	err := c.BindJSON(&updateShip)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if updateShip.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Такой id не найден"})
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err2 := h.Repository.Update_ship(&updateShip)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	err2 := h.Repository.Delete_ship(id)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// загрузка фото в минио
func (h *Handler) uploadInMinio(file *multipart.File, header *multipart.FileHeader, id string) (string, error) {
	newUrl, errMinio := h.ImageInMinio(file, header)
	if errMinio != nil {
		return "", errMinio
	}
	err := h.Repository.Update_image(id, newUrl)
	if err != nil {
		return "", err
	}
	return newUrl, nil
}

func (h *Handler) AddImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	id := c.Request.FormValue("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Id не найден"})
		return
	}
	if header == nil || header.Size == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Header не найден"})
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	defer func(file multipart.File) {
		errClose := file.Close()
		if errClose != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errClose})
			return
		}
	}(file)

	newUrl, errImage := h.uploadInMinio(&file, header, id)
	if errImage != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errImage})
		return
	}
	c.JSON(http.StatusOK, gin.H{"image_url": newUrl})
}
