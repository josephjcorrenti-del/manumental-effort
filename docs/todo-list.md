todo list

phase 0 - local environment setup
[x] install Go
[x] verify Go installation
[x] install Node.js + npm
[x] verify Node.js + npm installation
[x] install MongoDB locally
[x] verify MongoDB is running
[x] install Git if needed

phase 1 - project initialization
[x] create root directory: ~/manumental-effort
[x] create base directories:
        - docs/
        - client/web/
        - server/
        - shared/
        - deploy/
        - scripts/
[x] initialize git repository
[x] create README.md
[x] create docs/decisions.md
[x] copy current decisions into docs/decisions.md

phase 2 - core documentation
[x] create docs/v1-scope.md
[x] define v1 includes / excludes
[~] define what v1 excludes
[x] create docs/domain-model.md
[~] define core objects (good enough for v1, will evolve)
[ ] define core objects:
        - user
        - space
        - channel
        - message
        - membership
[x] create docs/system-architecture.md
[~] define architecture (correct, not exhaustive)
        - client / server / shared structure
        - REST + WebSocket model
        - MongoDB usage
        - Mint-hosted deployment
[x] create docs/api-principles.md
[~] define API rules (good baseline, will expand)
        - auth model
        - pagination model
        - error format
        - naming conventions

phase 2.5 - naming + repo conventions
[x] record repo naming conventions in docs/decisions.md
[x] record Go package/file naming conventions in docs/api-principles.md
[x] ensure repo-wide paths use hyphens where appropriate
[ ] ensure Go code follows normal Go naming conventions
[~] repo/docs naming conventions are aligned enough to proceed

phase 3 - backend bootstrap
[x] initialize Go module under server/
[x] add Gin dependency
[x] create API entrypoint under server/cmd/api/
[x] implement minimal HTTP server startup
[x] add /health endpoint
[x] verify server runs locally
[x] verify /health returns success

phase 4 - MongoDB integration
[x] add MongoDB Go driver
[x] create MongoDB connection layer
[x] create config file: server/configs/app-local.yaml
[x] load config at startup
[x] initialize MongoDB connection at startup
[x] verify backend connects successfully to MongoDB

phase 5 - users domain
[x] create server/internal/users/
[x] define user model
[x] implement user repository
[x] implement create user flow
[x] implement get user by id flow
[x] add POST /users
[x] add GET /users/{id}
[x] persist users in MongoDB
[x] verify users can be created and retrieved

phase 5.5 - users validation + uniqueness cleanup
[x] add basic create-user validation
[x] validate required fields:
        - username
        - display_name
        - email
[x] normalize user input where appropriate
[x] reject invalid Mongo user ids cleanly
[x] add uniqueness checks for:
        - username
        - email
[x] decide app-vs-db ownership for uniqueness enforcement
        - DB owns actual uniqueness
        - app owns friendly error handling
[x] add MongoDB indexes for:
        - username_normalized unique
        - email_normalized unique
[x] return clearer user-facing errors for duplicate username/email
[x] verify invalid input is rejected
[x] verify duplicate username is rejected
[x] verify duplicate email is rejected

phase 6 - authentication
[x] define minimal authentication approach

step 1 - credentials + password hashing
[x] add auth configuration:
        - jwt signing key
        - token expiry
[x] create server/internal/auth/
[x] define credential model
[x] implement password hashing
[x] implement credential repository
[x] update create user flow to store password credentials
[x] validate password on create user
[x] verify user + credential are both created

step 2 - login + jwt token issuance
[x] define login request/response models
[x] implement credential lookup by normalized email
[x] implement login service flow
[x] implement JWT token creation
[x] add POST /auth/login
[x] verify valid login returns token
[x] verify invalid email/password is rejected

step 3 - auth middleware + protected route
[x] implement JWT token parsing/verification
[x] add auth middleware
[x] attach authenticated user id to request context
[x] add protected test route
[x] verify authenticated requests work
[x] verify missing token is rejected
[x] verify invalid token is rejected
[x] make /auth/me return the current user profile

step 4 - auth usability + hardening
[x] create helper to safely get authenticated user id from context
[x] standardize context key usage (avoid raw strings)
[ ] optionally fetch user once in middleware (future-ready, not required now)
[x] define rule: which routes require auth vs public
[ ] apply middleware to at least one non-auth route (prove reuse)
[x] refactor /auth/me to use auth helper + users domain cleanly

phase 7 - spaces + channels + memberships
[x] create server/internal/spaces/
[x] define space model
[x] implement create space flow
[x] add POST /spaces
[x] persist spaces in MongoDB

[x] create server/internal/memberships/
[x] define membership model
[x] implement membership repository
[x] add unique index for:
        - space_id + user_id
[x] create owner membership when a space is created
[x] update create space flow to also create owner membership
[x] verify creating a space also creates owner membership

[x] implement join space flow
[x] add POST /spaces/{id}/join

[x] create server/internal/channels/
[x] define channel model
[x] implement channel repository
[x] implement create channel flow
[x] add POST /spaces/{id}/channels
[x] implement list channels in a space flow
[x] add GET /spaces/{id}/channels
[x] enforce membership-based access rules
[x] verify user can create and navigate a space

phase 8 - messages over REST
[x] create server/internal/messages/
[x] define message model
[x] implement send message flow
[x] implement fetch recent messages flow
[x] implement fetch older messages flow using cursor pagination
[x] enforce canonical ordering:
        - created_at ascending
        - tie-break by message_id
[x] verify messages can be posted and read through REST

phase 9 - realtime over WebSocket
[ ] create server/internal/realtime/
[ ] implement WebSocket endpoint
[ ] implement connection lifecycle handling
[ ] implement subscribe / unsubscribe flow for channel watching
[ ] enforce access checks on channel subscription
[ ] broadcast new messages to active channel subscribers
[ ] verify live messages appear in connected clients

phase 10 - frontend bootstrap
[ ] initialize React app under client/web/
[ ] create minimal application shell
[ ] implement login view
[ ] implement space list view
[ ] implement channel list view
[ ] implement message view
[ ] connect frontend to REST API
[ ] connect frontend to WebSocket updates
[ ] verify end-to-end user flow works locally

phase 11 - core UX + moderation basics
[ ] implement soft delete in message storage model
[ ] implement delete message flow
[ ] implement basic moderation delete capability
[ ] implement scroll-up loading for older messages
[ ] add logging
[ ] improve error handling
[ ] verify basic chat UX is stable

phase 12 - future-ready scaffolding
[ ] separate visibility from discoverability in the model
[ ] support:
        - private
        - public unlisted
        - public discoverable
[ ] leave room for channel allow-lists and space-level bans
[ ] leave room for user follows:
        - follow spaces
        - follow channels
        - follow users
[ ] leave room for future message metadata:
        - edits
        - likes / reactions
        - ranking / activity

to be prioritized 
    mongodb audit records:
        - created_by
        - updated_by
        - record_start_date
        - record_end_date

    username rules / identity
        - preserve display casing for username
        - enforce uniqueness on lowercase-normalized username
        - add username_normalized field
        - consider email_normalized field
        - validate allowed username characters and length
        - return friendly duplicate/validation errors

    init / rollback / install scripts
        - add local Linux init scripts
        - add local Linux rollback scripts
        - add local Linux install/setup scripts
        - plan AWS init scripts for later deployment phase
        - plan AWS rollback scripts for later deployment phase
        - decide ownership between deploy/ and scripts/
