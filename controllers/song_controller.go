package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"song_library/models"
	"song_library/utils"
)

// GetSongs godoc
// @Summary Получить список песен
// @Description Получение песен с пагинацией и фильтрацией
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
func GetSongs(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	log.Println("[INFO] Обработка запроса на получение списка песен")
	query := utils.ApplyPaginationAndFiltering(r, db)

	var songs []models.Song
	if err := query.Find(&songs).Error; err != nil {
		log.Printf("[ERROR] Не удалось получить список песен: %v\n", err)
		http.Error(w, "Не удалось получить список песен", http.StatusInternalServerError)
		return
	}

	log.Printf("[INFO] Найдено %d песен\n", len(songs))
	json.NewEncoder(w).Encode(songs)
}

// GetSong godoc
// @Summary Получить информацию о песне
// @Description Получение информации о конкретной песне по ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Success 200 {object} models.Song
// @Failure 404 {object} map[string]string
// @Router /song/{id} [get]
func GetSong(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	params := mux.Vars(r)
	id := params["id"]
	log.Printf("[INFO] Обработка запроса на получение песни с ID: %s\n", id)

	var song models.Song
	if err := db.First(&song, id).Error; err != nil {
		log.Printf("[ERROR] Песня с ID %s не найдена: %v\n", id, err)
		http.Error(w, "Песня не найдена", http.StatusNotFound)
		return
	}

	log.Printf("[INFO] Песня с ID %s успешно найдена\n", id)
	json.NewEncoder(w).Encode(song)
}

// AddSong godoc
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
func AddSong(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	log.Println("[INFO] Обработка запроса на добавление новой песни")
	var song models.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		log.Printf("[ERROR] Неверный формат запроса: %v\n", err)
		http.Error(w, "Неверный ввод", http.StatusBadRequest)
		return
	}

	log.Printf("[DEBUG] Добавляем песню: %+v\n", song)
	if err := db.Create(&song).Error; err != nil {
		log.Printf("[ERROR] Не удалось сохранить песню: %v\n", err)
		http.Error(w, "Не удалось сохранить песню", http.StatusInternalServerError)
		return
	}

	log.Println("[INFO] Песня успешно добавлена")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(song)
}

// UpdateSong godoc
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
func UpdateSong(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	params := mux.Vars(r)
	id := params["id"]
	log.Printf("[INFO] Обработка запроса на обновление песни с ID: %s\n", id)

	var song models.Song
	if err := db.First(&song, id).Error; err != nil {
		log.Printf("[ERROR] Песня с ID %s не найдена: %v\n", id, err)
		http.Error(w, "Песня не найдена", http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		log.Printf("[ERROR] Неверный формат запроса: %v\n", err)
		http.Error(w, "Неверный ввод", http.StatusBadRequest)
		return
	}

	log.Printf("[DEBUG] Обновляем песню с ID %s: %+v\n", id, song)
	if err := db.Save(&song).Error; err != nil {
		log.Printf("[ERROR] Не удалось обновить песню с ID %s: %v\n", id, err)
		http.Error(w, "Не удалось обновить песню", http.StatusInternalServerError)
		return
	}

	log.Printf("[INFO] Песня с ID %s успешно обновлена\n", id)
	json.NewEncoder(w).Encode(song)
}

// DeleteSong godoc
// @Summary Удалить песню
// @Description Удаление песни из библиотеки по ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Success 204
// @Failure 500 {object} map[string]string
// @Router /song/{id} [delete]
func DeleteSong(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	params := mux.Vars(r)
	id := params["id"]
	log.Printf("[INFO] Обработка запроса на удаление песни с ID: %s\n", id)

	if err := db.Delete(&models.Song{}, id).Error; err != nil {
		log.Printf("[ERROR] Не удалось удалить песню с ID %s: %v\n", id, err)
		http.Error(w, "Не удалось удалить песню", http.StatusInternalServerError)
		return
	}

	log.Printf("[INFO] Песня с ID %s успешно удалена\n", id)
	w.WriteHeader(http.StatusNoContent)
}
