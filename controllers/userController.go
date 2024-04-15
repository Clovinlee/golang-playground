package controllers

import (
	"chris/gochris/initializers"
	"chris/gochris/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func UserRegister(c *gin.Context) {

	var body struct {
		Email    string
		Password string
		Name     string
	}

	c.Bind(&body)

	if body.Email == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email and Password are required",
		})
		return
	}

	var user models.User
	resultEmail := initializers.DB.Where("email = ?", body.Email).First(&user)
	if resultEmail.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email already exists",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password hashing failed",
		})
		return
	}

	user = models.User{Email: body.Email, Password: string(hashedPassword), Name: body.Name}

	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name": user.Name,
	})
}

func UserLogin(c *gin.Context) {
	var body struct {
		Email    string
		Password string
	}

	c.Bind(&body)

	if body.Email == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email and Password are required",
		})
		return
	}

	var user models.User

	result := initializers.DB.Where("email = ?", body.Email).First(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Login Error",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil || result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Email / Password",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error generating token",
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*7, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{})

}
