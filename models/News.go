package models

import "gopkg.in/mgo.v2/bson"

// Struktur dari news
// bson digunakan untuk memberi tau mgo driver
type News struct {
	ID         bson.ObjectId
	ID_News    string `bson:"ID_News" json:"ID_News"`
	Website    string `bson:"Website" json:"Website"`
	News_Title string `bson:"News_Title" json:"News_Title"`
	Writer     string `bson:"Writer" json:"Writer"`
	Category   int    `bson:"Category" json:"Category"`
}

type MNews []News
