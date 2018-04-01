package main

import (
	mgo "gopkg.in/mgo.v2"
)

func drop_tables() {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.DB("test").C("articles").RemoveAll(nil)
}
