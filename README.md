# blog-in-golang
Started a professional blog using golang.

Published: https://app.mitchelletzel.com

Sample React App: https://app.mitchelletzel.com/realtor

Sample Local Test Commands:

 * go run main.go models.go handlers.go
 * ../../../../bin/dep.exe ensure
 * java -Djava.library.path=./DynamoDBLocal_lib -jar DynamoDBLocal.jar -sharedDb
 * go run main.go graphStore.go googleSRE.go 1

 Docker Helper Commands:

 * docker stop $(docker ps -aq)
 * docker rm $(docker ps -aq)
 * docker rmi --force $(docker images -q)
 * docker build -t blog .
 * docker run -d -p 8080:3000 blog
