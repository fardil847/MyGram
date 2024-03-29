package controllers

import (
	"MyGram/helpers"
	"MyGram/models"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func (u *UserRepository) UserRegister(c *gin.Context) {
	contentType := helpers.GetContentType(c)
	_, _ = u.DB, contentType

	User := models.User{}

	if contentType == "application/json" {
		c.ShouldBindJSON(&User)
	} else {
		c.ShouldBind(&User)
	}

	User.Created_at = time.Now()
	User.Updated_at = time.Now()

	if err := u.DB.Debug().Create(&User).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to create user data",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"id":    User.GormModel.Id,
		"email": User.Email,
		"name":  User.Name,
	})
}

func (u *UserRepository) GetAllUser(c *gin.Context) {
	contentType := helpers.GetContentType(c)
	_, _ = u.DB, contentType

	AllUser := []models.User{}

	if err := u.DB.Preload("Photos").Preload("Comments").Preload("SocialMedias").Find(&AllUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "can't find data",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    AllUser,
	})
}

func (u *UserRepository) UserLogin(c *gin.Context) {
	contentType := helpers.GetContentType(c)
	_, _ = u.DB, contentType

	User := models.User{}

	if contentType == "application/json" {
		c.ShouldBindJSON(&User)
	} else {
		c.ShouldBind(&User)
	}
	password := User.Password

	if err := u.DB.Debug().Where("email=?", User.Email).Take(&User).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "Invalid email/password",
		})
		return
	}

	fmt.Println((User.Password), (password))
	if comparePass := helpers.ComparePass([]byte(User.Password), []byte(password)); !comparePass {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "Invalid email/password",
		})
		return
	}
	token := helpers.GenerateToken(uint(User.GormModel.Id), User.Email)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})

}

func (u *UserRepository) UserUpdate(c *gin.Context) {
	GetId, _ := strconv.Atoi(c.Param("userId"))
	UserData := c.MustGet("userData").(jwt.MapClaims)
	UserId := UserData["id"].(float64)

	contextType := helpers.GetContentType(c)
	_, _ = u.DB, contextType

	User := models.User{}
	OldUser := models.User{}

	if contextType == "application/json" {
		c.ShouldBindJSON(&User)
	} else {
		c.ShouldBind(&User)
	}

	User.Updated_at = time.Now()
	User.Id = int(UserId)

	if err := u.DB.Where("id=?", GetId).Take(&OldUser).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "data not found",
			"message": err.Error(),
		})
		return
	}
	if err := u.DB.Preload("Photos").Preload("Comments").Preload("Medias").Model(&OldUser).Updates(&User).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to update data",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":         OldUser.Id,
		"email":      OldUser.Email,
		"username":   OldUser.Name,
		"age":        OldUser.Age,
		"updated_at": OldUser.Updated_at,
	})
}

func (u *UserRepository) UserDelete(c *gin.Context) {
	UserData := c.MustGet("userData").(jwt.MapClaims)
	UserId := int(UserData["id"].(float64))
	User := models.User{}

	if err := u.DB.First(&User, UserId).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "data not found",
			"message": err.Error(),
		})
		return
	}
	if err := u.DB.Preload("Photos").Preload("Comments").Preload("Medias").Model(&User).Delete(&User).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to delete data",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "your account has been succesfully deleted",
	})
}
