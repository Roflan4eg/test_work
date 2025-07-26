package main

import (
	"database/sql"
	"github.com/Roflan4eg/test_work/internal/api"
	"github.com/Roflan4eg/test_work/internal/storage"
	"log"
	"os"

	_ "github.com/Roflan4eg/test_work/docs"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	conn, err := storage.New(db)
	if err != nil {
		log.Fatal(err)
		return
	}
	handler := api.New(conn)

	r := gin.Default()
	r.Use(api.ErrorHandler())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes := r.Group("/api/v1/subscriptions")
	{
		routes.POST("", handler.CreateSubscription)
		routes.POST("/get_for_period", handler.GetPriceForPeriod)
		routes.GET("", handler.ListSubscriptions)
		routes.GET("/:id", handler.GetSubscription)
		routes.PUT("/:id", handler.UpdateSubscription)
		routes.DELETE("/:id", handler.DeleteSubscription)
	}

	if err = r.Run(":" + os.Getenv("PORT")); err != nil {
		log.Fatal(err)
		return
	}
}
