# What is Hexagonal Architecture?

Core Idea: Your business logic should not depend on external things (database, HTTP, APIs)
Instead, external things should depend on your business logic

Why Hexagonal? The hexagon shape is just a visual to represent that your core buisness logic is in the center,
and external things connect to it form all sides through "ports"

The example given will be based of a mini microservice project task-service. The basic project outline is defined below:
```
Personal Task Management System Project

## Services
1. User service - Handle user registration/authentication
2. Task service - Create, Update, list, delete tasks
3. Notification service - Send simple notifications when tasks are created/completed
```

## The Three Layers
1. Domain Layer (Center of Hexagon)
What it is: Your core business rules and entities, it contains:
* Business entities (struct like Task, User)
* Business rules (validation, calculations)
* Domain errors

Key rule: This layer knows NOTHING about HTTP, databases, or external services
```go
// 
```
https://claude.ai/chat/de06bf0f-543f-4d8f-a95f-284afbb6eab0
