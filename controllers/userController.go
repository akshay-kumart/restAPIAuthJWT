package controllers

import (
	"github.com/akshay-kumart/go-api/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var Users []models.User

func SignUp(c *gin.Context) {
	var Body struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := c.BindJSON(&Body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to Bind"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(Body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to Generate Hash"})
	}

	Users = append(Users, models.User{ID: rand.Int(), Username: Body.Username, Password: string(hash), Role: Body.Role})

	c.JSON(http.StatusOK, gin.H{"message": Users})

}

func Login(c *gin.Context) {
	var Body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&Body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to Bind"})
		return
	}

	id := 0
	role := ""
	for _, val := range Users {
		if val.Username == Body.Username {
			if err := bcrypt.CompareHashAndPassword([]byte(val.Password), []byte(Body.Password)); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": "Wrong Password"})
				return
			}
			id = val.ID
			role = val.Role
			break
		}
	}
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":    id,
		"expire": time.Now().Add(time.Hour * 24 * 30).Unix(),
		"role":   role,
	})

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Server misconfigured"})
		return
	}

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		log.Println("Error signing token:", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create token"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "login successfull"})

}

func Validate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Imm insidee yaay"})
}

func Role(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Imm the userr hehe"})
}
