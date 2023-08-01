# REST API for Simple Blog Project

### The task:

This is a simple-blog application with jwt authentication and CRUD operations for posts and comments. Only registered users can create posts and comments.

### Tech stack:

* Go;
* Gin-Gonic;
* REST; 
* Postgres; 
* Golang-Migrate;
* JWT; 
* Docker;
* Git.

### To run the blog please proceed next:

1. Initialise database:
``` BASH
make databaseinit
```
2. Make migrations:
``` BASH
make migrate_up
```
3. Run the Blog:
``` BASH
make run
```
4. Blog will run on http://localhost:8080

## API Documentation

### Authentication Endpoints

#### Sign Up

- Endpoint: **POST /auth/sign-up**
- Description: Allows users to sign up by providing their registration information.
- Request Body:

```JSON
{
    "email": "user@example.com",
    "password": "password123"
}
```

- Response:
    - **201 Created**: User successfully registered. Returns the id of the created user.
    - **400 Bad Request**: Invalid request or missing required fields.
    - **409 Conflict**: User already exists with the provided email or username.
    - **500 Internal Server Error**: An unexpected error occurred on the server.

#### Sign In
- Endpoint: **POST /auth/sign-in**
- Description: Allows registered users to sign in and obtain a JWT token for API access.
- Request Body:

```JSON 
{
    "email": "user@example.com",
    "password": "password123"
}
```

- Response:
    - **200 OK**: Authentication successful. Returns a JWT token in the response body.
    - **404 Not Found**: Invalid credentials.
    - **500 Internal Server Error**: An unexpected error occurred on the server.

### API Endpoints

> Note: All API endpoints below require a valid JWT token to be included in the Authorization header of the request.

#### Posts

##### Get All Posts

- Endpoint: **GET /api/posts/all**
- Description: Retrieves all posts available in the system.
- Response:
    - **200 OK**: Returns an array of posts.
    - **500 Internal Server Error**: An unexpected error occurred on the server.

##### Filter Posts By User

- Endpoint: **GET /api/posts**
- Description: Retrieves posts created by the authenticated user.
- Query Parameters:
    - user_id: ID of the user to filter posts.
    `?user_id=?`
- Response:
    - **200 OK**: Returns an array of posts created by the user.
    - **404 Not Found**: Post with specified user_iddoes not exist.
    - **500 Internal Server Error**: An unexpected error occurred on the server.

##### Get Post By ID

- Endpoint: **GET /api/posts/:id**
- Description: Retrieves a specific post by its unique ID.
- Response:
    - **200 OK**: Returns the post object.
    - **404 Not Found**: Post with the specified ID does not exist.
    - **500 Internal Server Error**: An unexpected error occurred on the server.

##### Create Post

- Endpoint: **POST /api/posts**
- Description: Creates a new post for the authenticated user.
- Request Body:
```JSON 
{
  "title": "Example Post",
  "body": "This is the content of the post."
}
```

- Response:
    - **201 Created**: Post successfully created. Returns the newly created post id.
    - **400 Bad Request**: Invalid request or missing required fields.
    - **500 Internal Server Error**: An unexpected error occurred on the server.

##### Update Post

- Endpoint: **PATCH /api/posts/:id**
- Description: Updates an existing post owned by the authenticated user.
- Request Body:
``` JSON
{
  "title": "Updated Post Title",
  "body": "This is the updated content of the post."
}
```

- Response:
    - **200 OK**: Post successfully updated.
    - **400 Bad Request**: Invalid request or missing required.
    - **404 Not Found**: Post with the specified ID does not exist.
    <!-- - 403 Forbidden: The authenticated user does not own the post. -->
    - **500 Internal Server Error**: An unexpected error occurred on the server.

##### Delete Post

- Endpoint: **DELETE /api/posts/:id**
- Description: Deletes an existing post owned by the authenticated user.
- Response:
    - **200 OK**: Post successfully deleted.
    - **404 Not Found**: Post with the specified ID does not exist.
    <!-- - 403 Forbidden: The authenticated user does not own the post. -->
    - **500 Internal Server Error**: An unexpected error occurred on the server.

#### Comments

##### Get All Comments

- Endpoint: **GET /api/comments/all**
- Description: Retrieves all comments available in the system.
- Response:
    - **200 OK**: Returns an array of comments.
    - **500 Internal Server Error**: An unexpected error occurred on the server.

##### Filter Comments By Post

- Endpoint: **GET /api/comments/post**
- Description: Retrieves comments associated with a specific post.
- Query Parameters:
    - post_id: ID of the post to filter comments.
    `?post_id=?`
- Response:
    - **200 OK**: Returns an array of comments associated with the post.
    - **404 Not Found**: Post with the specified ID does not exist.
    - **500 Internal Server Error**: An unexpected error occurred on the server.

##### Filter Comments By User

- Endpoint: **GET /api/comments/user**
- Description: Retrieves comments created by the authenticated user.
- Query Parameters:
    - user_id: ID of the user to filter comments.
    `?user_id=?`
- Response:
    - **200 OK**: Returns an array of comments created by the user.
    - **404 Not Found**: Post with the specified ID does not exist.
    - **500 Internal Server Error**: An unexpected error occurred on the server.

##### Get Comment By ID

- Endpoint: **GET /api/comments/:id**
- Description: Retrieves a specific comment by its unique ID.
- Response:
    - **200 OK**: Returns the comment object.
    - **404 Not Found**: Comment with the specified ID does not exist.
    - **500 Internal Server Error**: An unexpected error occurred on the server.

##### Create Comment

- Endpoint: **POST /api/comments**
- Description: Creates a new comment for a post by the authenticated user.
- Query Parameters:
    - post_id: ID of the post to create a comment.
    - user_id: ID of the user to create a comment.
    `?post_id=?&user_id=?`
- Request Body:
``` JSON
{
    "email": "user@example.com",
    "body": "This is a comment on the post."
}
```

- Response:
    - **201 Created**: Comment successfully created. Returns the newly created comment id.
    - **400 Bad Request**: Invalid request or missing required fields.
    - **404 Not Found**: Post or User with the specified ID does not exists.
    - **500 Internal Server Error**: An unexpected error occurred on the server.

##### Update Comment

- Endpoint: **PATCH /api/comments/:id**
- Description: Updates an existing comment owned by the authenticated user.
- Request Body:
``` JSON
{
    "email": "user@example.com",
    "body": "This is the updated content of the comment."
}
```

- Response:
    - **200 OK**: Comment successfully updated.
    - **400 Bad Request**: Invalid request or missing required.
    - **404 Not Found**: Comment with the specified ID does not exist.
    <!-- - 403 Forbidden: The authenticated user does not own the comment. -->
    - **500 Internal Server Error**: An unexpected error occurred on the server.

##### Delete Comment

- Endpoint: **DELETE /api/comments/:id**
- Description: Deletes an existing comment owned by the authenticated user.
- Response:
    - **200 OK**: Comment successfully deleted.
    - **404 Not Found**: Comment with the specified ID does not exist.
    <!-- - 403 Forbidden: The authenticated user does not own the comment. -->
    - **500 Internal Server Error**: An unexpected error occurred on the server.