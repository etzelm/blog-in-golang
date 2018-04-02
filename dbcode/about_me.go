package main

import (
	"log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func about_me() {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("articles")

	blurb := "Just a quick blurb about me and this blog"
	created := "March 31st, 2018"
	modified := "April 1st, 2018"
	hold := "<div class=\"container\">" +
		"<div class=\"row\">" +
		"<div class=\"col-md-6 form-group\">" +
		"<br><img src=\"/public/face.jpg\" alt=\"My Face\">" +
		"</div>" +
		"<div class=\"col-md-6 form-group\">" +
		"<h3 style=\"color:#A619FF;\">Education:</h3>" +
		"<h4>&emsp;&emsp;Bachelors of Science in Computer Science</h4>" +
		"<h4>&emsp;&emsp;University of California, Santa Cruz</h4>" +
		"<h4>&emsp;&emsp;Attended from June 2015 to March 2018</h4><br>" +
		"<h4>&emsp;&emsp;Associates of Science in Computer Science</h4>" +
		"<h4>&emsp;&emsp;Diablo Valley College, Pleasant Hill, CA</h4>" +
		"<h4>&emsp;&emsp;Attended from January 2013 to June 2015</h4>" +
		"</div>" +
		"</div>" +
		"<h3 style=\"color:#A619FF;\">About Me:</h3>" +
		"<h4>&emsp;&emsp;Hopefully you've guessed by now that my name is Mitchell " +
		"Etzel and as a recent college graduate I've decided to start this blog to " +
		"help foster and share my further eduction. I take these actions in the " +
		"hope that they might help others along a similar journey one day. The " +
		"goal is to dedicate this website to distributed systems and things " +
		"related to them as well. Although I make promises toward the fact that " +
		"I will probably get a little bit distracted with side projects along the " +
		"way, here and there.</h4>" +
		"<h3 style=\"color:#A619FF;\">Contact Me:</h3>" +
		"<h4>&emsp;&emsp;<a style=\"color:#CC8A14;\" href=\"mailto:etzelm@live.com\">Through Email</a></h4>" +
		"<h4>&emsp;&emsp;<a style=\"color:#9C6708;\" href=\"https://github.com/etzelm\">Through Github</a></h4>" +
		"</div>"

	/* err = c.Insert(&article{ID: 0, Title: "About Me",
	Created: created, Modified: modified, Blurb: blurb, Content: hold}) */

	err = c.Update(bson.M{"id": 0},
		&article{ID: 0, Title: "About Me",
			Created: created, Modified: modified, Blurb: blurb, Content: hold})

	if err != nil {
		log.Fatal(err)
	}
}
