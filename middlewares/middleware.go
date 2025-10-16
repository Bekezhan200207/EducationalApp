package middlewares

import (
	"go-EdTech/config"
	"go-EdTech/logger"
	"go-EdTech/models"
	"go-EdTech/repositories"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func AuthMiddleware(sessionsRepo *repositories.SessionsRepository, usersRepo *repositories.UsersRepository, rolesRepo *repositories.RoleRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := logger.GetLogger()

		var userUUID uuid.UUID
		var isSessionAuth bool

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// Извлекаем токен, удаляя префикс "Bearer "
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// Разбираем JWT токен, используя секретный ключ
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(config.Config.JwtSecretKey), nil
			})

			// Если токен невалиден, возвращаем ошибку
			if err != nil || !token.Valid {
				logger.Warn("Invalid token", zap.Error(err)) // Логируем предупреждение
				c.JSON(http.StatusUnauthorized, models.NewApiError("invalid token"))
				c.Abort() // Прерываем выполнение дальнейших middleware
				return
			}

			// Извлекаем claims (данные из токена)
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				logger.Warn("Invalid token claims")
				c.JSON(http.StatusUnauthorized, models.NewApiError("invalid token claims"))
				c.Abort()
				return
			}

			// Извлекаем идентификатор пользователя из токена (subject)
			subject, ok := claims["sub"].(string)
			if !ok {
				logger.Warn("Invalid subject in token")
				c.JSON(http.StatusUnauthorized, models.NewApiError("invalid token subject"))
				c.Abort()
				return
			}

			// Преобразуем строку в UUID для идентификатора пользователя
			userUUID, err = uuid.Parse(subject)
			if err != nil {
				logger.Warn("Invalid user ID format in token")
				c.JSON(http.StatusUnauthorized, models.NewApiError("invalid user ID format"))
				c.Abort()
				return
			}
		} else {
			// Если токена нет, пробуем аутентификацию через сессии
			isSessionAuth = true
			sessionToken, err := c.Cookie("session_token")
			if err != nil {
				logger.Warn("No session token found", zap.Error(err))
				c.JSON(http.StatusUnauthorized, models.NewApiError("no session token"))
				c.Abort()
				return
			}

			// Проверяем валидность сессионного токена
			session, _, err := sessionsRepo.GetSession(c.Request.Context(), sessionToken)
			if err != nil {
				logger.Warn("Invalid session token", zap.Error(err))
				c.JSON(http.StatusUnauthorized, models.NewApiError("invalid session token"))
				c.Abort()
				return
			}
			userUUID = session.UserUUID // Извлекаем ID пользователя из сессии
		}

		user, err := usersRepo.FindByUUID(c.Request.Context(), userUUID.String())
		if err != nil {
			logger.Warn("User not found", zap.Error(err))
			c.JSON(http.StatusUnauthorized, models.NewApiError("user not found"))
			c.Abort()
			return
		}

		// Получаем роль пользователя из базы данных
		role, err := rolesRepo.GetRoleByID(c.Request.Context(), user.Role_id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.NewApiError("couldn't find role"))
			c.Abort()
			return
		}

		// Сохраняем информацию о пользователе в контексте для использования в последующих middleware
		c.Set("userID", userUUID)
		c.Set("userRole", role)
		c.Set("isSessionAuth", isSessionAuth)

		logger.Info("User authenticated",
			zap.Any("userID", userUUID),
			zap.String("role", role.Name),
			zap.Bool("isSessionAuth", isSessionAuth))

		c.Next()
	}
}


func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := logger.GetLogger()

		roleObj, exists := c.Get("role")
		if !exists {
			logger.Warn("Role missing - access denied")
			c.JSON(http.StatusForbidden, models.NewApiError("access denied"))
			c.Abort()
			return
		}

		role, ok := roleObj.(*models.Role)
		if !ok {
			logger.Error("Invalid role type in context")
			c.JSON(http.StatusInternalServerError, models.NewApiError("role parsing error"))
			c.Abort()
			return
		}

		// Проверка роли
		if role.Name != requiredRole {
			logger.Warn("Access denied", zap.String("role", role.Name), zap.String("required", requiredRole))
			c.JSON(http.StatusForbidden, models.NewApiError("forbidden"))
			c.Abort()
			return
		}

		logger.Debug("Access granted")
		c.Next()
	}
}