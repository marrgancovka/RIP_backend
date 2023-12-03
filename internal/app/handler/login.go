package handler

import (
	"awesomeProject/internal/app/ds"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

const login = "name"
const password = "12345"
const jwtToken = "test"
const ExpiresIn = time.Hour

var SigningMethod = jwt.SigningMethodHS256

func (h *Handler) Login(gCtx *gin.Context) {
	// cfg := a.config
	req := &ds.LoginReq{}

	err := json.NewDecoder(gCtx.Request.Body).Decode(req)
	if err != nil {
		gCtx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if req.Login == login && req.Password == password {
		// значит проверка пройдена
		// генерируем ему jwt
		token := jwt.NewWithClaims(SigningMethod, &ds.JWTClaims{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(ExpiresIn).Unix(),
				IssuedAt:  time.Now().Unix(),
				Issuer:    "bitop-admin",
			},
			UserUUID: uuid.New(), // test uuid
			// Scopes:   []string{}, // test data
		})

		if token == nil {
			gCtx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("token is nil"))
			return
		}

		strToken, err := token.SignedString([]byte(jwtToken))
		if err != nil {
			gCtx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("cant create str token"))
			return
		}

		gCtx.JSON(http.StatusOK, ds.LoginResp{
			ExpiresIn:   int(ExpiresIn),
			AccessToken: strToken,
			TokenType:   "Bearer",
		})
	}

	gCtx.AbortWithStatus(http.StatusForbidden) // отдаем 403 ответ в знак того что доступ запрещен

}
