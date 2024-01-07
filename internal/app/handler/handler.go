package handler

import (
	"awesomeProject/internal/app/config"
	myminio "awesomeProject/internal/app/myMinio"
	"awesomeProject/internal/app/redis"
	"awesomeProject/internal/app/repository"
	"awesomeProject/internal/app/role"
	"fmt"
	"mime/multipart"

	_ "awesomeProject/docs"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	Config      *config.Config
	Logger      *logrus.Logger
	Repository  *repository.Repository
	MinioClient *minio.Client
	Redis       *redis.Client
}

func New(c *config.Config, l *logrus.Logger, r *repository.Repository, m *minio.Client, redis *redis.Client) *Handler {
	return &Handler{
		Config:      c,
		Logger:      l,
		Repository:  r,
		MinioClient: m,
		Redis:       redis,
	}
}

// иницилизируем запросы
func (h *Handler) Register(r *gin.Engine) {
	r.GET("/api/ships", h.WithoutAuth(role.Admin, role.Buyer), h.Get_ships)
	r.GET("/api/ships/:id", h.WithoutAuth(role.Admin, role.Buyer), h.Get_ship)
	r.POST("/api/ships", h.WithAuthCheck(role.Admin), h.Post_ship)
	r.POST("/api/ships/application", h.WithAuthCheck(role.Buyer, role.Admin), h.Post_application)
	r.PUT("/api/ships", h.WithAuthCheck(role.Admin), h.Put_ship)
	r.PUT("/api/ships/image", h.WithAuthCheck(role.Admin), h.AddImage)
	r.DELETE("/api/ships/:id", h.WithAuthCheck(role.Admin), h.Delete_ship)

	r.GET("/api/applications", h.WithAuthCheck(role.Admin, role.Buyer), h.get_applications)
	r.GET("/api/application/:id", h.WithAuthCheck(role.Admin, role.Buyer), h.get_application)
	r.PUT("/api/application/admin", h.WithAuthCheck(role.Admin), h.put_application_admin)
	r.PUT("/api/application/client", h.WithAuthCheck(role.Buyer, role.Admin), h.put_application_client)
	r.DELETE("/api/application/:id", h.WithAuthCheck(role.Buyer, role.Admin), h.delete_application)

	r.GET("/api/flights/cosmodroms", h.WithoutAuth(role.Admin, role.Buyer), h.get_cosmodroms)
	// r.PUT("/api/flights/date", h.WithAuthCheck(role.Buyer), h.put_flight_date)
	// r.PUT("/api/flights/cosmodrom/begin", h.WithAuthCheck(role.Buyer), h.put_cosmodrom_begin)
	// r.PUT("/api/flights/cosmodrom/end", h.WithAuthCheck(role.Buyer), h.put_cosmodrom_end)
	r.PUT("/api/flights", h.WithAuthCheck(role.Buyer, role.Admin), h.put_data_flights)
	r.DELETE("/api/flights/application/:id_application/:id_ship", h.WithAuthCheck(role.Buyer, role.Admin), h.delete_flight)

	r.LoadHTMLGlob("static/templates/*")
	r.Static("/styles", "./static/css")
	r.Static("/image", "./static/image")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/api/login", h.Login)
	r.POST("/api/sign_up", h.SignUp)
	r.GET("/api/logout", h.Logout)

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
