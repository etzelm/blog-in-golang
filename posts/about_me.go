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
	modified := "April 2nd, 2018"
	//		Start of the green well that backgrounds About Me post
	hold := "<div class=\"well\" style=\"background-color:#DFF0D8;\">" +
		//Start of the container for the face picture and education panel
		"<div class=\"container\">" +
		//Only used one row for the container
		"<div class=\"row\">" +
		//Start of the column with face picture in it
		"<div class=\"col-md-5 form-group\">" +
		//Spacing for face picture
		"<br><br>&emsp;&emsp;&emsp;&emsp;" +
		//Actual face picture
		"<img src=\"/public/face.jpg\" alt=\"My Face\" height=\"250\" width=\"250\"></div>" +
		//Start of the column with education panel in it
		"<div class=\"col-md-6 form-group\"><br>" +
		//Start of the education panel
		"<div class=\"panel panel-default\">" +
		"<div class=\"panel-heading\" style=\"color:#A619FF;font-size: 1.6em\">Education:</div>" +
		//Start of the list
		"<ul class=\"list-group\">" +
		"<li class=\"list-group-item\"><h4>&emsp;<b>&emsp;Bachelors of Science in Computer Science</b></h4>" +
		"<h4>&emsp;&emsp;University of California, Santa Cruz</h4>" +
		"<h4>&emsp;&emsp;Attended from June 2015 to March 2018</h4></li>" +
		"<li class=\"list-group-item\"><h4>&emsp;&emsp;<b>Associates of Science in Computer Science</b></h4>" +
		"<h4>&emsp;&emsp;Diablo Valley College, Pleasant Hill, CA</h4>" +
		"<h4>&emsp;&emsp;Attended from January 2013 to June 2015</h4></li>" +
		//End of the container for the face picture and education panel
		"</div></div></div></div>" +
		//Start of the About Me paragraph and contact info
		"<div class=\"panel panel-default\">" +
		"<div class=\"panel-heading\" style=\"color:#A619FF;font-size: 1.6em\">About Me & Contact Info:</div>" +
		"<ul class=\"list-group\">" +
		"<li class=\"list-group-item\">" +
		"<h4>&emsp;&emsp;Hopefully by now you've guessed that my name is Mitchell " +
		"Etzel and as a recent college graduate I've decided to start this blog to " +
		"help foster and share my further eduction. I take these actions in the " +
		"hope that they might help others along a similar journey one day. The " +
		"goal is to dedicate this website to distributed systems and things " +
		"related to them as well. Although I make promises toward the fact that " +
		"I will probably get a little distracted with side projects along the " +
		"way.</h4></li>" +
		"<li class=\"list-group-item\">" +
		"<h4>&emsp;&emsp;<a style=\"color:#CC8A14;\" href=\"mailto:etzelm@live.com\">" +
		"<img src=\"/public/email.png\" alt=\"Email\"  height=\"30\" width=\"45\"></a>" +
		"&emsp;&emsp;<a style=\"color:#9C6708;\" href=\"https://github.com/etzelm\">" +
		"<img src=\"/public/github.png\" alt=\"Github\"  height=\"45\" width=\"45\"></a></h4></li>" +
		"</div>" +
		"</div>" +
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
