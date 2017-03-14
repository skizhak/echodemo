### About

Simple learning project 

Backend: GO [Echo](https://echo.labstack.com/)

Frontend: NodeJS app
ES6, SCSS, Handlebars, Material Design, Webpack

### Installation

clone this repo.

```
cd echodemo
go get -u github.com/kardianos/govendor
govendor sync
```

```
cd web
npm install
npm run build 

(npm run dev for dev environment)
```
MySql

setup mysql and create user and database.
user: 'echodemo' with password 'demo123'
db: 'echodemo'

```
mysql -u root -p -e "create user 'echodemo'@'localhost' identified by 'demo123'; grant all privileges on * . * to 'echodemo'@'localhost'; flush privileges; drop database if exists echodemo; create database echodemo;"
```
Let's use goose for db migration.
```
go get bitbucket.org/liamstask/goose/cmd/goose
```
check you're able to connect to db using the config file under db
```
goose status
```
Now, apply the migration
```
goose up
```

Stripe API

```
update with your Stripe API Keys in resources/controller.go and web/dist/index.html
```

Start the server
```
go run main.go
```

Load broswer, point to [http://localhost:9001](http://localhost:9001)


APIs
```
"/api/hello"
    GET: Simple "hello world"
		
"/users"
    GET : List of users
    POST: Create new user

"/users/:id"
    GET: Get a particular user

"/users/:id/payments"
    GET:  Get payments of a user
    POST: Post payment for a user
```

### Dev  Note
sqlboiler is part of vendor, but if new models needs to be generated, install it
```
go get github.com/vattle/sqlboiler
```
update the db tables, and run
```
go generate
```