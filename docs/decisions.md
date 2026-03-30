Decisions

Product Direction
- Build manumental-effort as a centralized Discord-like chat platform first.
- Self-hosting is a later phase, not a v1 requirement.
- Use the product term "space" rather than "server".
- Design so later versions can support user-hosted servers/instances.

Hosting / Deployment
- v1 runs from the Linux Mint box.
- Design v1 so it can later move to AWS without major structural changes.
- Keep deployment simple for v1.
- Defer serverless and containerization decisions until after the Mint-hosted v1 is working.

Core Stack
- Use Go for the backend.
- Use Gin as the API/web framework for v1.
- Add guard rails so Gin handlers stay thin.
- Use React for the frontend.
- Use MongoDB for persistence in v1.
- Optimize around fast reads and clean query patterns.

Messaging Architecture
- Use REST for standard CRUD, auth, setup, and posting messages.
- Use WebSocket for live message delivery and live channel updates.
- Start with Option A: send messages over REST, receive broadcasts over WebSocket.
- Leave room for richer WebSocket-first messaging later if needed.

Access Model
- Public spaces/channels can be read without login.
- Posting requires an authenticated account and appropriate membership.
- Private channels require authentication and membership.
- Moderation/admin actions require role-based permissions.
- Separate visibility from discoverability:
  - private
  - public unlisted
  - public discoverable

Permissions / Moderation Direction
- Start with simple space roles:
  - owner
  - admin
  - member
- Start with simple channel visibility:
  - public
  - private by membership/allow-list
- Plan for future blacklist/ban support.
- Product principle: users should have tools to avoid content, words, or users they do not want to see.
- Prefer user filtering and community controls over heavy centralized language policing, except where hard safety boundaries are required.

Message Behavior
- Default channel display order is chronological: oldest to newest.
- Canonical ordering rule is:
  - created_at ascending
  - tie-break by message_id
- Load the most recent X messages first.
- Load older messages by cursor-based pagination when the user scrolls upward.
- New live messages appear in real time as they are received.
- Support delete in v1.
- Prefer soft delete in storage with clean UX in presentation.
- Editing should be in v1 if cheap; otherwise make it one of the first enhancements.
- Leave room for future message metadata such as likes/reactions and ranking.

User Experience Direction
- Users will later be able to follow spaces, channels, and other users.
- Followed items should appear before generally active content.
- Later phases may include activity views, search, and richer discovery.

Architecture Principles
- Use middleware for cross-cutting HTTP concerns only.
- Keep business logic, permissions, and persistence out of middleware.
- Middleware belongs to the server layer, not shared.
- shared/ contains contracts and schemas, not runtime logic.
- HTTP middleware is implemented under server/internal/platform/.

Repository / Structure
- Use a top-level repo structure based on:
  - client
  - server
  - shared
- Put frontend and backend code under those areas as needed.
- Include architecture and product docs in the repo from the start.

Documentation
- Add and maintain these core docs:
  - docs/v1-scope.md
  - docs/domain-model.md
  - docs/system-architecture.md
  - docs/api-principles.md
  - docs/decisions.md

Code / Naming Conventions
- File and path names are spelled out and separated by hyphens where appropriate.
- Prefer clear names over short names.
- Keep framework code separate from business logic and persistence logic.
- Repo-wide paths may use hyphens, but Go package/file naming should follow standard Go conventions.

Design Principle
- Build the simplest working system first, then evolve it.
- Avoid over-engineering in v1.
- Preserve flexibility for future expansion.

scripting principle
- add scripts after manual steps are understood and stable
- do not automate unstable setup too early
- prefer small, explicit scripts over large opaque setup scripts

phase 6 auth decisions
- use email + password login
- store password hashes only
- use bcrypt for password hashing
- use JWT bearer tokens for API authentication
- keep auth credentials separate from user profile documents
- normalize email for auth lookup
- attach authenticated user id to Gin request context

Auth rules

Public
- POST /users
- POST /auth/login
- GET /users/{id} (current behavior; revisit later if needed)

Authenticated
- GET /auth/me

Future
- all create/update/delete routes require authentication by default
- write routes should use authenticated user id as actor identity
