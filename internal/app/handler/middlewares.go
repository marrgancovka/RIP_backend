package handler

import (
	"awesomeProject/internal/app/ds"
	"awesomeProject/internal/app/role"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
)

const jwtPrefix = "Bearer "

func (h *Handler) WithAuthCheck(assignedRoles ...role.Role) func(ctx *gin.Context) {
	return func(gCtx *gin.Context) {
		jwtStr := gCtx.GetHeader("Authorization")
		if !strings.HasPrefix(jwtStr, jwtPrefix) { // если нет префикса то нас дурят!
			gCtx.AbortWithStatus(http.StatusForbidden) // отдаем что нет доступа
			gCtx.JSON(http.StatusForbidden, gin.H{"error": "jwt token не найден"})
			gCtx.AbortWithStatus(http.StatusForbidden)
			return // завершаем обработку
		}

		// отрезаем префикс
		jwtStr = jwtStr[len(jwtPrefix):]
		// проверяем jwt в блеклист редиса
		err := h.Redis.CheckJWTInBlacklist(gCtx.Request.Context(), jwtStr)
		if err == nil { // значит что токен в блеклисте
			gCtx.AbortWithStatus(http.StatusForbidden)

			return
		}
		if !errors.Is(err, redis.Nil) { // значит что это не ошибка отсуствия - внутренняя ошибка
			gCtx.AbortWithError(http.StatusInternalServerError, err)

			return
		}

		token, err := jwt.ParseWithClaims(jwtStr, &ds.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(h.Config.JWT.Token), nil
		})
		if err != nil {
			gCtx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			gCtx.AbortWithStatus(http.StatusForbidden)
			return
		}

		myClaims := token.Claims.(*ds.JWTClaims)

		for _, oneOfAssignedRole := range assignedRoles {
			if myClaims.Role == oneOfAssignedRole {
				gCtx.Set("user_id", myClaims.UserID)
				gCtx.Set("user_role", myClaims.Role)
				gCtx.Next()
				return

				// gCtx.AbortWithStatus(http.StatusForbidden)
				// log.Printf("role %s is not assigned in %s", myClaims.Role, assignedRoles)

				// return
			}
		}
		gCtx.AbortWithStatus(http.StatusForbidden)
		log.Printf("role %s is not assigned in %s", myClaims.Role, assignedRoles)
		return

	}

}

func (h *Handler) WithoutAuth(assignedRoles ...role.Role) func(ctx *gin.Context) {
	return func(gCtx *gin.Context) {
		jwtStr := gCtx.GetHeader("Authorization")
		if !strings.HasPrefix(jwtStr, jwtPrefix) {
			return
		}

		jwtStr = jwtStr[len(jwtPrefix):]
		err := h.Redis.CheckJWTInBlacklist(gCtx.Request.Context(), jwtStr)
		if err == nil {
			gCtx.JSON(http.StatusForbidden, err.Error())
			gCtx.AbortWithStatus(http.StatusForbidden)
			return
		}
		if !errors.Is(err, redis.Nil) {
			gCtx.JSON(http.StatusInternalServerError, err.Error())
			gCtx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		token, errParse := jwt.ParseWithClaims(jwtStr, &ds.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(h.Config.JWT.Token), nil
		})
		if errParse != nil {
			gCtx.JSON(http.StatusForbidden, err.Error())
			gCtx.AbortWithStatus(http.StatusForbidden)
			return
		}

		myClaims := token.Claims.(*ds.JWTClaims)
		for _, oneOfAssignedRole := range assignedRoles {
			if myClaims.Role == oneOfAssignedRole {
				gCtx.Set("user_id", myClaims.UserID)
				gCtx.Set("user_role", myClaims.Role)
				gCtx.Next()
				return
			}
		}
		gCtx.AbortWithStatus(http.StatusForbidden)
		gCtx.JSON(http.StatusForbidden, gin.H{"error": "нет доступа"})
		return
	}
}
