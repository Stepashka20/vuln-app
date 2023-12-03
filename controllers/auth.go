package controllers

import (
	"VulnApp/db"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type AuthController struct{}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (a *AuthController) Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	var storedUsername, storedPassword string

	err := db.GetDB().QueryRow("SELECT username, password FROM users WHERE username = '"+username+"'").Scan(&storedUsername, &storedPassword)
	if err != nil {
		fmt.Println(err)
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"message": "Неправильное имя пользователя или пароль"})
		return
	}

	if CheckPasswordHash(password, storedPassword) {
		session := sessions.Default(c)
		session.Set("username", storedUsername)
		session.Save()

		c.Redirect(http.StatusFound, "/profile")
		return
	}

	c.HTML(http.StatusBadRequest, "error.html", gin.H{"message": "Неправильное имя пользователя или пароль"})
}

func (a *AuthController) Register(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username != "" && password != "" {
		var userExists bool
		err := db.GetDB().QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", username).Scan(&userExists)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{"message": "Ошибка базы данных"})
			return
		}
		if userExists {
			c.HTML(http.StatusBadRequest, "error.html", gin.H{"message": "Пользователь уже существует"})
			return
		}
		hashedPassword, _ := HashPassword(password)
		_, err = db.GetDB().Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, hashedPassword)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{"message": "Ошибка базы данных"})
			return
		}

		session := sessions.Default(c)
		session.Set("username", username)
		session.Save()

		c.HTML(http.StatusOK, "register_success.html", gin.H{"username": username})
		return
	}

	c.HTML(http.StatusBadRequest, "error.html", gin.H{"message": "Введите имя пользователя и пароль"})
}

func (a *AuthController) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusFound, "/")
}

func (a *AuthController) Index(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")
	if username != nil {
		c.Redirect(http.StatusFound, "/profile")
		return
	}
	c.HTML(http.StatusOK, "index.html", nil)
}
