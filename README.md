# Movies CRUD app

This is my first ever go app.

The idea is to start small with an in-memory data structure and grow it to use a database and eventually send notifications etc, making this a playground around the golang and it's capabilities.

# How to run it
Right now you can download the code and run on your machine using:
- go build
- go run main.go

# Docker Mysql
Create a docker instance of mysql server

    docker run --name=mysql-go -p3306:3306 -e MYSQL_ROOT_HOST=% -e MYSQL_ROOT_PASSWORD=root -d mysql/mysql-server:5.7.41

Connect to it and run the script that can be found in movies.dll

    docker exec -it mysql-go mysql -uroot -p

# Kafka & Protobuf
You will need an up and running kafka cluster to be able to post created movie events.
Once you have it, create a topic called "movies".

Protobuf messages are being sent to this topic, so any streaming solution can
read from it and parse it back based on the proto message available in this repository.

The generated proto is committed in the repository, but if you want to modify and then
generate it again, you can do so by running this command:

    protoc --go_out=./internal/movies/ ./proto/movie_event.proto

# Test
For testing, please run the following command:

    go test -v ./...

This will ensure to run any tests on any directories