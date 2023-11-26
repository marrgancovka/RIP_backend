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
	r.GET("/api/ships", h.Get_ships)
	r.GET("/api/ships/:id", h.Get_ship)
	r.POST("/api/ships", h.Post_ship)
	r.POST("/api/ships/application", h.Post_application)
	r.PUT("/api/ships", h.Put_ship)
	r.PUT("/api/ships/image", h.AddImage)
	r.DELETE("/api/ships/:id", h.Delete_ship)

	r.GET("/api/applications", h.get_applications)
	r.GET("/api/applications/:id", h.get_application)
	r.PUT("/api/application/admin", h.put_application_admin)
	r.PUT("/api/application/client", h.put_application_client)
	r.DELETE("/api/application/:id", h.delete_application)

	r.GET("/api/flights/cosmodroms", h.get_cosmodroms)
	r.PUT("/api/flights/date", h.put_flight_date)
	r.PUT("/apiflights/cosmodrom/begin", h.put_cosmodrom_begin)
	r.PUT("/api/flights/cosmodrom/end", h.put_cosmodrom_end)
	r.DELETE("/api/flights/application:id_application/ship:id_ship", h.delete_flight)

	r.LoadHTMLGlob("static/templates/*")
	r.Static("/styles", "./static/css")
	r.Static("/image", "./static/image")

}

func (h *Handler) ImageInMinio(file *multipart.File, header *multipart.FileHeader) (string, error) {
	objectName := header.Filename

	if _, err := h.MinioClient.PutObject(myminio.BucketName, objectName, *file, header.Size, minio.PutObjectOptions{
		ContentType: header.Header.Get("Content-Type"),
	}); err != nil {
		return "", err
	}

	return fmt.Sprintf("http://%s/%s/%s", myminio.MinioHost, myminio.BucketName, objectName), nil
}
