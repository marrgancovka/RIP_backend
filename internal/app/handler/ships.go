package handler

import (
	"awesomeProject/internal/app/ds"
	"awesomeProject/internal/app/role"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Ships godoc
// @Summary Список кораблей
// @Description Получение списка кораблей
// @Tags Корабли
// @Accept json
// @Produce json
// @Param search query string false "Фильтрация поиска"
// @Success 200 {string} ds.Ship
// @Failure 400 {object} string "Неверный запрос"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /api/ships [get]
func (h *Handler) Get_ships(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		fmt.Println("not exists")
		search := c.Query("search")
		ships, app, err := h.Repository.Select_ships(search, 0)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": ships, "app": app})
		return
	}
	search := c.Query("search")
	ships, app, err := h.Repository.Select_ships(search, userId.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": ships, "app": app})
	return
}

// Get_ship godoc
// @Summary Получение информации о корабле
// @Description Получение информации о корабле по его идентификатору
// @Tags Корабли
// @Accept json
// @Produce json
// @Param id path int true "ID корабля" Format(int64) default(1)
// @Success 200 {object} map[string]interface{} "Успешный запрос"
// @Failure 400 {object} string "Неверный запрос"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /api/ships/{id} [get]
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

// Post_ship godoc
// @Summary Создание корабля
// @Security ApiKeyAuth
// @Tags Корабли
// @Description Создание космического корабля
// @Accept json
// @Security ApiKeyAuth
// @Produce json
// @Param newShip body ds.Ship true "Информация о новом корабле"
// @Success 201 {object} ds.Ship "Успешное создание космического корабля"
// @Failure 400 {object} string "Неверный запрос"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /api/ships [post]
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
	newid, errRep := h.Repository.Insert_ship(&newShip)
	if errRep != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errRep.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "Создан новый космический корабль", "id": newid})
}

// Post_application godoc
// @Summary Добавление услуги в заявку или создание новой заявки и добавление в нее услуги
// @Tags Корабли
// @Description Создание или обновление заявки с добавлением услуги
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body ds.ShipToAppReq true "Данные для добавления корабля в заявку"
// @Success 200 {object} string "Успешное добавление услуги в заявку"
// @Failure 400 {object} string "Неверный запрос"
// @Failure 401 {object} string "Ошибка авторизации"
// @Failure 403 {object} string "Доступ запрещен"
// @Router /api/ships/application [post]
func (h *Handler) Post_application(c *gin.Context) {
	var request struct {
		ShipId uint `json:"id_ship"`
	}

	userID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не найден"})
		return
	}
	fmt.Println(userID)

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID пользователя не может быть пустым"})
		return
	}
	fmt.Println(userIDUint)

	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.ShipId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Выберите космический корабль"})
		return
	}

	err2 := h.Repository.Insert_application(userIDUint, request.ShipId)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"comment": "Услуга добалена в заявку"})
	return
}

// Put_ship godoc
// @Summary Изменение информации о корабле
// @Security ApiKeyAuth
// @Tags Корабли
// @Description Изменение информации о корабле
// @Security ApiKeyAuth
// @Produce json
// @Param updateShip body ds.Ship true "Обновленная информация о корабле"
// @Success 201 {object} ds.Ship "Успешное создание космического корабля"
// @Failure 400 {object} string "Неверный запрос"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /api/ships [put]
func (h *Handler) Put_ship(c *gin.Context) {
	var updateShip ds.Ship
	err := c.BindJSON(&updateShip)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error1": err.Error()})
		return
	}
	if updateShip.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Такой id не найден"})
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error2": err.Error()})
		return
	}
	err2 := h.Repository.Update_ship(&updateShip)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error3": err2.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "seccess"})
	return
}

// Delete_ship godoc
// @Summary Удаление космического корабля
// @Security ApiKeyAuth
// @Tags Корабли
// @Description Удаление космического корабля
// @Security ApiKeyAuth
// @Produce json
// @Success 201 {object} ds.Ship "Успешное удаление космического корабля"
// @Failure 400 {object} string "Неверный запрос"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /api/ships/{id} [delete]
func (h *Handler) Delete_ship(c *gin.Context) {
	userRole, exists := c.Get("user_role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не найден"})
		return
	}
	if userRole != role.Admin {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Доступ запрещен"})
		return
	}

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

// AddImage godoc
// @Summary Загрузка изображения для корабля
// @Security ApiKeyAuth
// @Tags Корабли
// @Description Загрузка изображения для корабля
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Изображение в формате файла"
// @Param id formData string true "Идентификатор корабля"
// @Success 200 {object} ds.Ship "Успешно добавлено фото"
// @Failure 400 {object} string "Неверный запрос"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /api/ships/image [put]
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
