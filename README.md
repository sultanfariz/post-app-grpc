# Post App 
This is a simple post application that I developed as a part of learning gRPC in Go. The application has both server and client components.

Server
The server component of the app has two services: auth and post. The auth service provides login and registration functionality, while the post service provides the following methods:

GetAllPosts
GetPostById
CreatePost
DeletePost
The post entity has the following fields:

id
title
content
topic
created_at
updated_at
The CreatePost method pushes the created post to RabbitMQ for each topic queue.

The server app uses Clean Architecture and MySQL as RDBMS with Gorm as the ORM. An interceptor is implemented to act as an auth middleware to check if the user is logged in.

Client
The client component of the app is a gRPC client that also serves a RESTful API through an API gateway. It provides the following methods:

GetAllPosts
GetPostById
SubscribePostByTopic
SubscribePostByTopic works as a server-stream and sends new posts by the requested topic from the message broker.

The client app is implemented using buf as the compiler.

System Design Diagram
System Design Diagram

Note: The client app does not provide any authentication functionality.

Usage
Clone the repository.
Run go run main.go to start the server and go run client/main.go to start the client.