package handler

import (
	"awesomeProject/internal/app/ds"
	"awesomeProject/internal/app/role"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// Register godoc
// @Summary Регистрация пользователя
// @Description Регистрация нового пользователя.
// @Tags Пользователи
// @Accept json
// @Produce json
// @Param request body ds.RegisterReq true "Детали регистрации"
// @Router /api/sign_up [post]
func (h *Handler) SignUp(gCtx *gin.Context) {
	req := &ds.RegisterReq{}

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

	gCtx.JSON(http.StatusOK, &ds.RegisterResp{
		Ok: true,
	})
}

func generateHashString(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// Login godoc
// @Summary Аутентификация пользователя
// @Description Аутентификация пользователя.
// @Tags Пользователи
// @Accept json
// @Produce json
// @Param request body ds.LoginReq true "Детали входа"
// @Success 200 {object} ds.LoginResp "Успешная аутентификация"
// @Failure 400 {object} string "Неверный запрос"
// @Failure 401 {object} string "Неверные учетные данные"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /api/login [post]
func (h *Handler) Login(gCtx *gin.Context) {
	cfg := h.Config
	req := &ds.LoginReq{}

	err := json.NewDecoder(gCtx.Request.Body).Decode(req)
	if err != nil {
		gCtx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	user, err := h.Repository.GetUserByLogin(req.Login)
	if err != nil {
		gCtx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	fmt.Println(user)
	if req.Login == user.UserName && user.UserPassword == generateHashString(req.Password) {
		// значит проверка пройдена
		// генерируем ему jwt
		fmt.Println("ok")
		token := jwt.NewWithClaims(cfg.JWT.SigningMethod, &ds.JWTClaims{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(cfg.JWT.ExpiresIn).Unix(),
				IssuedAt:  time.Now().Unix(),
				Issuer:    "bitop-admin",
			},
			UserID: user.ID,
			Role:   role.Role(user.UserRole),
		})
		fmt.Println(token)

		if token == nil {
			gCtx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("пустой токен"))
			return
		}

		strToken, err := token.SignedString([]byte(cfg.JWT.Token))
		if err != nil {
			gCtx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("невозможно создать строку токена"))
			return
		}

		gCtx.JSON(http.StatusOK, gin.H{
			"ExpiresIn":   cfg.JWT.ExpiresIn,
			"AccessToken": strToken,
			"TokenType":   "Bearer",
			"Role":        user.UserRole,
			"Username":    user.UserName,
		})
	}

	gCtx.AbortWithStatus(http.StatusUnauthorized) // отдаем 403 ответ в знак того что доступ запрещен
}

// Logout godoc
// @Summary Выход пользователя
// @Description Завершение сеанса текущего пользователя.
// @Tags Пользователи
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {string} string "Успешный выход"
// @Failure 400 {object} string "Неверный запрос"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /api/logout [get]
func (h *Handler) Logout(gCtx *gin.Context) {
	// получаем заголовок
	jwtStr := gCtx.GetHeader("Authorization")
	if !strings.HasPrefix(jwtStr, jwtPrefix) { // если нет префикса то нас дурят!
		gCtx.AbortWithStatus(http.StatusBadRequest) // отдаем что нет доступа

		return // завершаем обработку
	}

	// отрезаем префикс
	jwtStr = jwtStr[len(jwtPrefix):]

	_, err := jwt.ParseWithClaims(jwtStr, &ds.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.Config.JWT.Token), nil
	})
	if err != nil {
		gCtx.AbortWithError(http.StatusBadRequest, err)
		log.Println(err)

		return
	}

	// сохраняем в блеклист редиса
	err = h.Redis.WriteJWTToBlacklist(gCtx.Request.Context(), jwtStr, h.Config.JWT.ExpiresIn)
	if err != nil {
		gCtx.AbortWithError(http.StatusInternalServerError, err)

		return
	}

	gCtx.Status(http.StatusOK)
}
