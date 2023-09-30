package handler

import (
	"awesomeProject/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Logger     *logrus.Logger
	Repository *repository.Repository
}

func New(l *logrus.Logger, r *repository.Repository) *Handler {
	return &Handler{
		Logger:     l,
		Repository: r,
	}
}

func (h *Handler) Register(r *gin.Engine) {
	r.GET("/home", h.AllShips)
	r.GET("/home/:id", h.ShipById)

	r.LoadHTMLGlob("static/templates/*")
	r.Static("/styles", "./static/css")
	r.Static("/image", "./static/image")
}
