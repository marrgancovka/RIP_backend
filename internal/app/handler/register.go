package handler

import (
	"awesomeProject/internal/app/ds"
	"awesomeProject/internal/app/role"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type registerReq struct {
	FirstName    string `json:"name"` // лучше назвать то же самое что login
	SecondName   string
	Phone        string
	UserName     string
	UserPassword string `json:"pass"`
}

type registerResp struct {
	Ok bool `json:"ok"`
}

func (h *Handler) SignUp(gCtx *gin.Context) {
	req := &registerReq{}

	err := json.NewDecoder(gCtx.Request.Body).Decode(req)
	if err != nil {
		gCtx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if req.UserPassword == "" {
		gCtx.AbortWithError(http.StatusBadRequest, fmt.Errorf("pass is empty"))
		return
	}

	if req.FirstName == "" {
		gCtx.AbortWithError(http.StatusBadRequest, fmt.Errorf("name is empty"))
		return
	}
	if req.SecondName == "" {
		gCtx.AbortWithError(http.StatusBadRequest, fmt.Errorf("name is empty"))
		return
	}
	if req.UserName == "" {
		gCtx.AbortWithError(http.StatusBadRequest, fmt.Errorf("name is empty"))
		return
	}
	if req.Phone == "" {
		gCtx.AbortWithError(http.StatusBadRequest, fmt.Errorf("name is empty"))
		return
	}

	err = h.Repository.Register(&ds.Users{
		UserRole:     string(role.Buyer),
		FirstName:    req.FirstName,
		SecondName:   req.SecondName,
		Phone:        req.Phone,
		UserName:     req.UserName,
		UserPassword: generateHashString(req.UserPassword), // пароли делаем в хешированном виде и далее будем сравнивать хеши, чтобы их не угнали с базой вместе
	})
	if err != nil {
		gCtx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	gCtx.JSON(http.StatusOK, &registerResp{
		Ok: true,
	})
}

func generateHashString(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
