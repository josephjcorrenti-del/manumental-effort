Domain Model

User
- represents a person using the platform
- can authenticate
- can join spaces
- can post messages
- can follow spaces, channels, and users later

Space
- top-level community container
- has visibility/discoverability settings
- contains channels
- contains memberships

Channel
- belongs to one space
- can be public or private
- contains messages

Message
- belongs to one channel
- created by one user
- ordered by created_at + message_id
- supports soft delete in storage

Membership
- connects a user to a space
- defines role:
        - owner
        - admin
        - member
- controls posting and moderation access
