package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"

	// importing all packages

	config "MNews/config"
	"MNews/dao"
	models "MNews/models"

	// import json
	"encoding/json"
)

var conf = config.Config{}
var mdao = dao.MNewsDAO{}

// GET list of all MNews
func FindAllNews(w http.ResponseWriter, r *http.Request) {
	MNews, err := mdao.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, MNews)
}

// POST a new news
func Createnews(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var news models.News
	if err := json.NewDecoder(r.Body).Decode(&news); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	news.ID = bson.NewObjectId()
	if err := mdao.Insert(news); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, news)
}

// GET a news by its ID_news
func FindNewsByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	news, err := mdao.FindByID(params["ID_news"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Berita tidak ditemukan!")
		return
	}
	respondWithJson(w, http.StatusOK, news)
}

// GET a news by its category

func FindNewsByCategory(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	news, err := mdao.FindByCategory(params["Category"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Berita tidak ditemukan atau Kategori Salah")
		return
	}
	respondWithJson(w, http.StatusOK, news)
}

// PUT update an existing news
func Updatenews(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var news models.News
	if err := json.NewDecoder(r.Body).Decode(&news); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid")
		return
	}
	if err := mdao.Update(news); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

// DELETE an existing news
func Deletenews(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var news models.News
	if err := json.NewDecoder(r.Body).Decode(&news); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid")
		return
	}
	if err := mdao.Delete(news); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

// Error handling
func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

// JSON responses
func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Parsing connection into config.toml
// Connecting config into db
func init() {
	conf.Read()

	mdao.Server = conf.Server
	mdao.Database = conf.Database
	mdao.Connect()
}

// Main
// Handling HTTP Routes
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/MNews", FindAllNews).Methods("GET")
	r.HandleFunc("/MNews", Createnews).Methods("POST")
	r.HandleFunc("/MNews", Updatenews).Methods("PUT")
	r.HandleFunc("/MNews", Deletenews).Methods("DELETE")
	r.HandleFunc("/MNews/{ID_news}", FindNewsByID).Methods("GET")
	r.HandleFunc("/MNews/Category/{Category}", FindNewsByCategory).Methods("GET")
	// Start server
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
