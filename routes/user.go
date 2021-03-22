package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nsukmana-dev/restapi/config"
	"github.com/nsukmana-dev/restapi/models"
)

func GetProfile(c *gin.Context) {
	var user models.User
	user_id := int(c.MustGet("jwt_user_id").(float64))

	item := config.DB.Where("id = ?", user_id).Preload("Articles", "user_id = ?", user_id).Find(&user)

	c.JSON(200, gin.H{
		"status": "Berhasil ke profile",
		"data":   item,
	})
}
