# Services Reorganization — Design Plan

**Date:** 2026-03-08
**Branch:** `dev/20260308-fycha-services-reorg`
**Status:** Draft
**Package:** packages/fycha-golang-ryta

---

## Overview

Move `document_service.go` and `storage_handler.go` from the fycha root package into `services/` sub-packages (`services/docprocessor/` and `services/storage/`), grouping all service-layer code under `services/` alongside the existing `services/doctemplate/`. Root package keeps domain types (routes, labels, filters, htmx helpers, assets). Views stay as-is.

---

## Motivation

The fycha root package currently mixes domain types (labels, routes, filters) with service-layer code (storage HTTP handler, document processing orchestrator). Moving service code into `services/` sub-packages:

1. **Clarifies package responsibility** — root = domain types & config, services/ = service-layer logic
2. **Consistent with doctemplate** — `services/doctemplate/` already lives under services/, the other two should too
3. **Cleaner imports** — consumer apps explicitly import the service they need rather than getting everything from one flat package
4. **Independent testability** — each service package can have focused tests without pulling in unrelated deps

---

## Current vs Target Structure

```
# Current                              # Target
fycha-golang-ryta/                     fycha-golang-ryta/
├── document_service.go     ──MOVE──►  ├── services/docprocessor/
│                                      │   └── document_service.go
├── storage_handler.go      ──MOVE──►  ├── services/storage/
│                                      │   └── storage_handler.go
├── services/doctemplate/   ──KEEP──►  ├── services/doctemplate/   (unchanged)
├── datasource.go           ──KEEP──►  ├── datasource.go
├── routes.go               ──KEEP──►  ├── routes.go
├── routes_config.go        ──KEEP──►  ├── routes_config.go
├── labels.go               ──KEEP──►  ├── labels.go
├── report_filter.go        ──KEEP──►  ├── report_filter.go
├── htmx.go                 ──KEEP──►  ├── htmx.go
├── assets.go               ──KEEP──►  ├── assets.go
├── package_dir.go          ──KEEP──►  ├── package_dir.go
├── views/                  ──KEEP──►  ├── views/  (unchanged)
```

---

## Symbols Being Moved

### To `services/docprocessor/` (package `docprocessor`)

| Symbol | Type | Current Consumer Apps |
|--------|------|---------------------|
| `DocumentService` | struct | **None** (newly created, no consumers yet) |
| `NewDocumentService` | func | **None** |
| `StorageReadWriter` | interface | **None** |

**Import change:** `fycha.DocumentService` → `docprocessor.DocumentService`
**Impact:** Zero — no consumer apps use these yet.

### To `services/storage/` (package `storage`)

| Symbol | Type | Current Consumer Apps |
|--------|------|---------------------|
| `StorageHandler` | struct | retail-client, service-client (views.go) |
| `NewStorageHandler` | func | retail-client, service-client (container.go) |
| `StorageReader` | interface | retail-client, service-client (container.go) |
| `StorageReadResult` | struct | retail-client, service-client (container.go) |
| `StorageRouteRegistrar` | interface | retail-client, service-client (container.go) |
| `ErrObjectNotFound` | error var | retail-client, service-client (container.go) |

**Import change:** `fycha.StorageHandler` → `storage.StorageHandler`
**Impact:** 4 files across 2 apps need import updates.

---

## Implementation Steps

### Phase 1: Create `services/docprocessor/` (zero-impact move)

1. Create `packages/fycha-golang-ryta/services/docprocessor/document_service.go`
   - Change package declaration to `package docprocessor`
   - Update internal import of doctemplate (path unchanged, just verifying)
2. Delete `packages/fycha-golang-ryta/document_service.go`
3. Verify: `go build ./...` in fycha-golang-ryta

### Phase 2: Create `services/storage/` + update consumers

1. Create `packages/fycha-golang-ryta/services/storage/storage_handler.go`
   - Change package declaration to `package storage`
   - Move all symbols: `StorageHandler`, `NewStorageHandler`, `StorageReader`, `StorageReadResult`, `StorageRouteRegistrar`, `ErrObjectNotFound`, `contentTypeFromExt`
