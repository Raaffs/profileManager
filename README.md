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

## GitHub Link:
https://github.com/Raaffs/IdentityService

---

### Built With

* [![Go][Go]][Go-url]
* [![React][React.js]][React-url]
* [![psql][psql]][psql-url]
* [![docker][docker]][docker-url]
---------
## Usage
### 1. Registration
<img width="1883" height="944" alt="image" src="https://github.com/user-attachments/assets/50744042-109e-4268-97b1-4177d82df7c3" />

---
### 2. Login
<img width="1883" height="944" alt="image" src="https://github.com/user-attachments/assets/533d91ba-a054-4872-91e1-c8e17ee0b238" />

---
### 3. Create Profile
<img width="1875" height="938" alt="image" src="https://github.com/user-attachments/assets/c21b843d-cfd6-4f2a-831c-b46e27c1348e" />

--- 
### 4. Get Profile
<img width="1867" height="936" alt="image" src="https://github.com/user-attachments/assets/ffd993b9-7a7a-410c-aa3c-f3a34f571cf1" />

---
### 5. Update Profile
<img width="1869" height="933" alt="image" src="https://github.com/user-attachments/assets/1563c6da-fe2a-4a1b-af78-e53004a540df" />

---

## 6. Implementation approach and core logic 

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

---

### Option 1: Using Docker (Recommended)
####  Prerequisites
1. Install **Docker & Docker Compose**: [Get Docker](https://docs.docker.com/get-docker/)

2. Build & run using following command:
   ```bash
    docker compose --env-file .env.example up --build
    ```
  The project should be live on localhost:3000
  
---
### Option 2: Manual Local Setup
####  Prerequisites

Ensure you have the following installed before proceeding:

* **Go** (Golang): [Download & Install](https://go.dev/doc/install)
* **Node.js**: [Download & Install](https://nodejs.org/)
* **PostgreSQL**: [Download & Install](https://www.postgresql.org/download/)

1. Environment Configuration
   - Create your local environment file and update the database connection string.

    ```bash
    cp .env.example .env
    # Open .env and edit your DB_URL
    ```
3. Database Migrations
   ```bash
     go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
   ```
    ```bash
    migrate -path server/migrations/ -database "your_db_url" up
    ```
4. Backend Setup
   Install Go modules and start the server.
    ```bash
    go mod download
    ```
    ```bash
    go run ./server/cmd/web
    ```
5. Frontend Setup
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

---
## AI USAGE LOG
### 1. AES Encryption & Test Scenarios  
**AI Agent:** ChatGPT 5 Mini (Copilot)  
**Score:** ⭐⭐⭐⭐⭐ (5 / 5)

- Correctly implemented AES-GCM encryption using base64-decoded keys and secure nonce generation.
- Assisted in creating comprehensive test cases covering edge scenarios such as empty inputs, invalid keys, repeated encryption, and mismatched keys.

---

### 2. Backend Error Handling  
**AI Agent:** ChatGPT 5 Mini (Copilot)  
**Score:** ⭐⭐☆☆☆ (2 / 5)

- Useful for generating generic boilerplate error handling patterns.
- Struggled with context-aware error mapping for custom domain errors and database-specific failures.
- Occasionally suggested incorrect HTTP status codes (e.g., using `Unauthorized` instead of `InternalServerError`).
- Required manual debugging and corrections to align error responses with actual application behavior.

---

### 3. Frontend Error Handling  
**AI Agent:** Gemini  
**Score:** ⭐⭐⭐⭐⭐ (5 / 5)

- Performed very well with React-based error handling flows.
- Correctly handled API error states, user-facing messages, and redirects (e.g., 404 pages).
- Integrated cleanly with Formik and Yup validation patterns.
- Required minimal cleanup before production use.
- Significantly improved development speed and UX consistency.

---

### 4. Regex Generation & Validation  
**AI Agent:** ChatGPT 5 Mini (Copilot)  
**Score:** ⭐⭐⭐⭐⭐ (5 / 5)

- Accurately inferred intent from function and variable names.
- Generated correct and readable regex patterns for phone numbers, email addresses, and other validations.

---

### 5. SQL Query Assistance  
**AI Agent:** ChatGPT 5 Mini (Copilot)  
**Score:** ⭐⭐⭐☆☆ (3.5 / 5)

- Generated structurally correct SQL queries.
- Overlooked Go-specific implementation details, particularly pointer semantics required when scanning query results with `pgx` library, which could've led to empty result sets if left uncorrected.
- Required manual fixes to ensure correct data retrieval.

---

### 6. Frontend Styling & UI Design  
**AI Agent:** Gemini  
**Score:** ⭐⭐⭐⭐⭐ (5 / 5)

- Extremely effective for modern UI layouts and Material-UI theming.
- Enabled rapid iteration on visual design and component styling.
- Major productivity boost for frontend polish and consistency.

---

### 7. JWT Setup  
**AI Agent:** Gemini  
**Score:** ⭐⭐⭐☆☆ (3.5 / 5)

- Provided a functional JWT setup.
- Claims structure was somewhat non-standard. Referencing the [official Echo](https://echo.labstack.com/docs/cookbook/jwt) documentation provided a more straightforward implementation. 
---

### 8. Generic UI Templates & Components  
**AI Agent:** Gemini  
**Score:** ⭐⭐⭐⭐⭐ (5 / 5)

- Excellent for generating reusable React components and layout scaffolding.
- Enabled rapid prototyping with clean, extendable UI components.

---

### 9. General Code Autocompletion  
**AI Agent:** ChatGPT (Copilot)  
**Score:** ⭐⭐⭐⭐☆ (4 / 5)

- Very effective for local context autocompletion within functions or files.
- Occasionally missed broader context involving external libraries or imports.

---

### 10. README & Documentation Assistance  
**AI Agent:** ChatGPT  
**Score:** ⭐⭐⭐⭐⭐ (5 / 5)

- Assisted in structuring and refining the project README.
- Helped clearly document architecture, setup instructions, API contracts, and development decisions.
- Improved readability of project documentation.


[Go]: https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white
[Go-url]: https://go.dev/
[React.js]: https://img.shields.io/badge/React-20232A?style=for-the-badge&logo=react&logoColor=61DAFB
[React-url]: https://reactjs.org/
[psql]: https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white
[psql-url]: https://www.postgresql.org/
[docker]: https://img.shields.io/badge/Docker-2CA5E0?style=for-the-badge&logo=docker&logoColor=white
[docker-url]: https://www.docker.com/
