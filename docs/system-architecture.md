System Architecture

Structure
- client/
- server/
- shared/

Client
- React web frontend
- uses REST for standard operations
- uses WebSocket for live updates

Server
- Go backend using Gin
- single deployable service for v1
- exposes REST API
- exposes WebSocket endpoint
- owns middleware, business logic, permissions, persistence, and realtime delivery

Shared
- contracts, schemas, docs, and shared definitions
- not runtime logic
- not middleware

Database
- MongoDB
- optimized for fast reads
- stores users, spaces, channels, memberships, and messages

Realtime
- clients watch channels over WebSocket
- server broadcasts new messages to active channel watchers

Deployment
- v1 runs on Linux Mint
- single-node local deployment
- MongoDB local
- AWS is a later deployment target
