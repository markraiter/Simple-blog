# REST API for simple blog (NIX Academy Trainee level task)

###The task:

* Receiving from the service [JSONPlaceholder](https://jsonplaceholder.typicode.com) posts with userid=7, and receiving comments from these posts with postId={postId} using package [net/http](https://pkg.go.dev/net/http);

* Recording of received posts and comments to the MySQL database using [GORM](https://gorm.io/);

* Creation of CRUD for posts and comments. The response from the API must implement the data presentation format: XML, JSON using [Echo framework](https://echo.labstack.com/);

* Adding [SWAGGER](https://swagger.io/) to the API;

* Adding the ability to register, authorize users using the [JWT standard](https://jwt.io/);

* Writing tests for the API, using the standard library for testing - [testing](https://pkg.go.dev/testing).

###Tech stack:

* Go;
* Echo Framework;
* REST; 
* MySQL; 
* GORM; 
* JWT; 
* Swagger; 
* Git.

> To see swagger docs please proceed this link: [swagger](http://localhost:8080/swagger/index.html)