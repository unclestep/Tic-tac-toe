# Tic-Tac-Toe REST API

Deployed application: [Tic-Tac-Toe on Railway.com](https://tic-tac-toe-production-5faf.up.railway.app)

> **Note:** The root URL (`/`) returns 404 — this is expected, the API has no UI.
> Use [Swagger UI](https://tic-tac-toe-production-5faf.up.railway.app/swagger/index.html) or `requests.http` to interact with the API.

## Getting Started

**Prerequisites:** Go 1.26+

1. **Run the server:**
   ```bash
   go run main.go
   ```
   *The server will start on `http://localhost:8080`*

2. **Open Swagger UI:**
   `http://localhost:8080/swagger/index.html`

## Testing the API

`requests.http` file contains http-requests to server.

### Option 1: Swagger UI (recommended)
Open `http://localhost:8080/swagger/index.html` in your browser.
All endpoints are documented and can be tested directly from the UI.

### Option 2: VS Code
1. Install the REST Client extension.
2. Open `requests.http`.
3. Click Send Request above any method.
   *Note: The `session_id` from the creation step is automatically captured and used in subsequent requests.*

### Option 3: Curl
To verify the server is running, try creating a game:
```bash
curl -X POST http://localhost:8080/game/create \
     -H "Content-Type: application/json" \
     -d '{"player_id": "uuid-player-1", "player_name": "player1", "rules": {"board_width": 3, "board_height": 3, "win_length": 3}, "seed": 0}'
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/game/create` | Initialize a new game session |
| `POST` | `/game/{session_id}/connect` | Join an existing game as the second player |
| `POST` | `/game/{session_id}/start` | Start the game once both players are ready |
| `POST` | `/game/{session_id}/move` | Submit a move (x, y coordinates) |
| `POST` | `/game/{session_id}/disconnect` | Leave the game session |
