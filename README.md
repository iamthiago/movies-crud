# Movies CRUD app

This is my first ever go app.

The idea is to start small with an in-memory data structure and grow it to use a database and eventually send notifications etc, making this a playground around the golang and it's capabilities.

# How to run it
Right now you can download the code and run on your machine using:
- go build
- go run main.go

# Docker Mysql
docker run --name=mysql -p3306:3306 -e MYSQL_ROOT_HOST=% -e MYSQL_ROOT_PASSWORD=root -d mysql/mysql-server:5.7.41