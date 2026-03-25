~/manumental-effort/
├── README.md
├── .gitignore
├── docs/
│   ├── v1-scope.md
│   ├── domain-model.md
│   ├── system-architecture.md
│   ├── api-principles.md
│   └── decisions.md
├── client/
│   └── web/
│       ├── package.json
│       ├── public/
│       │   └── index.html
│       └── src/
│           ├── app/
│           ├── components/
│           ├── features/
│           │   ├── auth/
│           │   ├── spaces/
│           │   ├── channels/
│           │   ├── messages/
│           │   └── users/
│           ├── lib/
│           ├── pages/
│           ├── services/
│           └── main.jsx
├── server/
│   ├── go.mod
│   ├── go.sum
│   ├── cmd/
│   │   └── api/
│   │       └── main.go
│   ├── internal/
│   │   ├── auth/
│   │   ├── users/
│   │   ├── spaces/
│   │   ├── channels/
│   │   ├── messages/
│   │   ├── memberships/
│   │   ├── follows/
│   │   ├── moderation/
│   │   ├── realtime/
│   │   ├── platform/
│   │   │   ├── config/
│   │   │   ├── http/
│   │   │   ├── mongodb/
│   │   │   ├── logging/
│   │   │   └── ids/
│   │   └── common/
│   ├── configs/
│   │   ├── app-example.yaml
│   │   └── app-local.yaml
│   ├── scripts/
│   │   ├── start-local.sh
│   │   └── test-local.sh
│   └── tests/
│       ├── integration/
│       └── e2e/
├── shared/
│   ├── api/
│   │   ├── openapi.yaml
│   │   └── websocket-events.md
│   ├── schemas/
│   │   ├── message.json
│   │   ├── channel.json
│   │   ├── space.json
│   │   └── user.json
│   └── glossary/
│       └── product-terms.md
├── deploy/
│   ├── systemd/
│   │   └── manumental-effort-api.service
│   ├── nginx/
│   │   └── manumental-effort.conf
│   └── aws/
│       └── README.md
└── scripts/
    ├── init-repo.sh
    └── backup-local.sh
