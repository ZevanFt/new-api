# new-api 运维归档（talenting.vip）

更新时间：2026-03-08

## 1. 部署现状

- 服务域名：`https://newapi.talenting.vip`
- 运行容器：`new-api`（镜像 `calciumion/new-api:latest`）
- 数据持久化：`/root/data/new-api:/data`
- SQLite 文件：`/root/data/new-api/one-api.db`

## 2. 资源与性能配置（已生效）

- 容器限制：
  - `--memory=512m`
  - `--memory-swap=512m`
  - `--cpus=0.75`
  - `--pids-limit=256`
  - `--log-opt max-size=10m`
  - `--log-opt max-file=3`
- 关键环境变量：
  - `SYNC_FREQUENCY=600`
  - `GOMAXPROCS=1`
  - `GOMEMLIMIT=384MiB`
  - `GOGC=90`
  - `RELAY_MAX_IDLE_CONNS=60`
  - `RELAY_MAX_IDLE_CONNS_PER_HOST=10`
  - `MAX_REQUEST_BODY_MB=32`
  - `TASK_QUERY_LIMIT=200`
  - `FORCE_STREAM_OPTION=false`
  - `CHANNEL_UPSTREAM_MODEL_UPDATE_TASK_ENABLED=false`
  - `UPDATE_TASK=false`

## 3. 渠道与兼容策略

- 渠道 1/2 已启用并配置 Codex 兼容覆盖：
  - `header_override`: `{"OpenAI-Beta":"responses=experimental","originator":"codex_cli_rs"}`
  - `param_override`: `{"instructions":"","store":false}`
- 全局 chat->responses 策略（options）：
  - `global.chat_completions_to_responses_policy={"enabled":true,"all_channels":true,"model_patterns":["(?i)codex","(?i)^gpt-5(\\.|-|$)"]}`

## 4. 高峰保护（管理员豁免）

- 模型请求限流已启用：
  - `ModelRequestRateLimitEnabled=true`
  - `ModelRequestRateLimitDurationMinutes=1`
  - `ModelRequestRateLimitCount=40`
  - `ModelRequestRateLimitSuccessCount=30`
  - `ModelRequestRateLimitGroup={"default":[40,30],"admin":[2147483647,2147483647]}`
- 管理员豁免：
  - 用户 `zevan` 已切换到 `admin` 组
  - `zevan` 的 token 已切换到 `admin` 组
  - 渠道 `1/2/3` 分组已扩展为 `default,admin`

## 5. 低负载后台任务配置

- 看板导出：
  - `DataExportEnabled=false`
  - `DataExportInterval=30`
- 自动渠道测试：
  - `monitor_setting.auto_test_channel_enabled=false`
- 上游模型自动同步任务：
  - `CHANNEL_UPSTREAM_MODEL_UPDATE_TASK_ENABLED=false`

## 6. 定时任务

- 渠道 1 保活（每 4 小时）：
  - `15 */4 * * * /root/scripts/ch1-keepalive.sh >> /root/data/new-api/logs/ch1-keepalive.log 2>&1`
- 服务巡检（每 15 分钟）：
  - `*/15 * * * * /root/scripts/new-api-healthcheck.sh >> /root/data/new-api/logs/new-api-healthcheck.log 2>&1`

## 7. 脚本与日志

- 脚本：
  - `/root/scripts/ch1-keepalive.sh`
  - `/root/scripts/new-api-healthcheck.sh`
  - `/root/scripts/new-api-recreate-current.sh`
- 配置：
  - `/root/.config/ch1-keepalive.env`
- 日志：
  - `/root/data/new-api/logs/ch1-keepalive.log`
  - `/root/data/new-api/logs/new-api-healthcheck.log`

## 8. 备份与回滚

- SQLite 备份：
  - `/root/backups/new-api/one-api-20260308-214521-pre-peak-tune.db`
- 环境快照：
  - `/root/backups/new-api/new-api-env-20260308-214521.txt`
- 一键重建：
  - 执行 `/root/scripts/new-api-recreate-current.sh`

## 9. 快速验证命令

```bash
docker ps | grep new-api
docker stats --no-stream | grep new-api
curl -I https://newapi.talenting.vip/
tail -n 50 /root/data/new-api/logs/ch1-keepalive.log
tail -n 50 /root/data/new-api/logs/new-api-healthcheck.log
```
