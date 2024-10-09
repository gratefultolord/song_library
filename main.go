// Online Song Library Service

// Описание проекта:
// - REST методы для онлайн библиотеки песен
// - Интеграция с внешним API для обогащения данных песен
// - PostgreSQL для хранения данных с миграциями
// - Документация Swagger для API

package main

import (
	"log"
	"net/http"
	"os"
	"song_library/database"
	_ "song_library/docs" // импорт сгенерированной документации
	"song_library/models"
	"song_library/router"
	"song_library/utils"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title API для онлайн библиотеки песен
// @version 1.0
// @description Это пример сервера для управления библиотекой песен.
// @termsOfService http://swagger.io/terms/

// @contact.name Поддержка API
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func main() {
	utils.LoadEnv()
	log.Println("[INFO] Загружены переменные окружения")

	dsn := os.Getenv("POSTGRES_CONN")
	log.Printf("[DEBUG] Подключение к базе данных с DSN: %s\n", dsn)
	db, err := database.ConnectDatabase(dsn)
	if err != nil {
		log.Fatalf("[ERROR] Не удалось подключиться к базе данных: %v", err)
	}
	log.Println("[INFO] Успешное подключение к базе данных")

	if err := db.AutoMigrate(&models.Song{}); err != nil {
		log.Fatalf("[ERROR] Не удалось выполнить миграции: %v", err)
	}
	log.Println("[INFO] Миграции выполнены успешно")

	r := router.SetupRouter(db)
	log.Println("[INFO] Маршрутизатор настроен")

	// Добавление маршрута для Swagger UI
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	log.Println("[INFO] Маршрут для Swagger UI добавлен")

	log.Println("[INFO] Сервер запускается...")
	log.Fatal(http.ListenAndServe(os.Getenv("SERVER_ADDRESS"), r))
}
