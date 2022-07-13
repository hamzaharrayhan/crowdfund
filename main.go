package main

import (
	"crowdfund/handler"
	"crowdfund/user"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "root:@tcp(127.0.0.1:3306)/crowdfund?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	//bikin router
	router := gin.Default()

	//bikin grup routing dengan prefix yang sama
	api := router.Group("/api/v1")

	//bikin route ke user dengan prefix /api/v1
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.LoginHandler)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	router.Run()
}
