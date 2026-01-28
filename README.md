# raft-chat

A basic chat service written in Go for demonstration purposes. The service provides HTTP endpoints for managing chat rooms and messages, with in-memory storage.

## Overview

raft-chat is a simple HTTP service that allows you to:
- Create chat rooms with names and descriptions
- Send messages to rooms
- Retrieve messages from rooms (up to 20 most recent)

All data is stored in memory and messages are stored chronologically by their receipt time.

## Architecture

The service is organized into clear layers:

- **API Layer** (`api/`): HTTP request handlers and routing using Go's `ServeMux`
- **Storage Layer** (`storage/`): In-memory data storage with a single-writer pattern for thread-safe message writes
- **Models** (`models/`): Core data structures (Room and Message)
- **Main** (`main.go`): Service initialization and wiring

## Getting Started

### Prerequisites

- Go 1.22 or later

### Running Locally

1. Clone the repository:

```bash
git clone https://github.com/max2sax/raft-chat.git
cd raft-chat
```

2. Build the service:

```bash
go build
```

3. Run the service:

```bash
./raft-chat
```

The service will start on `http://localhost:8080`

## API Endpoints

### Create a Room

```text
POST /rooms
Content-Type: application/json

{
  "name": "general",
  "description": "General discussion"
}
```

Response: `201 Created`

```json
{
  "id": "general",
  "name": "general",
  "description": "General discussion"
}
```

### Get All Rooms

``` text
GET /rooms
```

Response: `200 OK`

```json
[
  {
    "id": "general",
    "name": "general",
    "description": "General discussion"
  },
  {
    "id": "random",
    "name": "random",
    "description": "Random discussions"
  }
]
```

### Get a Room

``` text
GET /rooms/{roomID}
```

Response: `200 OK`

```json
{
  "id": "general",
  "name": "general",
  "description": "General discussion"
}
```

### Send a Message

```text
POST /rooms/{roomID}/messages
Content-Type: application/json

{
  "sender": "alice",
  "content": "Hello world!"
}
```

Response: `201 Created`

### Get Messages

```text
GET /rooms/{roomID}/messages
```

Response: `200 OK`

```json
[
  {
    "id": "01ARZ3NDEKTSV4RRFFQ69G5FAV",
    "roomName": "general",
    "timestamp": "2026-01-28T12:34:56.789Z",
    "sender": "alice",
    "content": "Hello world!"
  }
]
```

## Testing

You can test the endpoints using `curl`:

```bash
# Create a room
curl -X POST http://localhost:8080/rooms \
  -H "Content-Type: application/json" \
  -d '{"name":"general"}'

# Get all rooms
curl http://localhost:8080/rooms

# Get a specific room
curl http://localhost:8080/rooms/general

# Send a message
curl -X POST http://localhost:8080/rooms/general/messages \
  -H "Content-Type: application/json" \
  -d '{"sender":"user1","content":"Hello!"}'

# Get messages
curl http://localhost:8080/rooms/{roomID}/messages
```

## Design Notes

- **Message Concurrency**: Messages are written to storage through a single-writer pattern using channels, ensuring thread-safe concurrent access without race conditions
- **Room IDs**: Room names are used as unique identifiers
- **Message IDs**: Uses ULIDs for messages, providing sortable, timestamp-based identifiers
- **Message Limit**: Returns only the 20 most recent messages per room
- **Empty Collections**: Returns empty arrays instead of null when no data exists

## Future features

- **Search messages within and across rooms**
- **Add testing**: Start out with an end to end testing framework for testing integration of the api
- **Allow deleting rooms which will delete all messages as well**
- **Allow editing of messages**
- **Create users validate that the message senders are valid users**
- **Add database for permanent storage**
- **Add more logging**
- **...**
- **Scale and take over the world**
