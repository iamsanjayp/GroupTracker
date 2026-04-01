# GroupTracker

Team Productivity & Project-Based Learning Web App for college environments.

## Tech Stack

- **Frontend:** React 18 + Vite + React Router + Zustand + Axios
- **Backend:** Go (Fiber v2)
- **Database:** MySQL 8.0
- **Auth:** JWT (access + refresh) + Google OAuth 2.0

## Getting Started

### 1. Database Setup

```bash
# In MySQL, run:
mysql -u root -p < backend/internal/database/migrations/001_init.sql
```

### 2. Backend

```bash
cd backend

# Update .env with your MySQL credentials and JWT secret

# Install dependencies
go mod tidy

# Run the server
go run cmd/server/main.go
```

Backend starts on `http://localhost:8080`

### 3. Frontend

```bash
cd frontend

# Install dependencies
npm install

# Start dev server
npm run dev
```

Frontend starts on `http://localhost:5173` and proxies API calls to the backend.

## Environment Variables

Copy `backend/.env` and update:

| Variable | Description |
|----------|-------------|
| `DB_USER` | MySQL username |
| `DB_PASSWORD` | MySQL password |
| `DB_NAME` | Database name (default: `grouptracker`) |
| `JWT_SECRET` | Long random string for signing tokens |
| `GOOGLE_CLIENT_ID` | Google OAuth client ID |
| `GOOGLE_CLIENT_SECRET` | Google OAuth client secret |

## Features (Phase 1)

- JWT + Google OAuth authentication
- Team creation & join via invite code
- Role-based access (Captain, VC, Manager, Strategist, Member)
- Dashboard with personal & admin views
- Daily 7-hour activity logger with auto-points
- Project tracker with tasks & member share %
- Points system with leaderboard & PS records
