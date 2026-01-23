# Kuro VCS

**Kuro** is a local-first version control system built with **Go** and **SQLite**, focused on correctness, explicit state, and simple internal invariants.

Instead of scattering repository state across files and directories, Kuro stores all core metadata in a single SQLite database, making branching, snapshots, and history explicit and inspectable.

> *Kuro (黒) means “black” in Japanese — representing a clean slate and a deliberate starting point.*

---

## Why Kuro?

Git is powerful, but it relies heavily on implicit behavior, loose files, and historical complexity.

Kuro explores a different design space:

* **Explicit state over implicit magic**
* **SQLite-backed storage instead of ad-hoc files**
* **Clear separation of refs, snapshots, and configuration**
* **Local-first correctness before synchronization**

Kuro is not trying to “replace Git overnight”.
It is an experiment in building a **simpler, more transparent VCS core**.

---

## Design Philosophy

* **Branches are refs** — names that may or may not point to a snapshot
* **Commits are snapshots** — immutable state captured explicitly
* **HEAD never dangles** — it always points to a ref
* **Unborn branches are valid** — no fake commits, no hacks
* **Local correctness comes first** — syncing is layered on later

If something exists in the repository, it should be visible, queryable, and understandable.

---

## Architecture

Kuro is structured as a set of layered components:

* **Core**

  * SQLite-based storage engine
  * Refs, snapshots, objects, and config
  * No UI, no logging, no side effects

* **CLI**

  * Cobra-based command interface
  * Human-friendly output
  * Thin orchestration over core logic

* **Local API (planned)**

  * HTTP API for local repository access
  * Enables tooling and UI integration

* **Remote API (planned)**

  * Synchronization service
  * Push / pull semantics built on explicit state

* **Web UI (planned)**

  * Next.js-based interface
  * Repository inspection and management via local API

---

## Current Features

* Repository initialization
* SQLite-backed repository storage
* Branch creation, listing, deletion
* Explicit HEAD and ref management
* Cross-platform CLI (Windows, Linux, macOS)

---

## Planned Features

* File tracking and commits
* Snapshot history and diffing
* Checkout and detached HEAD
* Local HTTP API
* Remote synchronization (push / pull)
* Web-based repository browser
* Collaboration workflows

---

## Getting Started

### Prerequisites

* Go **1.19+**

### Build

```bash
git clone <repository-url>
cd kuro

cd cli
go build -o kuro
```

### Initialize a Repository

```bash
./kuro init
```

This creates a `.kuro/` directory containing the SQLite database and repository metadata.

---

## Branch Management

```bash
# List branches
./kuro branch list

# Create a branch
./kuro branch add dev

# Delete a branch
./kuro branch delete dev
```

Example output:

```
→ main
• dev
• feature-x
```

---

## Project Structure

```
kuro/
├── cli/           # Command-line interface
├── core/          # Core VCS logic and database
├── local-api/     # Local HTTP API (planned)
├── remote-api/    # Remote sync service (planned)
├── web/           # Web interface (planned)
├── go.work
└── .gitignore
```

---

## Development Notes

* `.kuro/` is intentionally ignored by Git
* The core package contains **no UI or logging**
* CLI output is handled explicitly via a UI layer
* SQLite schema models VCS concepts directly (refs, snapshots, objects)

---

## Status

Kuro is under **active development**.

The current focus is on:

* solid core abstractions
* clean branch and ref semantics
* building a foundation that can evolve safely

Breaking changes are expected at this stage.

---

## Contributing

Contributions, feedback, and design discussion are welcome.

If you’re interested in:

* VCS internals
* database-backed systems
* explicit state machines
* building tools from first principles

you’ll feel at home here.

---

## License

[MIT](./LICENSE)

---