package main

import (
	"errors"

	log "github.com/sirupsen/logrus"
	mgo "gopkg.in/mgo.v2"
)

type article struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

var articleList = []article{
	article{ID: 1, Title: "Article 1", Content: "Article 1 body"},
	article{ID: 2, Title: "Article 2", Content: "Article 2 body"},
	article{ID: 3, Title: "Article 3", Content: "Article 3 body"},
	article{ID: 4, Title: "Article 4", Content: "Article 4 body"},
	article{ID: 5, Title: "Article 5", Content: "Article 5 body"},
	article{ID: 6, Title: "Article 6", Content: "Article 6 body"},
	article{ID: 7, Title: "Article 7", Content: "Article 7 body"},
	article{ID: 8, Title: "Article 8", Content: "Article 8 body"},
	article{ID: 9, Title: "Article 9", Content: "Article 9 body"},
	article{ID: 10, Title: "Article 10", Content: "Article 10 body"},
	article{ID: 11, Title: "Article 11", Content: "Article 11 body"},
	article{ID: 12, Title: "Article 12", Content: "Article 12 body"},
	article{ID: 13, Title: "Article 13", Content: "Article 13 body"},
	article{ID: 14, Title: "Article 14", Content: "Article 14 body"},
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

	/* session.DB("test").C("articles").RemoveAll(nil) */

	/* err = c.Insert(&article{ID: 15, Title: "Article 15", Content: "Article 15 body"},
		&article{ID: 16, Title: "Article 16", Content: "Article 16 body"},
		&article{ID: 17, Title: "Article 17", Content: "Article 17 body"},
		&article{ID: 18, Title: "Article 18", Content: "Article 18 body"},
		&article{ID: 19, Title: "Article 19", Content: "Article 19 body"})
	if err != nil {
		log.Info("Test5")
		log.Fatal(err)
	} */

	results := []article{}
	err = c.Find(nil).All(&results)
	if err != nil {
		log.Fatal(err)
	}

	return results
}

func getArticleByID(id int) (*article, error) {
	for _, a := range articleList {
		if a.ID == id {
			return &a, nil
		}
	}
	return nil, errors.New("Article not found")
}
