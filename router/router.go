package router

import (
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
	"song_library/controllers"
)

// SetupRouter настраивает маршруты для обработки запросов
func SetupRouter(db *gorm.DB) *mux.Router {
	r := mux.NewRouter()

	// Маршрут для получения списка песен
	// @Summary Получить список песен
	// @Description Получение песен с фильтрацией и пагинацией
	// @Tags songs
	// @Accept json
	// @Produce json
	// @Param page query int false "Номер страницы"
	// @Param limit query int false "Количество элементов на странице"
	// @Param group query string false "Название группы"
	// @Param song query string false "Название песни"
	// @Param releaseDate query string false "Дата выпуска песни"
	// @Success 200 {array} models.Song
	// @Failure 500 {object} map[string]string
	// @Router /songs [get]
	r.HandleFunc("/songs", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetSongs(w, r, db)
	}).Methods("GET")

	// Маршрут для получения информации о конкретной песне
	// @Summary Получить информацию о песне
	// @Description Получение информации о конкретной песне по ID
	// @Tags songs
	// @Accept json
	// @Produce json
	// @Param id path int true "ID песни"
	// @Success 200 {object} models.Song
	// @Failure 404 {object} map[string]string
	// @Router /song/{id} [get]
	r.HandleFunc("/song/{id}", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetSong(w, r, db)
	}).Methods("GET")

	// Маршрут для добавления новой песни
	// @Summary Добавить новую песню
	// @Description Добавление новой песни в библиотеку
	// @Tags songs
	// @Accept json
	// @Produce json
	// @Param song body models.Song true "Добавить песню"
	// @Success 201 {object} models.Song
	// @Failure 400 {object} map[string]string
	// @Failure 500 {object} map[string]string
	// @Router /song [post]
	r.HandleFunc("/song", func(w http.ResponseWriter, r *http.Request) {
		controllers.AddSong(w, r, db)
	}).Methods("POST")

	// Маршрут для обновления данных о песне
	// @Summary Обновить информацию о песне
	// @Description Обновление данных о существующей песне по ID
	// @Tags songs
	// @Accept json
	// @Produce json
	// @Param id path int true "ID песни"
	// @Param song body models.Song true "Обновить песню"
	// @Success 200 {object} models.Song
	// @Failure 400 {object} map[string]string
	// @Failure 404 {object} map[string]string
	// @Failure 500 {object} map[string]string
	// @Router /song/{id} [put]
	r.HandleFunc("/song/{id}", func(w http.ResponseWriter, r *http.Request) {
		controllers.UpdateSong(w, r, db)
	}).Methods("PUT")

	// Маршрут для удаления песни
	// @Summary Удалить песню
	// @Description Удаление песни из библиотеки по ID
	// @Tags songs
	// @Accept json
	// @Produce json
	// @Param id path int true "ID песни"
	// @Success 204
	// @Failure 500 {object} map[string]string
	// @Router /song/{id} [delete]
	r.HandleFunc("/song/{id}", func(w http.ResponseWriter, r *http.Request) {
		controllers.DeleteSong(w, r, db)
	}).Methods("DELETE")

	return r
}
