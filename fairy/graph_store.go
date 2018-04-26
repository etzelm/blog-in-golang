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

func graph_store() {
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

	blurb := "A look at what goes into making a successful distributed system"
	created := "April 10th, 2018"
	modified := "April 25th, 2018"
	hold := "<div class=\"well\" style=\"background-color:#DFF0D8;\">" +
		//Start of the About Me paragraph and contact info
		"<div class=\"panel panel-default\">" +
		"<div class=\"panel-heading\" style=\"color:#A619FF;font-size: 1.17em\">" +
		"Creating a Scalable, Fault Tolerant, & Strongly Consistent Graph Store API:</div>" +
		"<ul class=\"list-group\" style=\"font-size: 1.07em;\">" +
		"<li class=\"list-group-item\">" +
		"<p>&emsp;&emsp;Collaborating on this project with me is my former distributed systems classmate, " +
		"<a style=\"color:#9C6708;\" href=\"https://github.com/bryandmc\" target=\"_blank\">Bryan McCoid</a>. " +
		"Our inspiration for this project was a desire to recreate some of our coursework, " +
		"created by <a style=\"color:#9C6708;\" href=\"https://github.com/palvaro\" target=\"_blank\">Peter Alvaro</a> " +
		"of <a style=\"color:#9C6708;\" href=\"https://disorderlylabs.github.io/\" target=\"_blank\">Disorderly Labs</a>, " +
		"for public view and to do so by completely starting " +
		"from scratch in order to implement algorithms that we had wanted to use during the course but had " +
		"ran out of time to do so. Several attempts at academic honesty have been made and we also strongly " +
		"discourage any current students who happen upon this content from using it in their own coursework, " +
		"not only because you will be caught but because you are depriving yourself of the education you pay " +
		"for. Bryan and I are aspiring Software Engineers and we take this project on in the hope of developing " +
		"and improving the skills needed to be successful in industry, especially when it comes to dealing with " +
		"distributed systems. The goal of this project is to provide a REST-accessible graph storage service that " +
		"runs on port 3000 and is available as a resource named gs. For example, the service " +
		"would listen at <a style=\"color:#9C6708;\" href=\"http://server-hostname:3000/gs\" target=\"_blank\">" +
		"http://server-hostname:3000/gs</a>.</p><p>&emsp;&emsp;We want to develop distributed system " +
		"software to support this service so that it can store an amount of data that would " +
		"not normally fit onto a single machine system. To accomplish this, we will simulate " +
		"our server code as if it is being run on multiple, separate hosts simultaneously, " +
		"using Docker to provide this functionality. A single server host in our system stores " +
		"only a certain subset of the graphs stored in the system as a whole. We also have " +
		"them keep track of a list of all the other server hostnames in the known system so " +
		"that they can forward requests they receive for graphs that aren't stored locally for " +
		"them. The plan is to distribute graphs among partitions that each have an active " +
		"amount of server hosts assigned to them based on the total number of server hosts " +
		"that exist in the system at the time of observation. This way each server host in a " +
		"partition can store the same subset of graphs assigned to that partition, providing " +
		"a measurable amount of fault-tolerance to the user if one of those hosts happens to " +
		"crash or experience a network partition.</p><p>&emsp;&emsp;Scalability is achieved by allowing for the " +
		"user to change the system environment by adding or removing server hosts, based on " +
		"their needs, using API calls which then have our distributed system software " +
		"automatically reshuffle our partitioning and graph distribution across all active  " +
		"server hosts to attain maximum fault-tolerance and minimize access latency. To ensure " +
		"strong consistency among server hosts in a partition that stores the same subset of " +
		"graphs in our system, we will use an algorithm called Raft that uses a 2 phase commit " +
		"sequence and timers to achieve consensus on a total causal order over any value given " +
		"to us by the user. Due to the CAP theorem, we know that using partitions to attain " +
		"fault tolerance means we cannot have a graph store that is both highly available and " +
		"strongly consistent. In this project, we will favor strong consistency over having " +
		"our system be highly available, meaning our service should only return responses to " +
		"requests if it can guarantee that it is using the most recent data available to it.</p></li>" +
		"<li class=\"list-group-item\" style=\"text-align: center;\">" +
		"Source code for this project can be located here:&emsp;&emsp;" +
		"<a href=\"https://github.com/etzelm/consistent-graph-store-api\" target=\"_blank\">" +
		"<img src=\"/public/github.png\" alt=\"Github\"  height=\"45\" width=\"45\"></a>" +
		"</div>" +
		//Start of the container for
		"<div class=\"container-fluid\">" +
		//Only used one row for the container
		"<div class=\"row row-centered\">" +
		//Start of the column with face picture in it
		"<div class=\"col-md-6 form-group\">" +
		"<div class=\"panel panel-default\">" +
		"<div class=\"panel-heading\" style=\"color:#A619FF;font-size: 1.17em\">Input Format Specifications:</div>" +
		//Start of the list
		"<ul class=\"list-group\" style=\"font-size: 1.07em;\">" +
		"<li class=\"list-group-item\"><b>Graph Names: </b>" +
		"chars: [a-zA-Z0-9] i.e. Alphanumeric w/ cases" +
		"<br>size: 1 to 250 characters</li>" +
		"<li class=\"list-group-item\"><b>Vertex Names: </b>" +
		"chars: [a-zA-Z0-9] i.e. Alphanumeric w/ cases" +
		"<br>size: 1 to 250 characters</li>" +
		"<li class=\"list-group-item\"><b>Edge Names: </b>" +
		"chars: [a-zA-Z0-9] i.e. Alphanumeric w/ cases" +
		"<br>size: 1 to 250 characters</li></div></div>" +
		//Start of the column with education panel in it
		"<div class=\"col-md-6 form-group\">" +
		//Start of the education panel
		"<div class=\"panel panel-default\">" +
		"<div class=\"panel-heading\" style=\"color:#A619FF;font-size: 1.17em\">Environment Variables Used:</div>" +
		//Start of the list
		"<ul class=\"list-group\" style=\"font-size: 1.07em;\">" +
		"<li class=\"list-group-item\">" +
		"<b><i>PARTITIONS:</i></b> Tracks all active server hosts in our system</li>" +
		"<li class=\"list-group-item\">" +
		"<b><i>IP:</i></b> Stores docker network ip/port used for inter-communication</li>" +
		"<li class=\"list-group-item\">" +
		"<b><i>PORT:</i></b> Stores local network port exposed by docker for the user</li>" +
		"<li class=\"list-group-item\">" +
		"<b><i>R:</i></b> Stores max number of hosts a partition can be given</li>" +
		//End of the container for the face picture and education panel
		"</div>" +
		"</div>" +
		"</div>" +
		"</div>" +
		//Start of the About Me paragraph and contact info
		"<div class=\"panel panel-default\">" +
		"<div class=\"panel-heading\" style=\"color:#A619FF;font-size: 1.17em\">Example Docker Commands:</div>" +
		"<ul class=\"list-group\"  style=\"font-size: 1.07em;\">" +
		"<li class=\"list-group-item\">" +
		"<b>Starting a system with 4 active server hosts and a maximum partition size of 2:</b><br><br>" +
		"docker run -p 3001:3000 --ip=10.0.0.21:3000 --net=mynet -e IP=\"10.0.0.21:3000\" -e PORT=\"3001\" " +
		"-e R=2 -e PARTITIONS=\"10.0.0.21:3000,10.0.0.22:3000,10.0.0.23:3000,10.0.0.24:3000\" mycontainer<br><br>" +
		"docker run -p 3002:3000 --ip=10.0.0.22:3000 --net=mynet -e IP=\"10.0.0.22:3000\" -e PORT=\"3002\" " +
		"-e R=2 -e PARTITIONS=\"10.0.0.21:3000,10.0.0.22:3000,10.0.0.23:3000,10.0.0.24:3000\" mycontainer<br><br>" +
		"docker run -p 3003:3000 --ip=10.0.0.23:3000 --net=mynet -e IP=\"10.0.0.23:3000\" -e PORT=\"3003\" " +
		"-e R=2 -e PARTITIONS=\"10.0.0.21:3000,10.0.0.22:3000,10.0.0.23:3000,10.0.0.24:3000\" mycontainer<br><br>" +
		"docker run -p 3004:3000 --ip=10.0.0.24:3000 --net=mynet -e IP=\"10.0.0.24:3000\" -e PORT=\"3004\" " +
		"-e R=2 -e PARTITIONS=\"10.0.0.21:3000,10.0.0.22:3000,10.0.0.23:3000,10.0.0.24:3000\" mycontainer<br>" +
		"</div></div>"

	info := ItemInfo{
		Title:    "Creating a Scalable, Fault Tolerant, & Strongly Consistent Graph Store API",
		Created:  created,
		Modified: modified,
		Blurb:    blurb,
		Content:  hold,
	}

	item := Item{
		ID:   1,
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
