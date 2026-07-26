package middleware

import (
	"net/http"

	"github.com/DimitriLaPoudre/MusicShack/internal/inbound/http/dto/response"
	"github.com/DimitriLaPoudre/MusicShack/internal/inbound/http/macro"
	"github.com/DimitriLaPoudre/MusicShack/internal/inbound/http/utils"
	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TargetUserMiddleware(user model.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		me, err := utils.GetFromContext[model.User](c, macro.Me)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		targetUser := me

		if targetUserIDRaw := c.Param(macro.UserID); targetUserIDRaw != "" {
			targetUserID, err := uuid.Parse(targetUserIDRaw)
			if err != nil {
				c.JSON(http.StatusBadRequest, response.NewError(err))
				return
			}

			if me.ID != targetUserID {
				if me.Role != model.UserRoleAdmin {
					c.AbortWithStatus(http.StatusForbidden)
				}

				user, err := user.GetUserByFilter(c.Request.Context(), model.UserFilter{ID: &targetUserID})
				if err != nil {
					utils.Error(c, err)
					return
				}

				targetUser = user
			}
		}

		c.Set("target_user", targetUser)

		c.Next()
	}
}
