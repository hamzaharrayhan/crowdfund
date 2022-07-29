package main

import (
	"crowdfund/auth"
	"crowdfund/campaign"
	"crowdfund/handler"
	"crowdfund/helper"
	"crowdfund/transaction"
	"crowdfund/user"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
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

	// user
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	authService := auth.NewJwtService()

	userHandler := handler.NewUserHandler(userService, authService)

	//campaign
	campaignRepository := campaign.NewRepository(db)
	campaignService := campaign.NewService(campaignRepository)
	campaignHandler := handler.NewCampaignHandler(campaignService)

	// transactions
	transactionRepository := transaction.NewRepository(db)
	transactionService := transaction.NewService(transactionRepository, campaignRepository)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	//bikin router
	router := gin.Default()

	//bikin grup routing dengan prefix yang sama
	api := router.Group("/api/v1")
	router.Static("/images", "./images")
	//bikin route ke user dengan prefix /api/v1
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.LoginHandler)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadUserAvatar)

	//route ke campaigns
	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaigns/:id", campaignHandler.GetCampaign)
	api.POST("/campaigns", authMiddleware(authService, userService), campaignHandler.CreateCampaign)
	api.PUT("/campaigns/:id", authMiddleware(authService, userService), campaignHandler.UpdateCampaign)
	api.POST("/campaign-images", authMiddleware(authService, userService), campaignHandler.SaveCampaignImage)

	// route ke campaign transactions
	api.GET("/campaigns/:id/transactions", authMiddleware(authService, userService), transactionHandler.GetCampaignTransactions)
	api.GET("/transactions", authMiddleware(authService, userService), transactionHandler.GetUserTransactions)

	router.Run()
}

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.JSONResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		arrayAuthHeader := strings.Split(authHeader, " ")
		token := ""
		if len(arrayAuthHeader) == 2 {
			token = arrayAuthHeader[1]
		}

		validatedToken, err := authService.ValidateToken(token)
		if err != nil {
			response := helper.JSONResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claims, ok := validatedToken.Claims.(jwt.MapClaims)

		if !ok || !validatedToken.Valid {
			response := helper.JSONResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claims["user_id"].(float64))
		user, err := userService.GetUserByID(userID)

		if err != nil {
			response := helper.JSONResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)
	}
}
