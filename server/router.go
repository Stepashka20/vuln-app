package server

import (
	"VulnApp/controllers"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func NewRouter() *gin.Engine {
	router := gin.Default()

	store := cookie.NewStore([]byte(os.Getenv("COOKIE_SECRET")))
	router.Use(sessions.Sessions("mysession", store))

	router.LoadHTMLGlob("./templates/*.html")

	auth := new(controllers.AuthController)
	router.POST("/login", auth.Login)
	router.POST("/register", auth.Register)
	router.GET("/logout", auth.Logout)
	router.GET("/", auth.Index)

	user := new(controllers.UserController)
	router.GET("/profile", user.Profile)
	router.GET("/delete", user.Delete)
	router.POST("/upload", user.UploadFile)

	misc := new(controllers.MiscController)
	router.POST("/ping", misc.Ping)
	router.GET("/files", misc.Files)

	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})
	router.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", nil)
	})
	router.GET("/ping", func(c *gin.Context) {
		c.HTML(http.StatusOK, "ping.html", nil)
	})

	return router

}
