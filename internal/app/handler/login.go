package handler

// import (
// 	"awesomeProject/internal/app/ds"
// 	"awesomeProject/internal/app/role"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/golang-jwt/jwt"
// )

// func (h *Handler) Login(gCtx *gin.Context) {
// 	cfg := h.Config
// 	req := &ds.LoginReq{}

// 	err := json.NewDecoder(gCtx.Request.Body).Decode(req)
// 	if err != nil {
// 		gCtx.AbortWithError(http.StatusBadRequest, err)
// 		return
// 	}

// 	user, err := h.Repository.GetUserByLogin(req.Login)
// 	if err != nil {
// 		gCtx.AbortWithError(http.StatusInternalServerError, err)
// 		return
// 	}

// 	if req.Login == user.UserName && user.UserPassword == generateHashString(req.Password) {
// 		// значит проверка пройдена
// 		// генерируем ему jwt
// 		token := jwt.NewWithClaims(cfg.JWT.SigningMethod, &ds.JWTClaims{
// 			StandardClaims: jwt.StandardClaims{
// 				ExpiresAt: time.Now().Add(cfg.JWT.ExpiresIn).Unix(),
// 				IssuedAt:  time.Now().Unix(),
// 				Issuer:    "bitop-admin",
// 			},
// 			UserID: user.ID,
// 			Role:   role.Role(user.UserRole),
// 		})
// 		if token == nil {
// 			gCtx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("пустой токен"))
// 			return
// 		}

// 		strToken, err := token.SignedString([]byte(cfg.JWT.Token))
// 		if err != nil {
// 			gCtx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("невозможно создать строку токена"))
// 			return
// 		}

// 		gCtx.JSON(http.StatusOK, gin.H{
// 			"ExpiresIn":   cfg.JWT.ExpiresIn,
// 			"AccessToken": strToken,
// 			"TokenType":   "Bearer",
// 		})
// 	}

// 	gCtx.AbortWithStatus(http.StatusForbidden) // отдаем 403 ответ в знак того что доступ запрещен
// }
