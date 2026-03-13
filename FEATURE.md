# Feature Design

## Multi-Device Sync

This document is about future sync architecture only.

Current local architecture:
- the kernel is a local Go daemon
- local clients talk to it over a Unix socket
- this document does not change that local IPC model

Future sync is a separate concern layered on top of the local kernel.

### Problem

Crona is local-first with a SQLite database per device. Multiple devices create divergent state. Direct DB sync isn't viable — SQLite files aren't merge-friendly, and there's no natural conflict resolution at the row level.

### Why Op Logs Are the Right Foundation

Every mutation in Crona already produces an immutable `Op` record:

```
{ id, user_id, device_id, entity, entity_id, action, payload, timestamp }
```

This is an append-only log of all state changes, ordered by timestamp, tagged by device. Syncing Crona across devices reduces to: **share the op log, replay it on each device, derive consistent state**.

This is structurally similar to how distributed databases (CRDTs, event sourcing) handle sync — no central authority required.

---

### Proposed Solution: Layered Sync

Three sync modes, each builds on the previous:

#### Layer 1 — File-Based Sync (No Server Required)

Export the op log as an append-only NDJSON file to a shared folder watched by any cloud sync provider:

```
~/<cloud-provider>/crona-sync/
  ops-<device-id>.ndjson       # This device's ops, append-only
  ops-<device-id-2>.ndjson     # Other device's ops, written by them
```

On startup (and periodically), each device:
1. Reads all `ops-*.ndjson` files from other devices
2. Filters ops it hasn't seen yet (by `id` or `timestamp + device_id`)
3. Replays them against the local DB in timestamp order
4. Appends its own new ops to `ops-<its-device-id>.ndjson`

**Works with**: iCloud Drive, Dropbox, Google Drive, Syncthing, any folder sync tool.
**No server needed.** Sync latency = cloud sync latency (typically seconds).

#### Layer 2 — Self-Hosted Sync Relay (Optional)

A lightweight relay server that acts as a dumb op store — no business logic, no Crona domain knowledge.

```
POST /ops          # Push a batch of ops
GET  /ops?since=   # Pull ops since timestamp, optionally filtered by device
```

Devices push their ops on mutation and pull on startup or reconnect. The relay is stateless beyond storing op records.

**Deployment**: Single Docker container, runs on any VPS, Raspberry Pi, or home server. No managed infrastructure.

```yaml
# docker-compose.yml
services:
  crona-relay:
    image: crona/relay
    ports: ["3001:3001"]
    volumes: ["./data:/data"]   # ops stored as SQLite or flat files
```

Kernel sync config could add:

```json
{
  "sync": {
    "mode": "relay",
    "url": "http://my-server:3001",
    "token": "<shared-secret>"
  }
}
```

#### Layer 3 — Local Network P2P (Same Network, No Cloud)

When devices are on the same network, sync directly without any relay:

- Kernel advertises itself via mDNS (`_crona._tcp.local`)
- Devices discover peers on the local network
- Push and pull ops directly over a peer transport

Falls back to relay or file sync when off-network.

---

### Conflict Resolution

Since ops are immutable and timestamped, conflicts are resolved deterministically without coordination:

| Entity | Strategy | Rationale |
|--------|----------|-----------|
| Repo / Stream / Issue fields | Last-write-wins by `timestamp` | Simple fields, low conflict risk |
| Issue status transitions | Ordered by timestamp | Status is a state machine; latest wins |
| Active context | Per device, never merged | Each device has independent context |
| Sessions | Device-scoped, no merge needed | Sessions have `device_id`, owned by one device |
| Stash | Namespaced by device | Stashes are device-local by design |
| Deletes | Soft-delete wins over any prior update | `deleted_at` is terminal |

**Key invariant**: Sessions and active context are device-scoped. Two devices working simultaneously never conflict on timer state — each device owns its own session independently.

The only real conflict surface is **shared entity fields** (issue title, repo name) edited on two offline devices. Last-write-wins by timestamp is acceptable here — these are low-frequency mutations, not concurrent collaborative edits.

---

### What Doesn't Sync

- Scratchpad file **contents** (only metadata syncs via ops)
  - Scratchpad files themselves should go in the shared cloud folder alongside op files, or stay device-local
- Active context (always device-local)
- Timer state (derived, not stored — each device recomputes from its own sessions)
- Auth tokens (per-kernel, per-device)

---

### Implementation Notes

- **Op replay must be idempotent** — replaying an op that's already applied is a no-op (check by `op.id`)
- **Clock skew**: Use the op's `timestamp` as-is from the originating device. Minor skew (seconds) is acceptable for LWW. For correctness, consider Lamport timestamps if needed later.
- **Initial sync**: On first connect, pull all ops from relay/files, replay from scratch. The local DB is fully derivable from the op log — this is the key invariant to maintain.
- **Op log completeness**: Every mutation must produce an op. Gaps in the log = sync gaps. Audit this before implementing sync.
