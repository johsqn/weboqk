package main

import (
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/johsqn/weboqk/internal/repository"
	"github.com/johsqn/weboqk/internal/repository/dao"
	"github.com/johsqn/weboqk/internal/service"
	"github.com/johsqn/weboqk/internal/web"
	"github.com/johsqn/weboqk/internal/web/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	db := initDB()
	server := initWebServer()

	initUserHdl(db, server)

	server.Run(":8080")
}

func initUserHdl(db *gorm.DB, server *gin.Engine) {

	userDao := dao.NewUserDAO(db)
	userPO := repository.NewUserRepository(userDao)
	userServ := service.NewUserService(userPO)
	userHdl := web.NewUserHandler(userServ)
	userHdl.RegisterRoutes(server)
}

func initDB() *gorm.DB {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		panic("MYSQL_DSN is not set")
	}
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(err)
	}

	err = dao.InitUserTables(db)
	if err != nil {
		panic(err)
	}
	return db
}

func initWebServer() *gin.Engine {

	server := gin.Default()

	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000", "https://your_company.com"},
		AllowCredentials: true,

		AllowHeaders: []string{"Content-Type", "Authorization"},
		//AllowHeaders: []string{"content-type"},
		//AllowMethods: []string{"POST"},
		MaxAge: 12 * time.Hour,
	}))

	// login := &middleware.LoginBuilder{}
	// 存储数据的，也就是你 userId 存哪里
	// 直接存 cookie
	store := cookie.NewStore([]byte("check-login-secret"))
	server.Use(sessions.Sessions("ssid", store), middleware.CheckLogin())
	return server

}
