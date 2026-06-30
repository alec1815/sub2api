# sub2api 开发环境排错手册

> **目标**：从零搭建开发环境 → 成功运行前后端 → 登录中转站系统
>
> **全程覆盖**：记录每一步遇到的问题、排查思路、尝试过的方案、最终解决方案、涉及的关键代码和参考资料。

---

## 目录结构

| 文件 | 阶段 | 说明 |
|------|------|------|
| [00-overview.md](./00-overview.md) | 总览 | 全流程时序图 + 问题清单索引 + 经验总结 |
| [01-env-prep.md](./01-env-prep.md) | ① 环境准备 | PostgreSQL + Redis + Node.js/pnpm |
| [02-database-setup.md](./02-database-setup.md) | ② 数据库初始化 | psql 命令行陷阱、密码重置、migration 执行 |
| [03-backend-startup.md](./03-backend-startup.md) | ③ 后端首次启动 | config.yaml 陷阱、NeedsSetup 逻辑、启动流程分析 |
| [04-admin-account.md](./04-admin-account.md) | ④ 管理员创建 | AUTO_SETUP 流程、Windows DATA_DIR 问题、辅助脚本方案 |
| [05-login-verification.md](./05-login-verification.md) | ⑤ 登录验证 | 数据库确认、health check、前端登录流程 |
| [06-project-architecture.md](./06-project-architecture.md) | ⑥ 架构分析 | 基于官方博客的 7 步调用链路、核心设计、源码阅读路线 |
| [07-deepseek-platform-analysis.md](./07-deepseek-platform-analysis.md) | ⑦ 平台添加分析 | DeepSeek 平台添加：从 UI 到网关的 5 层 14 处改动全链路源码追踪 |
| [appendix.md](./appendix.md) | 附录 | 关键代码解读、GetDataDir 路径优先级、启动决策树 |
| [08-ent-repo-pitfalls.md](./08-ent-repo-pitfalls.md) | ⑧ P1/P2 企业开发排错 | Ent Schema 命名、部分唯一索引、自引用 edge、OrderOption、分层架构、Entity 映射、跨 entity 查询 |
| [09-ent-service-pitfalls.md](./09-ent-service-pitfalls.md) | ⑨ P3 Service 排错 | 接口耦合、结构体字段同步、类型匹配、未使用函数、变量重定义、双重校验、月度用量占位 |
| [10-enterprise-service-analysis.md](./10-enterprise-service-analysis.md) | ⑩ P3 Service 分析 | 6 服务架构模式、业务逻辑流程、依赖图、设计对照、TODO 追踪 |
| [11-p4-middleware-billing-analysis.md](./11-p4-middleware-billing-analysis.md) | ⑪ P4 中间件+计费路由 | 企业认证中间件、计费分流器、三池隔离门、网关计费管道集成 |
| [12-p5-handler-route-wire-analysis.md](./12-p5-handler-route-wire-analysis.md) | ⑫ P5 Handler+Route+Wire | 8 Handler 文件、37 端点、双重鉴权体系、Wire 注入链路、4 缺失路由修复 |
| [13-p7-frontend-development.md](./13-p7-frontend-development.md) | ⑬ P7 前端开发记录 | 7.1-7.4c 前端类型/API/路由/视图等 |
| [14-p7-comprehensive-analysis.md](./14-p7-comprehensive-analysis.md) | ⑭ P7 综合分析 | PRD符合度 / 前后端字段一致性 / P6+P8测试可行性 / 日志统计企业维度 |

---

## 快速定位问题

| 现象 | 可能原因 | 参考文档 |
|------|----------|----------|
| `psql: connection refused` | PostgreSQL 服务未启动 | [01-env-prep.md](./01-env-prep.md) |
| `FATAL: password authentication failed` | 密码错误 / pg_hba.conf 配置 | [02-database-setup.md](./02-database-setup.md) |
| psql 连 `localhost` 很慢 | IPv6 优先 → 回退 IPv4 | [02-database-setup.md](./02-database-setup.md) |
| 启动后端提示 "First run" | config.yaml 不存在（正常） | [03-backend-startup.md](./03-backend-startup.md) |
| 登录后无管理员账号 | 手动创建了 config.yaml，跳过了 Install | [04-admin-account.md](./04-admin-account.md) |
| `Access is denied: /app/data/.installed` | Windows 未设 DATA_DIR | [04-admin-account.md](./04-admin-account.md) |
| pnpm install 报 EPERM | npm 残留的 node_modules 冲突 | [01-env-prep.md](./01-env-prep.md) |
| bcrypt hash 在 psql 中被截断 | PowerShell 的 `$` 转义 | [02-database-setup.md](./02-database-setup.md) |
| `go run` 输出无法捕获 | Go 子进程不退出的特性 | [04-admin-account.md](./04-admin-account.md) |

---

## 阅读建议

1. **新成员首次搭建**：按 01 → 02 → 03 → 04 → 05 顺序阅读
2. **只遇到某个具体问题**：看 README 的快速定位表
3. **想了解排查方法论**：看 [00-overview.md](./00-overview.md) 的经验总结

---

| `go generate ./ent` 生成代码后编译失败 | Schema 字段名变更未同步三层 | [08-ent-repo-pitfalls.md](./08-ent-repo-pitfalls.md) |
| Repository 接口 import 循环依赖 | 接口放 repo 包 → service 层 import | [08-ent-repo-pitfalls.md](./08-ent-repo-pitfalls.md) |
| Ent 排序不生效 | OrderOption 类型使用错误 | [08-ent-repo-pitfalls.md](./08-ent-repo-pitfalls.md) |
| APIKey struct 企业字段缺失 | Schema 加了字段但 Service 层未同步 | [09-ent-service-pitfalls.md](./09-ent-service-pitfalls.md) |
| 企业 Profile 月度用量全是 0 | 网关计费层 P4 未实现，TODO 占位 | [09-ent-service-pitfalls.md](./09-ent-service-pitfalls.md) |
| 分页计算 int→int64 类型错误 | PaginationParams 返回 int 但分页切分需要 int64 | [09-ent-service-pitfalls.md](./09-ent-service-pitfalls.md) |
| 某 Service 方法不知道有哪些依赖 | 查依赖关系图和接口定义位置 | [10-enterprise-service-analysis.md](./10-enterprise-service-analysis.md) |
| 设计文档 VS 实际代码不一致 | 逐项对照表 §五 | [10-enterprise-service-analysis.md](./10-enterprise-service-analysis.md) |
| 企业 Key 扣了个人余额 / 反过来 | 隔离门未生效 → 检查 PoolType 传递链 | [11-p4-middleware-billing-analysis.md](./11-p4-middleware-billing-analysis.md) |
| 企业中间件 403 | 检查 JWT → enterprise_members → enterprise status 链 | [11-p4-middleware-billing-analysis.md](./11-p4-middleware-billing-analysis.md) |

---

> **最后更新**: 2026-06-29 · 基于企业功能 P1-P5 开发过程记录
