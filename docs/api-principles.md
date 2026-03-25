API Principles

Auth
- token-based authentication
- authenticated user attached to request context

Pagination
- cursor-based pagination
- fetch most recent messages first
- fetch older messages using a before-cursor

Ordering
- messages ordered by:
        - created_at ascending
        - message_id as tie-breaker

Errors
- use consistent error response format
- return clear client-facing messages
- avoid leaking internal details

Naming
- use clear, descriptive names
- repo-wide paths may use hyphens
- Go package/file naming follows normal Go conventions

Middleware
- middleware handles cross-cutting HTTP concerns only
- business logic, permissions, and persistence stay out of middleware
