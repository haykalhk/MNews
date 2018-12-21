package dao

import (
	models "MNews/models"
	"log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Starting establish connection wth server

type MNewsDAO struct {
	Server   string
	Database string
}

// For MongoDB Database
var db *mgo.Database

// Collect News
const (
	COLLECTION = "MNews"
)

// Establish connections to MongoDB Database
func (c *MNewsDAO) Connect() {
	session, err := mgo.Dial(c.Server)
	if err != nil {
		log.Fatal(err)
	}
	//create session
	db = session.DB(c.Database)
}

// Query-ing database

// Finding all of News
func (c *MNewsDAO) FindAll() ([]models.News, error) {
	var MNews []models.News
	err := db.C(COLLECTION).Find(bson.M{}).All(&MNews)
	return MNews, err
}

// Finding a News by ID
func (c *MNewsDAO) FindByID(ID_News string) (models.News, error) {
	var News models.News
	err := db.C(COLLECTION).Find(bson.M{"ID_News": ID_News}).One(&News)
	return News, err
}

//  CRUD for News DB
// Insert News into database
func (c *MNewsDAO) Insert(News models.News) error {
	err := db.C(COLLECTION).Insert(&News)
	return err
}

// Delete an existing News
func (c *MNewsDAO) Delete(News models.News) error {
	err := db.C(COLLECTION).Remove(&News)
	return err
}

// Update an existing News
func (c *MNewsDAO) Update(News models.News) error {
	err := db.C(COLLECTION).Update(bson.M{"ID_News": News.ID_News}, &News)
	return err
}
