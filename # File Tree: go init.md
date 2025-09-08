# File Tree: go init

Generated on: 08/09/2025, 13:11:11
Root path: `/Users/rou/Documents/Github/go init`

```
├── .git/ 🚫 (auto-hidden)
├── .github/ 🚫 (auto-hidden)
├── cli-generator/
│   ├── .DS_Store 🚫 (auto-hidden)
│   ├── .gitignore 🚫 (auto-hidden)
│   ├── LICENSE
│   ├── Makefile
│   ├── README.md 🚫 (auto-hidden)
│   ├── go.mod 🚫 (auto-hidden)
│   ├── goinit
│   ├── goinit-generator
│   ├── install.sh 🚫 (auto-hidden)
│   └── main.go 🚫 (auto-hidden)
├── gin/
│   ├── api/
│   │   ├── common/
│   │   │   ├── dto/
│   │   │   │   └── auth_dto.go 🚫 (auto-hidden)
│   │   │   └── middleware/
│   │   │       ├── auth.go 🚫 (auto-hidden)
│   │   │       └── session.go 🚫 (auto-hidden)
│   │   └── protocol/
│   │       ├── http/
│   │       │   ├── handler/
│   │       │   │   ├── admin_handler.go 🚫 (auto-hidden)
│   │       │   │   ├── auth_handler.go 🚫 (auto-hidden)
│   │       │   │   ├── error_response.go 🚫 (auto-hidden)
│   │       │   │   ├── health.go 🚫 (auto-hidden)
│   │       │   │   ├── password_reset.go 🚫 (auto-hidden)
│   │       │   │   ├── user.go 🚫 (auto-hidden)
│   │       │   │   └── user_image.go 🚫 (auto-hidden)
│   │       │   ├── router/
│   │       │   │   └── router.go 🚫 (auto-hidden)
│   │       │   └── routes/
│   │       │       ├── auth_routes.go 🚫 (auto-hidden)
│   │       │       └── realtime_routes.go 🚫 (auto-hidden)
│   │       ├── sse/
│   │       │   └── handler.go 🚫 (auto-hidden)
│   │       └── ws/
│   │           └── handler.go 🚫 (auto-hidden)
│   ├── cmd/
│   │   └── api/
│   │       └── main.go 🚫 (auto-hidden)
│   ├── config/
│   │   └── env.go 🚫 (auto-hidden)
│   ├── docs/
│   │   ├── docs.go 🚫 (auto-hidden)
│   │   ├── swagger.json 🚫 (auto-hidden)
│   │   └── swagger.yaml 🚫 (auto-hidden)
│   ├── internal/
│   │   ├── app/
│   │   │   └── user/
│   │   │       ├── reset.go 🚫 (auto-hidden)
│   │   │       └── user_service.go 🚫 (auto-hidden)
│   │   ├── apperr/
│   │   │   └── errors.go 🚫 (auto-hidden)
│   │   ├── data/
│   │   │   └── user/
│   │   │       ├── model/
│   │   │       │   └── gorm/
│   │   │       │       └── user_gorm.go 🚫 (auto-hidden)
│   │   │       └── repo/
│   │   │           └── user_repo.go 🚫 (auto-hidden)
│   │   ├── db/
│   │   │   └── db.go 🚫 (auto-hidden)
│   │   ├── di/
│   │   │   ├── README.md 🚫 (auto-hidden)
│   │   │   ├── bench_test.go 🚫 (auto-hidden)
│   │   │   ├── container.go 🚫 (auto-hidden)
│   │   │   ├── di_container.go 🚫 (auto-hidden)
│   │   │   └── di_test.go 🚫 (auto-hidden)
│   │   ├── domain/
│   │   │   ├── model/
│   │   │   │   └── base.go 🚫 (auto-hidden)
│   │   │   └── user/
│   │   │       ├── model/
│   │   │       │   └── user.go 🚫 (auto-hidden)
│   │   │       └── repo/
│   │   │           └── user.go 🚫 (auto-hidden)
│   │   ├── lib/
│   │   │   ├── email/
│   │   │   │   ├── README_LOCAL_DEV.md 🚫 (auto-hidden)
│   │   │   │   ├── email_service.go 🚫 (auto-hidden)
│   │   │   │   └── local_email_service.go 🚫 (auto-hidden)
│   │   │   ├── id/
│   │   │   │   └── id.go 🚫 (auto-hidden)
│   │   │   ├── jwt/
│   │   │   │   ├── database_blacklist.go 🚫 (auto-hidden)
│   │   │   │   ├── jwt_service.go 🚫 (auto-hidden)
│   │   │   │   └── redis_blacklist.go 🚫 (auto-hidden)
│   │   │   ├── pwreset/
│   │   │   │   ├── database_service.go 🚫 (auto-hidden)
│   │   │   │   └── service.go 🚫 (auto-hidden)
│   │   │   └── storage/
│   │   │       ├── local.go 🚫 (auto-hidden)
│   │   │       ├── s3.go 🚫 (auto-hidden)
│   │   │       └── storage.go 🚫 (auto-hidden)
│   │   ├── logger/
│   │   │   └── logger.go 🚫 (auto-hidden)
│   │   └── server/
│   │       └── server.go 🚫 (auto-hidden)
│   ├── tmp/ 🚫 (auto-hidden)
│   ├── .air.toml 🚫 (auto-hidden)
│   ├── .dockerignore 🚫 (auto-hidden)
│   ├── .env 🚫 (auto-hidden)
│   ├── .env.example 🚫 (auto-hidden)
│   ├── .gitignore 🚫 (auto-hidden)
│   ├── Dockerfile
│   ├── ENV_CONFIG_EXAMPLE.md 🚫 (auto-hidden)
│   ├── Makefile
│   ├── docker-compose.yml 🚫 (auto-hidden)
│   ├── docker-env.example 🚫 (auto-hidden)
│   └── init-db.sql 🚫 (auto-hidden)
├── .DS_Store 🚫 (auto-hidden)
├── .gitignore 🚫 (auto-hidden)
├── README.md 🚫 (auto-hidden)
├── go.mod 🚫 (auto-hidden)
├── go.sum 🚫 (auto-hidden)
└── install.sh 🚫 (auto-hidden)
```

---

_Generated by FileTree Pro Extension_
