# Jungle RPG

A single-player browser-based jungle exploration RPG built with Go and React.

## Stack

- **Backend:** Go, chi router, sqlc, modernc SQLite (pure Go, no CGO)
- **Frontend:** React + TypeScript + Vite, Tailwind CSS
- **Auth:** Google OAuth2, server-side sessions (gorilla/sessions)
- **Deploy:** Fly.io with persistent volume for SQLite

## Gameplay

Explore a 15x15 jungle grid, discover locations, complete quests, and earn enough gold to purchase the village (1500g) to win.

**Controls:**
- Arrow keys — move
- E — interact (enter village/cavern, rest, weave baskets)
- W — exit village/cavern, enter village
- C — build camp (50g, 2 uses)

**Locations:** Village, Capybara Den, Ruins, Waterfall, Cavern

**Hazards:** Beehive (-25 HP), Jaguar (-65 HP), Anaconda (HP → 5), Quicksand (-45 Energy)

**Village NPCs:**
- Village Elder — assigns quests (find locations for gold rewards)
- Shopkeeper — Health Potions (50g, +20 HP) and Stimulants (25g, +10 Energy)
- Basket Weaver — earn +6 Gold per interaction
- Guest Hut — rest (+30 Energy, +20 HP)
- Frontiersman — buy the village for 1500g to win

**Cavern:** A 3x8 dungeon with Deep Holes, Hidden Treasure (+400g), and Hidden Springs (+75 HP, +100 Energy).

## Development

### Prerequisites

- Go 1.23+
- Node.js 20+
- [sqlc](https://sqlc.dev/) (optional, generated code is committed)

### Environment Variables

```
GOOGLE_CLIENT_ID=...
GOOGLE_CLIENT_SECRET=...
GOOGLE_REDIRECT_URL=http://localhost:5173/auth/google/callback
SESSION_KEY=<random-32-byte-hex>
DATABASE_PATH=./jungle.db
PORT=8080
```

### Run locally

```bash
# Install frontend dependencies
cd web && npm install && cd ..

# Start both servers (Vite dev server proxies API to Go)
make dev
```

Frontend runs on `http://localhost:5173`, Go backend on `http://localhost:8080`.

### Build for production

```bash
make build
```

### Regenerate sqlc

```bash
make sqlc
```

## Deploy to Fly.io

```bash
fly launch
fly secrets set GOOGLE_CLIENT_ID=... GOOGLE_CLIENT_SECRET=... SESSION_KEY=... GOOGLE_REDIRECT_URL=https://<app>.fly.dev/auth/google/callback
fly volumes create jungle_data --size 1 --region ord
fly deploy
```

## Project Structure

```
├── cmd/server/main.go          # Entry point
├── internal/
│   ├── auth/                   # Google OAuth2 + session middleware
│   ├── game/                   # Pure game logic (no HTTP deps)
│   ├── api/                    # HTTP handlers for game and saves
│   ├── repository/             # SQLite setup + sqlc generated code
│   └── server/                 # Router and server setup
├── web/src/
│   ├── api/                    # API client layer
│   ├── types/                  # TypeScript types
│   ├── components/             # MapGrid, StatsBar, Inventory, etc.
│   └── pages/GamePage.tsx      # Main game page
├── Dockerfile                  # Multi-stage build
├── fly.toml                    # Fly.io config
└── sqlc.yaml                   # sqlc config
```
