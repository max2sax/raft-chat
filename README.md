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
```
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
  "id": "01ARZ3NDEKTSV4RRFFQ69G5FAV",
  "name": "general",
  "description": "General discussion"
}
```

### Send a Message
```
POST /rooms/{roomID}/messages
Content-Type: application/json

{
  "sender": "alice",
  "content": "Hello world!"
}
```

Response: `201 Created`

### Get Messages
```
GET /rooms/{roomID}/messages
```

Response: `200 OK`
```json
[
  {
    "id": "01ARZ3NDEKTSV4RRFFQ69G5FAV",
    "roomID": "01ARZ3NDEKTSV4RRFFQ69G5FAV",
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

# Send a message (replace {roomID} with the ID from above)
curl -X POST http://localhost:8080/rooms/{roomID}/messages \
  -H "Content-Type: application/json" \
  -d '{"sender":"user1","content":"Hello!"}'

# Get messages
curl http://localhost:8080/rooms/{roomID}/messages
```

## Design Notes

- **Message Concurrency**: Messages are written to storage through a single-writer pattern using channels, ensuring thread-safe concurrent access without race conditions. A syncmap is used for the message storage to prevent reads during a write operation.
- **Message IDs**: Uses ULIDs for both rooms and messages, providing sortable, timestamp-based identifiers
- **Message Limit**: Returns only the 20 most recent messages per room

## Future features

- **Search messages within and across rooms**
- **Add testing**: Start out with an end to end testing framework for testing integration of the api
- **Allow deleting rooms which will delete all messages as well**
- **Allow editing of messages**
- **Create users validate that the message senders are valid users**
- **Add database for permanent storage**
- **...**
- **Scale and take over the world**
