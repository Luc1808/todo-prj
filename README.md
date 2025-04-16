# To-Do API

A RESTful API for managing tasks and to-do items. This project serves as a learning experience for building APIs from scratch.

## Overview

This To-Do API provides a backend service for task management applications, allowing clients to create, read, update, and delete tasks. It includes features like task categorization, priority levels, filtering, and pagination.

## Core Features

### Data Model

Each task in the system has the following properties:

- `id`: Unique identifier (automatically generated)
- `title`: Short description of the task
- `description`: Optional detailed information
- `completed`: Boolean indicating completion status
- `dueDate`: When the task should be completed by
- `priority`: Level of importance (e.g., "low", "medium", "high")
- `category`: Tag or grouping for the task
- `createdAt`: Timestamp of when the task was created

### API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/tasks` | Create a new task |
| GET | `/tasks` | List all tasks (with pagination) |
| GET | `/tasks/:id` | Get a specific task by ID |
| PUT | `/tasks/:id` | Update an existing task |
| DELETE | `/tasks/:id` | Remove a task |
| GET | `/tasks/completed` | List all completed tasks |
| GET | `/tasks/category/:category` | List tasks by category |
| GET | `/tasks/priority/:level` | List tasks by priority level |

### Query Parameters

For the list endpoints (`GET /tasks`), the following query parameters are supported:

- `page`: Page number (default: 1)
- `limit`: Number of items per page (default: 10)
- `sortBy`: Field to sort by (default: "createdAt")
- `sortOrder`: "asc" or "desc" (default: "desc")
- `completed`: Filter by completion status (true/false)
- `search`: Search term to filter tasks by title or description

