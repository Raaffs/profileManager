<a name="readme-top"></a>


<br />

<div align="center">

  <h3 align="center">Identity Management Service</h3>

  <p align="center">
</div>



<!-- TABLE OF CONTENTS -->


<!-- ABOUT THE PROJECT -->
## About The Project
This is a full-stack identity management system built to handle user registration, login, and secure profile storage. The project is split into three main parts: a Go backend API, a React (Vite) frontend dashboard, and a PostgreSQL database.
The core focus of this system is the Protection of Personally Identifiable Information (PII). Specifically, it implements AES-256 encryption to ensure that sensitive data like Aadhaar/National ID numbers remain encrypted even if the database is compromised.


### Built With

* [![Go][Go]][Go-url]
* [![React][React.js]][React-url]
* [![psql][psql]][psql-url]
* [![docker][docker]][docker-url]
---------
## Implementation approach and core logic 

### 1.  Backend Implementation

a. **Authentication & Authorization**  
  - Users are authenticated via JWTs, generated using `GenerateToken(userID)` and validated in each request with `GetUserJWT()`.  
  - JWT claims include the `UserID` and standard JWT expiration.  
  - Tokens use the HS256 signing algorithm with a 72-hour expiration.  

b. **Data Layer**  
  - User and profile data is managed through separate repository interfaces (`UserRepository` and `ProfileRepository`) for clean separation of concerns.  
  - CRUD operations are abstracted behind the repository layer to allow easy swapping of database backends.  
  - Database schema is initialized via the `migrations/` folder; `init.sql` is used for Docker container setup.  

c. **Data Security**  
  - Sensitive fields are encrypted using AES-GCM via `EncryptFields` and `DecryptFields`.  
  - Uses Go’s standard libraries: `crypto/aes` (AES block cipher), `crypto/cipher` (GCM mode), and `crypto/rand` (secure nonces).  
  - AES-256 secret keys are stored securely via environment variables.  
  - Passwords are hashed using `bcrypt` with a nonce to protect against brute-force and rainbow table attacks.  

d. **Input Validation**  
  - Data is validated using the `Validator` utility before saving to the database.  
  - Checks include name length, Aadhaar number format, phone number format, and date correctness.  
  - Errors are collected in an `Errors` map for consistent handling of invalid inputs.  

e. **Error Handling**  
  - Standard HTTP errors are defined via `HttpResponseMsg` constants (`ErrBadRequest`, `ErrUnauthorized`, etc.).  
  - Repository methods return custom errors (`NotFound`, `AlreadyExists`) to provide a consistent interface for the service layer.  
  - DB-specific errors (e.g., `pgx.ErrNoRows`, unique constraint violations) are mapped to these custom errors, keeping the service layer database-agnostic.  

### 2. Frontend Implementation

a. **API Integration**  
  - Axios is used for HTTP requests, with the base URL dynamically set from `VITE_API_BASE_URL` and a fallback to `http://localhost:8080/api`.  
  - JWT tokens from `localStorage` are automatically attached to requests via an Axios request interceptor (`Authorization: Bearer <token>`).  

b. **Routing & Route Protection**  
  - Routes are managed with `react-router-dom`.  
  - The `ProtectedRoute` component ensures only authenticated users can access certain pages (e.g., Profile).  
  - Unauthenticated users are redirected to the login page, and unknown routes fall back to a 404 error page.
    
c. **Form Handling & Validation**  
- Forms are implemented using `Formik` for state management and submission handling.  
- Input validation is done with `Yup` schemas, enforcing rules like name length, Aadhaar number format, phone number format, and valid dates.  
- On submit, forms either create or update profiles via the API (`POST` or `PUT` requests to `restricted/profile`).  
- Server responses are handled gracefully, showing success or error messages based on API results.


d. **Theme & Styling**  
- The app uses Material-UI (`@mui`) with a custom theme applied via `ThemeProvider`.  
- Input fields are styled with `textFieldSx` and `activeTextFieldSx` for focused, hover, and disabled states.  
- Buttons use `btnstyle` with gradients, rounded borders, shadows, and hover effects for a modern look.

  ---
## Development Setup Guide

Follow these instructions to get the project up and running on your local machine.

### 🛠 Prerequisites

Ensure you have the following installed before proceeding:

* **Go** (Golang): [Download & Install](https://go.dev/doc/install)
* **Node.js**: [Download & Install](https://nodejs.org/)
* **PostgreSQL**: [Download & Install](https://www.postgresql.org/download/)
* **Docker & Docker Compose**: [Get Docker](https://docs.docker.com/get-docker/)
---
### Setup Instructions

#### Option 1: Using Docker (Recommended)
  ```bash
  docker compose --env-file .env.example up --build
  ```
---
#### Option 2: Manual Local Setup
1. Environment Configuration
     Create your local environment file and update the database connection string.
    ```bash
    cp .env.example .env
    # Open .env and edit your DATABASE_URL
    ```
2. Database Migrations
   ```bash
     go install -tags 'postgres' [github.com/golang-migrate/migrate/v4/cmd/migrate@latest](https://github.com/golang-migrate/migrate/v4/cmd/migrate@latest)
   ```
    ```bash
    migrate -path server/migrations/ -database "your_db_url" up
    ```
3. Backend Setup
   Install Go modules and start the server.
    ```bash
    go mod tidy
    ```
    ```bash
    go run ./server/cmd/web
    ```
4. Frontend Setup
   ```bash
    cd client
    npm install
    npm run dev
    ```
## API Endpoints
| Endpoint | Method | Auth Required | Request Body (JSON) | Description & Key Logic |
| :--- | :--- | :---: | :--- | :--- |
| `/api/login` | `POST` | ❌ No | `{"email": "...", "password": "..."}` | Authenticates user and returns a JWT token. |
| `/api/register` | `POST` | ❌ No | `{"email": "...", "password": "..."}` | Creates a new user account in the database. |
| `/api/restricted/profile` | `GET` | ✅ Yes | None | Fetches the profile associated with the authenticated user ID. |
| `/api/restricted/profile` | `POST` | ✅ Yes | `{"full_name": "...", "date_of_birth": "...", "aadhaar_number": "...", "phone_number": "...", "address": "..."}` | Initializes a new profile record for the authenticated user. |
| `/api/restricted/profile` | `PUT` | ✅ Yes | `{"full_name": "...", "date_of_birth": "...", "aadhaar_number": "...", "phone_number": "...", "address": "..."}` | Updates existing profile details. Validates via JWT `sub` claim. |

[Go]: https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white
[Go-url]: https://go.dev/
[React.js]: https://img.shields.io/badge/React-20232A?style=for-the-badge&logo=react&logoColor=61DAFB
[React-url]: https://reactjs.org/
[psql]: https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white
[psql-url]: https://www.postgresql.org/
[docker]: https://img.shields.io/badge/Docker-2CA5E0?style=for-the-badge&logo=docker&logoColor=white
[docker-url]: https://www.docker.com/
