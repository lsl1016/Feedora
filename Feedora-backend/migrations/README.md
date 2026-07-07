# migrations

阶段一使用 GORM AutoMigrate 自动建表（见 `internal/model/all.go`），无需手动执行 SQL。

阶段二如需版本化迁移，可在此放置 `NNN_*.sql`（golang-migrate 格式），例如 outbox、通知、积分、幂等等表。
