## Context

The codebase currently uses "package" terminology for what is essentially a Docker Compose role/configuration. This creates confusion with traditional software package managers. The rename to "role" better reflects the semantic meaning in deployment orchestration.

Current state:
- Database table `packages` stores Docker Compose content
- API endpoint `/packages` exposes CRUD operations
- TUI has "Packages" menu item
- All code references use "package" terminology

## Goals / Non-Goals

**Goals:**
- Rename all code, API, and database references from "package" to "role"
- Maintain existing functionality (no behavior changes)
- Update OpenAPI spec and regenerate types
- Update tests to reflect new terminology

**Non-Goals:**
- Add new functionality
- Migrate existing data (clean slate, no production data)
- Update documentation outside code (README, etc.)

## Decisions

### 1. Rename strategy: Full rename vs Alias

**Decision**: Full rename (not alias)

**Rationale**: 
- Simpler codebase with single terminology
- No technical debt from maintaining both names
- Implementation already complete, tests pass

### 2. Database migration approach

**Decision**: Rewrite migration file (not incremental ALTER)

**Rationale**:
- No production data exists (per user confirmation)
- Simpler than writing reversible migration
- The migration file is fresh (created for this project)

### 3. API endpoint path

**Decision**: `/roles` (not `/compose` or similar)

**Rationale**:
- Simple, intuitive naming
- Matches database table name
- Consistent with REST conventions

## Risks / Trade-offs

- **[Risk]** API breaking change
  - **Mitigation**: Document in proposal, clear error messages if old endpoint accessed
  
- **[Risk]** Stale references in documentation
  - **Mitigation**: Documentation updates out of scope for this change

## Migration Plan

Since this is a clean implementation (no production data):

1. Rewrite migration file with new table/column names
2. Rename SQL query files and update references
3. Update OpenAPI spec with new paths/schemas
4. Run code generation (`make generate-api`, `make generate-sql`)
5. Rename and update handler files
6. Rename and update TUI model files
7. Rename and update test files
8. Run tests to verify

## Open Questions

None - implementation complete and verified.
