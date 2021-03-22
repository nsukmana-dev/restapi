package routes

import (
	"strconv"
	"time"

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

	var olditem models.Article
	var slug = slug.Make(c.PostForm("title"))

	if !config.DB.First(&olditem, "slug = ?", slug).RecordNotFound() {
		slug = slug + "-" + strconv.FormatInt(time.Now().Unix(), 10)
	}

	item := models.Article{
		Title:  c.PostForm("title"),
		Desc:   c.PostForm("desc"),
		Tag:    c.PostForm("tag"),
		Slug:   slug,
		UserID: uint(c.MustGet("jwt_user_id").(float64)),
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

func GetArticleByTag(c *gin.Context) {
	tag := c.Param("tag")
	items := []models.Article{}

	config.DB.Where("tag LIKE ?", "%"+tag+"%").Find(&items)

	c.JSON(200, gin.H{"data": items})
}

func UpdateArticle(c *gin.Context) {
	id := c.Param("id")

	var item models.Article

	if config.DB.First(&item, "id = ?", id).RecordNotFound() {
		c.JSON(404, gin.H{"status": "error", "message": "record not found"})
		c.Abort()
		return
	}

	if uint(c.MustGet("jwt_user_id").(float64)) != item.UserID {
		c.JSON(403, gin.H{"status": "error", "message": "this data is forbidden"})
		c.Abort()
		return
	}

	config.DB.Model(&item).Where("id = ?", id).Updates(models.Article{
		Title: c.PostForm("title"),
		Desc:  c.PostForm("desc"),
		Tag:   c.PostForm("tag"),
	})

	c.JSON(200, gin.H{
		"status": "berhasil update",
		"data":   item,
	})
}

func DeleteArticle(c *gin.Context) {
	id := c.Param("id")
	var article models.Article

	config.DB.Where("id = ?", id).Delete(&article)
	c.JSON(200, gin.H{
		"status": "berhasil delete",
		"data":   article,
	})
}
