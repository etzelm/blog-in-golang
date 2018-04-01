package main

import (
	"errors"

	log "github.com/sirupsen/logrus"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type article struct {
	ID       int    `json:"id"`
	Created  string `json:"created"`
	Modified string `json:"modified"`
	Title    string `json:"title"`
	Blurb    string `json:"blurb"`
	Content  string `json:"content"`
}

// Return a list of all the articles
func getAllArticles() []article {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("articles")

	results := []article{}
	err = c.Find(nil).Sort("-id").All(&results)
	if err != nil {
		log.Fatal(err)
	}
	return results
}

func getArticleByID(id int) (*article, error) {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("articles")

	result := article{}
	err = c.Find(bson.M{"id": id}).One(&result)
	if err != nil {
		log.Fatal(err)
		return nil, errors.New("Article not found")
	}
	return &result, nil
}
