Project Scope and Goal

Objective
Build a real-time chat platform where men and people who like men can connect, talk, and support each other through shared spaces.

Core Model
- Users join spaces
- Spaces contain channels
- Channels contain messages
- Access is controlled by membership and roles

v1 Goal
Deliver a working, Mint-hosted chat system with:
- account creation and login
- space creation and membership
- public and private channels
- message posting and reading
- live message updates

Access Model
- Public channels can be read without login
- Posting requires an authenticated user
- Private channels require membership
- Moderation actions require roles (owner/admin)

Channel Visibility
- private
- public unlisted (link access)
- public discoverable (future)

Messaging Behavior
- Messages ordered by created_at + message_id
- Load most recent messages first
- Load older messages via cursor pagination
- New messages appear in real time
- Soft delete supported in storage

Moderation Philosophy
- Users control what they see
- Provide tools to avoid users/content
- Keep platform-level enforcement minimal (except hard safety boundaries)
- Community-level moderation via roles

v1 Includes
- user signup and login
- create and join spaces
- create public and private channels
- send and read messages
- live message updates
- basic moderation delete

v1 Excludes
- voice/video
- file uploads
- reactions/likes
- threads
- bots
- federation
- self-hosting
- advanced moderation tools

Non-Functional Goals
- Fast reads for message retrieval
- Simple, understandable architecture
- Single deployable backend service
- Designed to scale later without redesign

Technology Stack
- Backend: Go + Gin
- Frontend: React (web)
- Database: MongoDB
- Runtime: Linux Mint (local host for v1)	

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

Repository / Structure
- Create root directory: ~/manumental-effort
- Use a top-level repo structure based on:
  - client
  - server
  - shared
- Put frontend and backend code under those areas as needed.
- Include architecture and product docs in the repo from the start.

Documentation
- Add and maintain these core docs:
  - docs/v1_scope.md
  - docs/domain_model.md
  - docs/system_architecture.md
  - docs/api_principles.md

Code / Naming Conventions
- File and path names are spelled out and separated by hyphens where appropriate.
- Prefer clear names over short names.
- Keep framework code separate from business logic and persistence logic.
- Repo-wide paths may use hyphens, but Go package/file naming should follow standard Go conventions.

Architecture Principle
- Middleware belongs to the server layer, not shared.
- shared/ contains contracts and schemas, not runtime logic.
- HTTP middleware is implemented under server/internal/platform/.
