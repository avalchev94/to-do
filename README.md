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
  go build
  ./api
  ```
# API Documentation
* POST /register 
  ```json
  {"name":"some_name", "password":"...", "email":"...", "avatar":"..."}
  ```
* POST /login
  ```json
  {"name":"some_name", "password":"the_password"}
  ```
* GET /user
  Returns the current logged user.
* Get /user/:id
  Returns the specified user.