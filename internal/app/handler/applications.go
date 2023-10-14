package handler

import "github.com/gin-gonic/gin"

func (h *Handler) get_applications(c *gin.Context) {
	// получить все заявки
}

func (h *Handler) get_application(c *gin.Context) {
	// получить  заявкy
}

func (h *Handler) put_application_admin(c *gin.Context) {
	// изменить админа заявки
}

func (h *Handler) put_application_client(c *gin.Context) {
	// изменить клиента заявки
}

func (h *Handler) put_application(c *gin.Context) {
	// изменить заявку
}

func (h *Handler) delete_application(c *gin.Context) {
	// удалить заявку
}
