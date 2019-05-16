package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	log "github.com/sirupsen/logrus"
)

func about_me() {
	id := os.Getenv("AWS_ACCESS_KEY_ID")
	key := os.Getenv("AWS_SECRET_ACCESS_KEY")
	log.Info("id: ", id)
	log.Info("key: ", key)
	var my_credentials = credentials.NewStaticCredentials(id, key, "")

	sess, err := session.NewSession(&aws.Config{
		Credentials: my_credentials,
		Region:      aws.String("us-west-1"),
		//Endpoint:    aws.String("http://localhost:8000"),
	})
	if err != nil {
		log.Println(err)
		return
	}
	dbSvc := dynamodb.New(sess)

	blurb := "Just a quick blurb about me and this blog"
	created := "March 31st, 2018"
	modified := "May 16th, 2019"
	//		Start of the green well that backgrounds About Me post
	hold := "<div class=\"well\" style=\"background-color:#DFF0D8;\">" +
		//Start of the container for the face picture and education panel
		"<div class=\"container-fluid\">" +
		//Only used one row for the container
		"<div class=\"row\">" +
		//Start of the column with face picture in it
		"<div class=\"col-md-5 form-group\" style=\"text-align: center;\">" +
		//Spacing for face picture
		"<br>" +
		//Actual face picture
		"<img src=\"/public/face.png\" alt=\"My Face\" height=\"250\" width=\"250\"></div>" +
		//Start of the column with education panel in it
		"<div class=\"col-md-6 form-group\"><br>" +
		//Start of the education panel
		"<div class=\"panel panel-default\">" +
		"<div class=\"panel-heading\" style=\"color:#A619FF;font-size: 1.17em\">Education:</div>" +
		//Start of the list
		"<ul class=\"list-group\" style=\"font-size: 1.07em;\">" +
		"<li class=\"list-group-item\"><b>Bachelors of Science in Computer Science</b><br>" +
		"UCSC, Santa Cruz, CA<br>" +
		"Attended: June '15 to March '18</li>" +
		"<li class=\"list-group-item\"><b>Associates of Science in Computer Science</b><br>" +
		"DVC, Pleasant Hill, CA<br>" +
		"Attended: Janurary '13 to June '15</li>" +
		//End of the container for the face picture and education panel
		"</div></div></div></div>" +
		//Start of the About Me paragraph and contact info
		"<div class=\"panel panel-default\">" +
		"<div class=\"panel-heading\" style=\"color:#A619FF;font-size: 1.17em\">About Me & Contact Info:</div>" +
		"<ul class=\"list-group\" style=\"font-size: 1.07em;\">" +
		"<li class=\"list-group-item\">" +
		"&emsp;&emsp;By now you've hopefully guessed that my name is Mitchell Etzel and as a college graduate, " +
		"as well as a recent member of the cyber-security industry, I've decided to maintain this blog to help foster and " +
		"share some of my experiences in the hope that others may find them useful. The goal is to dedicate this website " +
		"to the study of distributed systems and cyber-security, in addition to all topics I deem related along the way. " +
		"Distributed systems are all around us and in most of the web services we use nowadays. Considering that they have " +
		"only gotten more popular over the last twenty years, I believe they deserve a more thorough and discerning study. " +
		"That also means that as an industry, through this study, we are quickly discovering that the web services we use " +
		"in our day to day life are not as secure as we once imagined them to be. With different exploits like Spectre and " +
		"Meltdown, boundaries that we once thought were insurmontable, like those between two virtual machines running on " +
		"the same host, are now demonstratively not so impenetrable. One of the best goals that this blog could help to " +
		"achieve is shedding a better light on how to make distributed systems more secure than they currently are.</li>" +

		"<li class=\"list-group-item\" style=\"text-align: center;\">" +
		"<a href=\"https://github.com/etzelm\"target=\"_blank\">" +
		"<img src=\"/public/github.png\" alt=\"Github\"height=\"45\" width=\"45\"></a>" +
		"&emsp;&emsp;<a href=\"mailto:etzelm@live.com\">" +
		"<img src=\"/public/email.png\" alt=\"Email\"height=\"30\" width=\"45\"></a>" +
		"&emsp;&emsp;<a href=\"https://www.linkedin.com/in/etzelm/\" target=\"_blank\">" +
		"<img src=\"/public/linkedin.png\" alt=\"LinkedIn\"height=\"45\" width=\"45\"></a></li>" +

		"</div>" +

		"<div class=\"panel panel-default\">" +
		"<div class=\"panel-heading\" style=\"color:#A619FF;font-size: 1.17em\">Relevant Work History:</div>" +
		"<ul class=\"list-group\" style=\"font-size: 1.07em;\">" +
		"<li class=\"list-group-item\">" +
		"<b><p style=\"text-align:left;\">Software Engineer - ShieldX Networks" +
		"<span style=\"float:right;\">June 2018 - June 2019</span></p></b>" +
		"</li><li class=\"list-group-item\">" +
		"<p style=\"text-align:left;\">Software Engineer/Threat Researcher on the Security Platform:" +
		"<span style=\"float:right;\">January 2019 - June 2019</span></p>" +
		"Team Focus: Ensuring the Functionality and Quality of the ShieldX Security Offerings" +
		"<br>Main Languages: C, Regular Expressions, XML" +
		"<br>Main Threat Report Companies Worked With: Telus, Idappcom, MITRE(CVE)" +
		"<br>Main Solo Projects:" +
		"<ul><li>False Positive Fixes/Monitoring Customer Data</li>" +
		"<li>Honeypot Setup/Research into Newly Seen Exploits</li>" +
		"<li>Generate a Majority of the New 2019 HTTP Threat Coverage</li>" +
		"</ul></li><li class=\"list-group-item\">" +
		"<p style=\"text-align:left;\">Software Engineer on the Orchestration(Back-end) Team:" +
		"<span style=\"float:right;\">June 2018 - January 2019</span></p>" +
		"Team Focus: Integrating the ShieldX Product with multiple On & Off Premise Cloud Platforms" +
		"<br>Main Language: Java" +
		"<br>Main Cloud Platforms Worked With: VMware, Azure, AWS" +
		"<br>Main Solo Projects:" +
		"<ul><li>Azure Load Balancer Sandwich Setup Automation</li>" +
		"<li>Deploy Microservice Instance to a VM Folder (VMware)</li>" +
		"<li>Multi-cloud VM Tag Discovery and use with Resource Groups</li>" +
		"</ul></li>" +
		"<li class=\"list-group-item\">" +
		"<p style=\"text-align:left;\"><b>Student Internship - Product Owner" +
		"<span style=\"float:right;\">January 2017 - June 2017</span></p></b>" +
		"Internship through the UCSC Baskin School of Engineering’s Corporate Sponsorship Program:<br>" +
		"<ul><li>SmartRevenue, a market research firm based out of Connecticut, came to us as a group of students " +
		"requesting help prototyping a system for determining a consumer’s “digital path to purchase”</li>" +
		"<li>Lead a team of 4 other students as Product Owner for the Data Collection Team, we focused on the " +
		"URLs visited and Android Applications used when participants opted into our software</li>" +
		"<li>Our team developed both a Chrome Extension and Android Application to help collect data from the " +
		"pool of over 50,000 people they use in their demographic studies</li>" +
		"<li>Created a website using Amazon Web Services to act as an Administration Panel/Web Interface for our Applications and Database</li>" +
		"</ul><br></li>" +
		"</div>" +
		"</div>" +
		"<p style=\"text-align: center;\">Source code for this website can be located: " +
		"<a style=\"color:#A619FF;\" href=\"https://github.com/etzelm/blog-in-golang\" target=\"_blank\">" +
		"here</a></p>" +
		"</div>"

	info := ItemInfo{
		Title:    "About Me",
		Created:  created,
		Modified: modified,
		Blurb:    blurb,
		Content:  hold,
	}

	item := Item{
		ID:   0,
		Info: info,
	}

	av, err := dynamodbattribute.MarshalMap(item)

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("Articles"),
	}

	_, err = dbSvc.PutItem(input)

	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
