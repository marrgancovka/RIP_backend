package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) get_cosmodroms(c *gin.Context) {
	cosmodroms, err := h.Repository.Select_cosmodroms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": cosmodroms})
}

func (h *Handler) put_flight_date(c *gin.Context) {
	// изменить дату полета
}

func (h *Handler) put_cosmodrom_begin(c *gin.Context) {
	// изменить космодром взлета
}

func (h *Handler) put_cosmodrom_end(c *gin.Context) {
	// изменить космодром приземления
}

func (h *Handler) delete_flight(c *gin.Context) {
	// удалить полет
}
