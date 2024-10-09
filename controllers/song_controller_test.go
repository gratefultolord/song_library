package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"song_library/models"
)

// Функция для настройки тестовой базы данных
func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("Не удалось подключиться к тестовой базе данных")
	}
	db.AutoMigrate(&models.Song{})
	return db
}

// Тест для метода GetSongs
func TestGetSongs(t *testing.T) {
	db := setupTestDB()
	// Добавляем тестовую запись
	db.Create(&models.Song{Group: "TestGroup", Song: "TestSong"})

	req, _ := http.NewRequest("GET", "/songs", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/songs", func(w http.ResponseWriter, r *http.Request) {
		GetSongs(w, r, db)
	}).Methods("GET")
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Ожидаемый статус код %d, получено %d", http.StatusOK, rr.Code)
	}

	var songs []models.Song
	if err := json.Unmarshal(rr.Body.Bytes(), &songs); err != nil {
		t.Errorf("Не удалось декодировать ответ: %v", err)
	}

	if len(songs) == 0 {
		t.Errorf("Ожидался хотя бы один результат")
	}
}

// Тест для метода AddSong
func TestAddSong(t *testing.T) {
	db := setupTestDB()

	newSong := models.Song{
		Group: "NewGroup",
		Song:  "NewSong",
	}
	body, _ := json.Marshal(newSong)

	req, _ := http.NewRequest("POST", "/song", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/song", func(w http.ResponseWriter, r *http.Request) {
		AddSong(w, r, db)
	}).Methods("POST")
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Ожидаемый статус код %d, получено %d", http.StatusCreated, rr.Code)
	}

	var song models.Song
	db.First(&song)

	if song.Group != newSong.Group || song.Song != newSong.Song {
		t.Errorf("Добавленная песня не соответствует ожиданиям")
	}
}

// Тест для метода GetSong
func TestGetSong(t *testing.T) {
	db := setupTestDB()
	song := models.Song{Group: "TestGroup", Song: "TestSong"}
	db.Create(&song)

	req, _ := http.NewRequest("GET", "/song/1", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/song/{id}", func(w http.ResponseWriter, r *http.Request) {
		GetSong(w, r, db)
	}).Methods("GET")
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Ожидаемый статус код %d, получено %d", http.StatusOK, rr.Code)
	}

	var fetchedSong models.Song
	if err := json.Unmarshal(rr.Body.Bytes(), &fetchedSong); err != nil {
		t.Errorf("Не удалось декодировать ответ: %v", err)
	}

	if fetchedSong.Group != song.Group || fetchedSong.Song != song.Song {
		t.Errorf("Полученная песня не соответствует ожиданиям")
	}
}

// Тест для метода UpdateSong
func TestUpdateSong(t *testing.T) {
	db := setupTestDB()
	song := models.Song{Group: "TestGroup", Song: "TestSong"}
	db.Create(&song)

	updatedSong := models.Song{Group: "UpdatedGroup", Song: "UpdatedSong"}
	body, _ := json.Marshal(updatedSong)

	req, _ := http.NewRequest("PUT", "/song/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/song/{id}", func(w http.ResponseWriter, r *http.Request) {
		UpdateSong(w, r, db)
	}).Methods("PUT")
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Ожидаемый статус код %d, получено %d", http.StatusOK, rr.Code)
	}

	var updated models.Song
	db.First(&updated, song.ID)

	if updated.Group != updatedSong.Group || updated.Song != updatedSong.Song {
		t.Errorf("Обновленная песня не соответствует ожиданиям")
	}
}

// Тест для метода DeleteSong
func TestDeleteSong(t *testing.T) {
	db := setupTestDB()
	song := models.Song{Group: "TestGroup", Song: "TestSong"}
	db.Create(&song)

	req, _ := http.NewRequest("DELETE", "/song/1", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/song/{id}", func(w http.ResponseWriter, r *http.Request) {
		DeleteSong(w, r, db)
	}).Methods("DELETE")
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("Ожидаемый статус код %d, получено %d", http.StatusNoContent, rr.Code)
	}

	var deleted models.Song
	result := db.First(&deleted, song.ID)
	if result.Error == nil {
		t.Errorf("Песня не была удалена")
	}
}
