#!/usr/bin/env bash
# 应用 migrations/ 目录下尚未执行的 NNN_*.sql 版本化迁移脚本。
#
# 用法（在 Feedora-backend 目录下执行）:
#   ./scripts/apply-migrations.sh [脚本目录，默认 migrations]
#
# 已应用记录写入 schema_migrations 表（version 主键），每个脚本只执行一次。
# 默认通过 docker compose 的 feedora-mysql 容器执行；连本机 MySQL 时改 MYSQL_CMD 为
# 形如 "mysql -uroot -proot feedora" 的客户端命令即可。
#
# 说明：不在 Go 代码里执行脚本文件内容（安全扫描要求 SQL 与输入严格分离），
# 迁移执行统一收敛在此脚本；AutoMigrate 仍由 go run ./cmd/migrate 负责。

set -euo pipefail

DIR="${1:-migrations}"
MYSQL_CMD=(docker exec -i feedora-mysql mysql -uroot -proot feedora)

# 文件名字符集校验：通过者不可能包含引号/分号等破坏 SQL 的字符。
valid_name() {
  case "$1" in
    [0-9][0-9][0-9]_[a-z0-9._-]*.sql) return 0 ;;
    *) return 1 ;;
  esac
}

"${MYSQL_CMD[@]}" -e "CREATE TABLE IF NOT EXISTS schema_migrations (version varchar(191) PRIMARY KEY, applied_at datetime(3) NOT NULL)"

applied_tmp=$(mktemp)
trap 'rm -f "$applied_tmp"' EXIT
"${MYSQL_CMD[@]}" -N -B -e "SELECT version FROM schema_migrations" > "$applied_tmp"

count=0
for f in "$DIR"/[0-9][0-9][0-9]_*.sql; do
  [ -e "$f" ] || continue
  name=$(basename "$f")
  valid_name "$name" || { echo "跳过（命名不符合 NNN_描述.sql）: $name"; continue; }
  if grep -Fxq "$name" "$applied_tmp"; then continue; fi

  echo "应用迁移: $name"
  "${MYSQL_CMD[@]}" < "$f"
  "${MYSQL_CMD[@]}" -e "INSERT INTO schema_migrations (version, applied_at) VALUES ('$name', NOW(3))"
  echo "$name" >> "$applied_tmp"
  count=$((count + 1))
done

echo "完成：本次应用 $count 个迁移。已应用版本："
"${MYSQL_CMD[@]}" -N -B -e "SELECT version FROM schema_migrations ORDER BY version"
