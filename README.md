# blog-in-golang
Starting a blog using golang.

Sample Local Test Commands:

 * go run main.go models.go handlers.go
 * ../../../../bin/dep.exe ensure
 * java -Djava.library.path=./DynamoDBLocal_lib -jar DynamoDBLocal.jar -sharedDb
 * go run main.go about_me.go graph_store.go google_SRE.go drop_tables.go create_table.go 1

 Docker Helper Commands:

 * docker stop $(docker ps -aq)
 * docker rm $(docker ps -aq)
 * docker rmi --force $(docker images -q)
 * docker build -t blog .
 * docker run -d -p 8080:3000 blog

Color Scheme:

 * Current Green: 3CB371
 * Complimentary Purple: A619FF
 * Related Light Green: 00FF72
 * Complimentary Light Tan: CC8A14
 * Complimentary Dark Tan: 9C6708