2. Delete `packages/fycha-golang-ryta/storage_handler.go`
3. Update consumer apps (4 files):
   - `apps/retail-client/internal/composition/container.go` — add `storage` import, update type references
   - `apps/retail-client/internal/composition/views.go` — update `fycha.StorageHandler` → `storage.StorageHandler`
   - `apps/service-client/internal/composition/container.go` — add `storage` import, update type references
   - `apps/service-client/internal/composition/views.go` — update `fycha.StorageHandler` → `storage.StorageHandler`
4. Verify: `go build ./...` in each affected app

### Phase 3: Clean up `routes.go` constant

1. Move `StorageImagesPrefix` constant from `routes.go` to `services/storage/storage_handler.go`
   - This constant is storage-specific and belongs with the storage package
   - Check if any consumer references `fycha.StorageImagesPrefix` — if so, update
2. Verify: `go build ./...`

---

## File References

| File | Change | Phase |
|------|--------|-------|
| `packages/fycha-golang-ryta/services/docprocessor/document_service.go` | **New file** — moved from root | 1 |
| `packages/fycha-golang-ryta/document_service.go` | **Delete** | 1 |
| `packages/fycha-golang-ryta/services/storage/storage_handler.go` | **New file** — moved from root | 2 |
| `packages/fycha-golang-ryta/storage_handler.go` | **Delete** | 2 |
| `packages/fycha-golang-ryta/routes.go` | Remove `StorageImagesPrefix` constant | 3 |
| `apps/retail-client/internal/composition/container.go` | Update fycha storage imports | 2 |
| `apps/retail-client/internal/composition/views.go` | Update `StorageHandler` type ref | 2 |
| `apps/service-client/internal/composition/container.go` | Update fycha storage imports | 2 |
| `apps/service-client/internal/composition/views.go` | Update `StorageHandler` type ref | 2 |

---

## Context & Sub-Agent Strategy

**Estimated files to read:** ~15
**Estimated files to modify:** 7 (2 new, 2 delete, 4 update, 1 constant move)
**Estimated context usage:** Low (<30 files)

No sub-agents needed. Single session is sufficient. The moves are mechanical — change package declaration, update imports.

---

## Risk & Dependencies

| Risk | Impact | Mitigation |
|------|--------|------------|
| Consumer apps that use `fycha.StorageHandler` break | Medium — 4 files in 2 apps | Phase 2 updates all consumers in same commit |
| Views importing from root `fycha` package lose access to moved types | None — views use domain types (labels, routes, filters), not services |
| `datasource.go` left at root despite being service-adjacent | Low — it's an interface, not service logic | Could move to `services/reporting/` in future if needed |

**Dependencies:**
- Phase 1 is independent (no consumers)
- Phase 2 depends on Phase 1 being committed (or done in same session)
- Phase 3 is independent but logically follows Phase 2

---

## Acceptance Criteria

- [ ] `services/docprocessor/document_service.go` exists with package `docprocessor`
- [ ] `services/storage/storage_handler.go` exists with package `storage`
- [ ] Original root files (`document_service.go`, `storage_handler.go`) deleted
- [ ] `StorageImagesPrefix` moved to storage package
- [ ] `go build ./...` passes for `packages/fycha-golang-ryta`
- [ ] `go build ./...` passes for `apps/retail-client`
- [ ] `go build ./...` passes for `apps/service-client`
- [ ] No remaining references to `fycha.StorageHandler` or `fycha.DocumentService` in consumer apps
- [ ] Existing tests pass: `go test ./services/doctemplate/ -v`

---

## Design Decisions

**Why not move `datasource.go`?** It's a thin interface (3 methods) that defines the data contract for report views. It's more of a domain port than a service. Moving it to `services/reporting/` would create a package with just one interface file — not worth the churn. Can revisit when fycha gets more data access patterns.

**Why `docprocessor` not `document`?** "document" is too generic and could conflict with naming in other contexts. "docprocessor" clearly describes the responsibility: orchestrating document template processing with storage I/O. It parallels "doctemplate" (template engine) — one processes, one templates.

**Why keep root files?** The remaining root files (routes, labels, filters, htmx, assets) are domain configuration — they define what fycha is, not how it does things. Consumer apps import these for wiring (route constants, label structs, filter types). They belong at the package identity level.
