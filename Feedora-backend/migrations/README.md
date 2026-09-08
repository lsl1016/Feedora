# migrations

版本化表结构变更 SQL 脚本目录，与 GORM AutoMigrate 协同工作：

- **AutoMigrate**（`go run ./cmd/migrate`）：按 Go 模型自动建新表、加新列，负责"从零到有"。
- **本目录**：存放需要**显式留痕、可审查、可回溯**的变更脚本（索引、约束、数据回填、列改名等 AutoMigrate
  做不了或做不好的操作），负责"从版本 N 到版本 N+1"。

## 约定

1. 文件名 `NNN_描述.sql`（三位递增编号 + 下划线 + 小写字母/数字/`._-` 字符），编号即执行顺序，
   已应用的编号不复用。
2. 执行方式：`./scripts/apply-migrations.sh`（在 Feedora-backend 目录下运行）按编号顺序执行**未应用**
   的脚本，并把文件名记入 `schema_migrations` 表（version 主键），每个脚本只执行一次。
   默认走 docker 的 `feedora-mysql` 容器，连本机 MySQL 时改脚本里的 `MYSQL_CMD`。
3. 脚本内容为普通 SQL，多条语句用 `;` 分隔；不要写含内嵌分号的过程体（存储过程 / 触发器）。
4. 脚本应保持单一职责、幂等优先（能用 `IF NOT EXISTS` 的地方就用）。
5. `baseline_schema.sql` 是当前全量表结构的参考快照（由 `mysqldump --no-data` 生成），**不编号、不参与执行**，
   仅用于 diff 审查；重大变更后可重新生成。

## 常用命令

```bash
./scripts/apply-migrations.sh                  # 应用未执行的 SQL 脚本
docker exec feedora-mysql mysql -uroot -proot feedora \
  -e "SELECT * FROM schema_migrations ORDER BY version"   # 查看已应用版本
docker exec feedora-mysql sh -c \
  'mysqldump -uroot -proot --no-data --skip-comments --compact feedora' \
  | sed 's/ AUTO_INCREMENT=[0-9]*//' > migrations/baseline_schema.sql  # 重新生成基线
```
