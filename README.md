# Kuro VCS

**Kuro** is a local-first version control system built in **Go** with **SQLite**.  
This repository focuses on the **core engine** and **CLI**, with a **remote API** under development.

> *Kuro (黒) means “black” in Japanese — a clean slate and deliberate starting point.*

---

## Index

1. [Overview](#overview)
2. [Project Structure](#project-structure)
3. [Core Concepts](#core-concepts)
4. [Features](#features)
5. [Getting Started](#getting-started)
6. [CLI Usage](#cli-usage)
7. [Configuration & Storage](#configuration--storage)
8. [Remote API](#remote-api)
9. [Development Notes](#development-notes)
10. [License](#license)

---

## Overview

Kuro explores a transparent, explicit VCS design:

- Repository state is stored in a single SQLite database
- Refs, snapshots, and config are **first-class** concepts
- The CLI is a thin orchestration layer over the core engine

---

## Project Structure

```
kuro/
├── core/          # Core VCS engine (SQLite-backed)
├── cli/           # Command-line interface
├── api/remote/    # Remote API server
├── go.work
└── README.md
```

---

## Core Concepts

- **Refs**: branch names that point to snapshots (or remain unborn)
- **Snapshots**: immutable commits captured as explicit records
- **Objects**: content-addressed blobs stored in SQLite
- **HEAD**: always points to a ref (never a detached orphan)

---

## Features

- Initialize a repository
- SQLite-backed storage (refs, snapshots, objects)
- Branch create / list / delete
- Add & stage files
- Commit snapshots
- Checkout refs or snapshots (workspace reset with `--ws`)
- Status & logs
- Diff for staged files (`diff`)
- Raw SQL queries against the repo database (`sql`)
- Config and auth management
- Remote management and push
- Identity check (`whoami`)

---

## Getting Started

### Prerequisites
- Go **1.19+** (or your target Go version)

### Build the CLI
```
git clone <repository-url>
cd kuro/cli
go build -o kuro
```

### Initialize a Repo
```
./kuro init
```

This creates a `.kuro/` directory containing the SQLite database and metadata.

---

## CLI Usage

### Stage Files
```
./kuro add .
```

### Commit
```
./kuro commit -m "Initial snapshot"
```

### Status
```
./kuro status --stage
```
With `--stage`, Kuro lists staged files and unstaged files (prefixed with `-`).

### Diff
```
./kuro diff
./kuro diff -f path/to/file
```

### Logs
```
./kuro logs
./kuro logs --branch main
```

### Branches
```
./kuro branch list
./kuro branch create dev
./kuro branch delete dev
```

### Checkout
- Switch HEAD only (no workspace changes):
```
./kuro checkout dev
```

- Reset workspace to snapshot:
```
./kuro checkout dev --ws
```

### Raw SQL
```
./kuro sql "SELECT name, snapshot_hash FROM refs"
```

---

## Configuration & Storage

### Repository Layout
- `.kuro/kuro.db` — SQLite database
- `.kuro/.kuroignore` — ignore rules

### User Config
Stored at:
- `~/.kuro/config.json`

Fields:
- `name` — author name used in commits
- `token` — auth token used for remote API

Set via:
```
./kuro config --name "Your Name"
./kuro config --token "your-token"
```

### `.kuroignore`
Default entries created on `init`:
```
.kuro
.git
node_modules
dist
build
```

---

## Remote API

The remote API is a Gin server backed by PostgreSQL.

### Environment

See `.env.example` in:
- `/` (root)
- `api/remote/`
- `cli/`
- `core/`

The remote API uses:
- `DATABASE_URL` (required)
- `REMOTE_HTTP_ADDR`
- `REMOTE_HTTP_SHUTDOWN_TIMEOUT`
- `REMOTE_LOG_LEVEL`
- `REMOTE_LOG_DEV`

### Auth

All `/api/*` endpoints are protected by auth middleware:

- `Authorization: Bearer <token>`  
**or**
- session cookies (`authjs.session-token` / `__Secure-authjs.session-token`)

### Endpoints

Base: `http://localhost:8080/api`

#### Utility
- `GET /health` — health check
- `GET /version` — build info
- `GET /ping` — returns auth user context

#### Repositories
- `GET /repositories?name=&user_id=` — list repositories
- `GET /repositories/:id` — get repository by id
- `POST /repositories` — upload repo database

**Upload requirements:**
- `Content-Type: application/octet-stream`
- `Authorization: Bearer <token>`
- `X-Remote: <user>/<repo>`

#### Refs
- `GET /repositories/:id/refs` — list refs
- `GET /repositories/:id/refs/:ref` — get ref by name

#### Objects
- `GET /repositories/:id/objects` — list object hashes
- `GET /repositories/:id/objects/:hash` — get object content

#### Snapshots
- `GET /repositories/:id/snapshots`
- `GET /repositories/:id/snapshots/:snapshot_id`
- `GET /repositories/:id/snapshots/:snapshot_id/files`
- `GET /repositories/:id/snapshots/:snapshot_id/files/*file_id`

#### Users
- `GET /users/me` — current user
- `GET /users?name=&page=&limit=` — list users (paged)

---

## Development Notes

- This repo uses Go workspaces (`go.work`) for local module development.
- The CLI depends on the core module.
- The remote API is isolated under `api/remote`.

---

## License

[MIT](./LICENSE)