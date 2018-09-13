## About
ToDo App

## Requirements
* Golang [Install](https://golang.org/doc/install)
* PostgreSQL [Install](https://www.postgresql.org/download/)
* Redis [Install](https://redis.io/download#installation)

## Install/Configure
After you have every of the above requirements:
1. Download the project:
  ```bash
    go get -u github.com/avalchev94/to-do-api
  ```
2. Update **./config/setup_env.sh** enviroment variables.
3. Open **.bashrc***(on ubuntu ~/.bashrc)* and add the following code(replace **GOPATH**):
  ```bash
  if [ -f ~/GOPATH/src/github.com/avalchev94/to-do-app/config/setup_env.sh ]; 
  then
    . ~/GOPATH/src/github.com/avalchev94/to-do-app/config/setup_env.sh
  fi
  ```
4. Create the database and import the schema:
  ```bash
    dbcreate -U $PG_USER todo_app
    #go to to-do-app main directory
    psql -U $PG_USER todo_app < config/todoapp.sql
  ```
5. Go to /api folder:
  ```bash
  go build*
  ./api
  ```
# API Documentation 
* **GET /login** Returns the currently logged user:
  ```json
  {"id": 1, "name":"avalchev", "email":"avl@bg.bg", "avatar":...}
  ```
* **POST /login** Start new login session:
  ```json
  {"name":"some_name", "password":"the_password"}
  ```
* **POST /user** Register new user:
  ```json
  {"name":"some_name", "password":"...", "email":"...", "avatar":"..."}
  ```
* **GET /user/:id** Returns the specified user: 
  ```json
  {"id": 1, "name":"avalchev", "email":"avl@bg.bg", "avatar":...}
  ```
* **GET /user/:id/labels** Returns the labels that the user created:
  ```json
  [
    {"id":1,"user_id":1,"name":"Work","color":"123456"},
    {"id":2,"user_id":1,"name":"School","color":"ABCDEF"}
  ]
  ```
* **GET /user/:id/tasks** Returns the tasks that the user created:
  ```json
  [{"id":2,"user_id":1,"title":"sth","data":"I should..","labels":[{"id":1,"name":"Work","color":"123456"}],"labels_id":[1],"created":"2018-09..."}]
  ```
* **POST /label** Create new label:
  ```json
  {"name":"Work", "color":"000fff"}
  ```
* **GET /label/:id** Get label with id:
  ```json
  {"id":1,"user_id":1,"name":"Work","color":"123456"}
  ```
* **POST /task** Create new task:
  ```json
  {"title":"..", "data":"what u should do", "labels":[1, 5, 10]}
  ```
* **GET /task/:id** Get task with id:
  ```json
  {"id":2,"user_id":1,"title":"Finish..","data":"I should finish..","labels":[{"id":1,"name":"Work","color":"123456"}],"labels_id":[1],"created":"2018-09..."}
  ```