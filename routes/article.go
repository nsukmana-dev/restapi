package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/nsukmana-dev/restapi/config"
	"github.com/nsukmana-dev/restapi/models"
)

func GetHome(c *gin.Context) {
	items := []models.Article{}
	config.DB.Find(&items)

	c.JSON(200, gin.H{
		"status": "berhasil ke halaman home",
		"data":   items,
	})
}

func GetArticle(c *gin.Context) {
	slug := c.Param("slug")

	var item models.Article

	if config.DB.First(&item, "slug = ?", slug).RecordNotFound() {
		c.JSON(404, gin.H{"status": "error", "message": "record not found"})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"status": "berhasil",
		"data":   item,
	})
}

func PostArticle(c *gin.Context) {
	item := models.Article{
		Title: c.PostForm("title"),
		Desc:  c.PostForm("desc"),
		Slug:  slug.Make(c.PostForm("title")),
	}

	//kalau slugnya sama, maka generate random slug
	//ngecek database apakah sudah ada slug yang sama
	//judul-pertama-stringrandom

	config.DB.Create(&item)

	c.JSON(200, gin.H{
		"status": "berhasil ngepost",
		"data":   item,
	})
}
