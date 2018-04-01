package main

import (
	"log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func how_its_made() {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("articles")

	blurb := "We'll besides the general allure of writing blogs in and of them self, " +
		"Golang is actually pretty easy to write a blog or pretty much any kind of website in"
	created := "April 1st, 2018"
	modified := "April 1st, 2018"
	hold := "<h3>Who Would Want to Write a Blog in Go?</h3>" +
		"<h4>&emsp;&emsp;Hopefully you've guessed by now that my name is Mitchell " +
		"Etzel and as a recent college graduate I've decided to start this blog to " +
		"help foster and share my further eduction. I take these actions in the " +
		"hope that they might help others along a similar journey one day. The " +
		"goal is to dedicate this website to distributed systems and things " +
		"related to them as well. Although I make promises toward the fact that " +
		"I will probably get a little bit distracted with side projects along the " +
		"way, here and there.</h4>" +
		"</div>"

	/* err = c.Insert(&article{ID: 1, Title: "Who Would Want to Write a Blog in Go?",
	Created: created, Modified: modified, Blurb: blurb, Content: hold}) */

	err = c.Update(bson.M{"id": 1},
		&article{ID: 1, Title: "Who Would Want to Write a Blog in Go?",
			Created: created, Modified: modified, Blurb: blurb, Content: hold})

	if err != nil {
		log.Fatal(err)
	}
}
