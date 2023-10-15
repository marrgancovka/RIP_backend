package handler

import (
	myminio "awesomeProject/internal/app/myMinio"
	"awesomeProject/internal/app/repository"
	"fmt"
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go"
	"github.com/sirupsen/logrus"
)

const (
	ShipsAll = "index.tmpl"
	ShipOne  = "second.tmpl"
)

type Handler struct {
	Logger      *logrus.Logger
	Repository  *repository.Repository
	MinioClient *minio.Client
}

func New(l *logrus.Logger, r *repository.Repository, m *minio.Client) *Handler {
	return &Handler{
		Logger:      l,
		Repository:  r,
		MinioClient: m,
	}
}

// иницилизируем запросы
func (h *Handler) Register(r *gin.Engine) {
	r.GET("/ships", h.Get_ships)
	r.GET("/ships/:id", h.Get_ship)
	r.POST("/ships", h.Post_ship)
	r.POST("/ships/application", h.Post_application)
	r.PUT("/ships", h.Put_ship)
	r.DELETE("/ships/:id", h.Delete_ship)

	r.GET("/applications", h.get_applications)
	r.GET("/applications/:id", h.get_application)
	r.PUT("/application/admin", h.put_application_admin)
	r.PUT("/application/client", h.put_application_client)
	r.DELETE("/application/:id", h.delete_application)

	r.GET("/flights/cosmodroms", h.get_cosmodroms)
	r.PUT("/flights/date", h.put_flight_date)
	r.PUT("flights/cosmodrom/begin", h.put_cosmodrom_begin)
	r.PUT("/flights/cosmodrom/end", h.put_cosmodrom_end)
	r.DELETE("/flights/application:id_application/ship:id_ship", h.delete_flight)

	r.LoadHTMLGlob("static/templates/*")
	r.Static("/styles", "./static/css")
	r.Static("/image", "./static/image")
	r.Static("/docs", "/home/margarita/Документы/DevelopmentNetworkApplication_Golang/cmd/main/docs/swagger.json")
}

func (h *Handler) ImageInMinio(file *multipart.File, header *multipart.FileHeader) (string, error) {
	objectName := header.Filename

	if _, err := h.MinioClient.PutObject("vikings-server", objectName, *file, header.Size, minio.PutObjectOptions{
		ContentType: header.Header.Get("Content-Type"),
	}); err != nil {
		return "", err
	}

	return fmt.Sprintf("http://%s/%s/%s", myminio.MinioHost, myminio.BucketName, objectName), nil
}
