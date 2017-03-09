### About

Simple learning project 

Backend: GO [Echo](https://echo.labstack.com/)

Frontend: NodeJS app
ES6, SCSS, Handlebars loader, Webpack

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

```
update Stripe API Keys
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
