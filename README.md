# Tic-Tac-Toe REST API

## Getting Started

**Prerequisites:** Go 1.26+

1. **Run the server:**
   ```bash
   go run main.go
   ```
   *The server will start on `http://localhost:8080`*

2. **Send any request from `requests.http`**

## Testing the API

`requests.http` file contains http-requests to server.

### Option 1: VS Code
1. Install the REST Client extension.
2. Open `requests.http`.
3. Click Send Request above any method.
   *Note: The `session_id` from the creation step is automatically captured and used in subsequent requests.*

### Option 2: Curl
To verify the server is running, try creating a game:
```bash
curl -X POST http://localhost:8080/game/create \
     -H "Content-Type: application/json" \
     -d '{"player_id": "uuid-player-1", "player_name": "player1", "rules": {"board_width": 3, "win_length": 3}, "seed": 0}'
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/game/create` | Initialize a new game session |
| `POST` | `/game/{session_id}/connect` | Join an existing game as the second player |
| `POST` | `/game/{session_id}/start` | Start the game once both players are ready |
| `POST` | `/game/{session_id}/move` | Submit a move (x, y coordinates) |
| `POST` | `/game/{session_id}/disconnect` | Leave the game session |
