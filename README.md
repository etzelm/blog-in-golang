
# blog-in-golang

Started a professional blog using golang.

Published: <https://mitchelletzel.com>

Sample React App: <https://mitchelletzel.com/realtor>

Sample Local Test Commands:

* cd realtor && npm install
* yarn build
* cd ../ && go mod download
* java -Djava.library.path=./DynamoDBLocal_lib -jar DynamoDBLocal.jar -sharedDb
* go run app.go
* cd daemon && go run app.go 1

 Docker Helper Commands:

* docker stop $(docker ps -aq)
* docker rm $(docker ps -aq)
* docker rmi --force $(docker images -q)
* docker build -t blog:develop .
* docker run -d -p 8080:3000 blog
* docker compose up --force-recreate -d
* docker compose down
* docker system prune -a
