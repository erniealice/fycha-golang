# Services Reorganization — Progress Log

**Plan:** [plan.md](./plan.md)
**Started:** 2026-03-08
**Branch:** `dev/20260308-fycha-services-reorg`

---

## Phase 1: Create `services/docprocessor/` — NOT STARTED

- [ ] Create `services/docprocessor/document_service.go` (move from root, change package to `docprocessor`)
- [ ] Delete `document_service.go` from root
- [ ] Verify `go build ./...` in fycha-golang-ryta

---

## Phase 2: Create `services/storage/` + update consumers — NOT STARTED

- [ ] Create `services/storage/storage_handler.go` (move from root, change package to `storage`)
- [ ] Delete `storage_handler.go` from root
- [ ] Update `apps/retail-client/internal/composition/container.go` imports
- [ ] Update `apps/retail-client/internal/composition/views.go` imports
- [ ] Update `apps/service-client/internal/composition/container.go` imports
- [ ] Update `apps/service-client/internal/composition/views.go` imports
- [ ] Verify `go build ./...` in retail-client and service-client

---

## Phase 3: Move `StorageImagesPrefix` constant — NOT STARTED

- [ ] Move `StorageImagesPrefix` from `routes.go` to `services/storage/storage_handler.go`
- [ ] Update any consumer references to `fycha.StorageImagesPrefix`
- [ ] Verify `go build ./...`

---

## Summary

- **Phases complete:** 0 / 3
- **Files modified:** 0 / 9

---

## Skipped / Deferred

| Item | Reason |
|------|--------|
| `datasource.go` move | Thin interface, domain port not service logic — defer to future |

---

## How to Resume

To continue this work:
1. Read this progress file and the [plan](./plan.md)
2. Check git status for any uncommitted changes
3. Start from the first incomplete phase above
4. Update checkboxes and summary as you complete steps

Key files to read first:
- `packages/fycha-golang-ryta/document_service.go` (source for Phase 1)
- `packages/fycha-golang-ryta/storage_handler.go` (source for Phase 2)
- `apps/retail-client/internal/composition/container.go` (consumer to update)
- `apps/service-client/internal/composition/container.go` (consumer to update)
