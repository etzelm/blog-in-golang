# blog-in-golang
Starting a blog using golang.

Sample Local Test Commands:

 * go run main.go models.article.go handlers.article.go
 * ../../../../bin/dep.exe ensure
 * java -Djava.library.path=./DynamoDBLocal_lib -jar DynamoDBLocal.jar -sharedDb
 * go run main.go about_me.go how_its_made.go drop_tables.go 1

 Docker Helper Commands:

 * docker stop $(docker ps -aq)
 * docker rm $(docker ps -aq)
 * docker rmi --force $(docker images -q)
 * docker build -t hw4 .
 * docker run -d hw4

Color Scheme:

 * Current Green: 3CB371
 * Complimentary Purple: A619FF
 * Related Light Green: 00FF72
 * Complimentary Light Tan: CC8A14
 * Complimentary Dark Tan: 9C6708