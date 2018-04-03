for i in {1..1} ;
do
    ( docker stop $(docker ps -aq) );
    ( docker rm $(docker ps -aq) );
    ( docker build -t blog . );
    ( docker run -d -p 8080:3000 blog );
    ( docker ps -a );
done