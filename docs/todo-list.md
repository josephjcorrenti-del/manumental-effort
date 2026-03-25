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
[~] define what v1 includes
[~] define what v1 excludes
[x] create docs/domain-model.md
[ ] define core objects:
        - user
        - space
        - channel
        - message
        - membership
[ ] create docs/system-architecture.md
[ ] define:
        - client / server / shared structure
        - REST + WebSocket model
        - MongoDB usage
        - Mint-hosted deployment
[~] create docs/api-principles.md
[ ] define:
        - auth model
        - pagination model
        - error format
        - naming conventions

phase 2.5 - naming + repo conventions
[ ] record repo naming conventions in docs/decisions.md
[ ] record Go package/file naming conventions in docs/api-principles.md
[ ] ensure repo-wide paths use hyphens where appropriate
[ ] ensure Go code follows normal Go naming conventions

phase 3 - backend bootstrap
[ ] initialize Go module under server/
[ ] add Gin dependency
[ ] create API entrypoint under server/cmd/api/
[ ] implement minimal HTTP server startup
[ ] add /health endpoint
[ ] verify server runs locally
[ ] verify /health returns success

phase 4 - MongoDB integration
[ ] add MongoDB Go driver
[ ] create MongoDB connection layer
[ ] create config file: server/configs/app-local.yaml
[ ] load config at startup
[ ] initialize MongoDB connection at startup
[ ] verify backend connects successfully to MongoDB

phase 5 - users domain
[ ] create server/internal/users/
[ ] define user model
[ ] implement create user flow
[ ] implement get user by id flow
[ ] add POST /users
[ ] add GET /users/{id}
[ ] persist users in MongoDB
[ ] verify users can be created and retrieved

phase 6 - authentication
[ ] define minimal authentication approach
[ ] implement login flow
[ ] implement token-based authentication
[ ] add auth middleware
[ ] attach authenticated user context to requests
[ ] protect routes where needed
[ ] verify authenticated requests work

phase 7 - spaces + channels + memberships
[ ] create server/internal/spaces/
[ ] create server/internal/channels/
[ ] create server/internal/memberships/
[ ] define space model
[ ] define channel model
[ ] define membership model
[ ] implement create space flow
[ ] implement join space flow
[ ] implement create channel flow
[ ] implement list channels in a space flow
[ ] enforce membership-based access rules
[ ] verify user can create and navigate a space

phase 8 - messages over REST
[ ] create server/internal/messages/
[ ] define message model
[ ] implement send message flow
[ ] implement fetch recent messages flow
[ ] implement fetch older messages flow using cursor pagination
[ ] enforce canonical ordering:
        - created_at ascending
        - tie-break by message_id
[ ] verify messages can be posted and read through REST

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
