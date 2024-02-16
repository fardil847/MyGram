package controllers

import (
	"MyGram/helpers"
	"MyGram/models"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SocialMediaRepository struct {
	DB *gorm.DB
}

func (m *SocialMediaRepository) UploadSocialMedia(c *gin.Context) {
	contentType := helpers.GetContentType(c)
	SocialMedia := models.SocialMedia{}

	userData := c.MustGet("userData").(jwt.MapClaims)
	userId := userData["id"].(float64)

	if contentType == "application/json" {
		c.ShouldBindJSON(&SocialMedia)
	} else {
		c.ShouldBind(&SocialMedia)
	}

	SocialMedia.User_id = uint(userId)
	SocialMedia.Created_at = time.Now()
	SocialMedia.Updated_at = time.Now()

	if err := m.DB.Debug().Create(&SocialMedia).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to upload social media",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":               SocialMedia.Id,
		"name":             SocialMedia.Name,
		"social_media_url": SocialMedia.Social_media_url,
		"user_id":          SocialMedia.User_id,
		"created_at":       SocialMedia.Created_at,
	})
}

func (m *SocialMediaRepository) GetSocialMedia(c *gin.Context) {
	SocialMedias := []models.SocialMedia{}

	if err := m.DB.Debug().Find(&SocialMedias).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "can't find media",
			"message": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"socialmedia": SocialMedias,
	})

}

func (m *SocialMediaRepository) UpdateSocialMedia(c *gin.Context) {
	contentType := helpers.GetContentType(c)

	userData := c.MustGet("userData").(jwt.MapClaims)
	userId := userData["id"].(float64)

	getId, _ := strconv.Atoi(c.Param("socialMediaId"))

	SocialMedia := models.SocialMedia{}
	OldSocialMedia := models.SocialMedia{}

	if contentType == "application/json" {
		c.ShouldBindJSON(&SocialMedia)
	} else {
		c.ShouldBind(&SocialMedia)
	}

	SocialMedia.Updated_at = time.Now()
	SocialMedia.User_id = uint(userId)

	if err := m.DB.Debug().First(&OldSocialMedia, getId).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "media not found",
			"message": err.Error(),
		})
		return
	}

	if err := m.DB.Debug().Model(&OldSocialMedia).Updates(&SocialMedia).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to update media",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":               OldSocialMedia.Id,
		"name":             OldSocialMedia.Name,
		"social_media_url": OldSocialMedia.Social_media_url,
		"user_id":          OldSocialMedia.User_id,
		"updated_at":       OldSocialMedia.Updated_at,
	})

}

func (m *SocialMediaRepository) DeleteSocialMedia(c *gin.Context) {
	getId, _ := strconv.Atoi(c.Param("socialMediaId"))
	SocialMedia := models.SocialMedia{}

	if err := m.DB.Debug().First(&SocialMedia, getId).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "media not found",
			"message": err.Error(),
		})
		return
	}

	if err := m.DB.Debug().Delete(&SocialMedia).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to delete media",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "your social media has been successfully deleted",
	})

}
