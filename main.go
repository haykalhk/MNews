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
var cobadao = dao.MNewsDAO{}

// POST a new news
// Decodes the request body into a news object, assign it an ID_MNews, and uses the DAO Insert method to create a news in database
func Createnews(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var news models.News
	if err := json.NewDecoder(r.Body).Decode(&news); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	news.ID = bson.NewObjectId()
	if err := cobadao.Insert(news); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, news)
}

// GET list of all MNews
// Uses FindAll method of DAO Library to fetch list of MNews from database
func AllMNews(w http.ResponseWriter, r *http.Request) {
	MNews, err := cobadao.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, MNews)
}

// CRUD

// GET a news by its ID_news
// Using mux library to get parameters that the users passed in with the request
func Findnews(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	news, err := cobadao.FindByID(params["ID_news"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "news tidak ditemukan!")
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
	if err := cobadao.Update(news); err != nil {
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
	if err := cobadao.Delete(news); err != nil {
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

	cobadao.Server = conf.Server
	cobadao.Database = conf.Database
	cobadao.Connect()
}

// Main
// Handling HTTP Routes
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/MNews", AllMNews).Methods("GET")
	r.HandleFunc("/MNews", Createnews).Methods("POST")
	r.HandleFunc("/MNews", Updatenews).Methods("PUT")
	r.HandleFunc("/MNews", Deletenews).Methods("DELETE")
	r.HandleFunc("/MNews/{ID_news}", Findnews).Methods("GET")
	// Start server
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
