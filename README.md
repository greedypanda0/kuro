# Kuro VCS

**Kuro** is a local-first version control system built in **Go** with **SQLite**.  
This repository currently focuses on the **core engine** and the **CLI**. The API is under development.

> *Kuro (黒) means “black” in Japanese — a clean slate and deliberate starting point.*

---

## What Kuro Is (Today)

Kuro explores a **transparent, explicit** VCS design:

- Repository state is stored in a single SQLite database
- Refs, snapshots, and config are **first-class** concepts
- The CLI is a thin orchestration layer over the core

---

## Core Concepts

- **Refs**: branch names that point to snapshots (or remain unborn)
- **Snapshots**: immutable commits captured as explicit records
- **Objects**: content-addressed blobs stored in SQLite
- **HEAD**: always points to a ref (never a detached orphan)

---

## Project Structure

```
kuro/
├── core/          # Core VCS engine (SQLite-backed)
├── cli/           # Command-line interface
├── api/           # API layer (in development)
├── go.work
└── README.md
```

---

## Features (Core + CLI)

- Initialize a repository
- SQLite-backed storage (refs, snapshots, objects)
- Branch create / list / delete
- Add & stage files
- Commit snapshots
- Checkout refs or snapshots (workspace reset with `--ws`)
- Status & logs
- Raw SQL queries against the repo database (`sql`)
- Config and auth management
- Remote management and push
- Identity check (`whoami`)

---

## Getting Started

### Prerequisites
- Go **1.19+**

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

## CLI Examples

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

### Logs
```
./kuro logs
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

---

## CLI Command Reference

- `init` — initialize a repository
- `add <path>` — stage a file or directory
- `remove <path>` — unstage a file or directory
- `remove .` — clear the entire stage
- `status` — show current branch and last commit
- `status --stage` — list staged files
- `commit -m "<message>"` — create a snapshot from staged files
- `logs` — show commit history for HEAD
- `logs --branch <name>` — show logs for a specific branch
- `branch list` — list branches
- `branch create <name>` — create a branch
- `branch delete <name>` — delete a branch
- `checkout [branch|commit]` — move HEAD to a branch or snapshot
- `checkout [branch|commit] --ws` — reset workspace to a snapshot
- `config --name "<name>"` — set local user name
- `config --token "<token>"` — set auth token
- `remote` — show current remote
- `remote add <user>/<repo>` — set remote
- `remote remove` — remove remote
- `push` — push the local database to the remote
- `whoami` — show local user and verify token remotely
- `sql "<query>"` — run raw SQL against the repository database

## Development Status

The **core engine** and **CLI** are active and usable.  
The **API** is under development and intentionally out of scope for this README.

---

## License

[MIT](./LICENSE)