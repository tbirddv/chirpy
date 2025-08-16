# Chirpy API

Chirpy is a simple social media API for posting short messages ("chirps"), user authentication, and metrics tracking.

## Endpoints

### Health Check
- **GET /api/healthz**  
  Returns `OK` if the server is running.

### Chirps
- **POST /api/chirps**  
  Create a new chirp (requires authentication).

### Health Check
- **GET /api/healthz**  
  Returns `OK` if the server is running.

### Chirps
- **POST /api/chirps**  
  Create a new chirp (requires authentication).
- **GET /api/chirps**  
  List all chirps. Supports optional query parameter:
  - `author_id`: Filter chirps by author (e.g. `/api/chirps?author_id=123`)
- **GET /api/chirps/{id}**  
  Get a chirp by ID.

### Users
- **POST /api/users**  
  Register a new user.

### Authentication
- **POST /api/login**  
  Log in and receive access/refresh tokens.
- **POST /api/refresh**  
  Refresh your access token using a refresh token.
- **POST /api/revoke**  
  Revoke a refresh token.

### Admin
- **GET /admin/metrics**  
  View server metrics (file server hits).
- **POST /admin/reset**  
  Reset metrics and delete all users (dev platform only).

## License

MIT

---

For more details, see [main.go](main.go) and the endpoint handlers in [chirps.go](chirps.go), [users.go](users.go), and [auth.go](auth.go).

