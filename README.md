# MiniOpsAgent

MiniOpsAgent 是一个面向运维学习和本地自动化实践的小型 AI Agent CLI。它运行在终端中，支持工作区文件操作、代码搜索、命令执行、联网检索、MCP 工具接入、Skill 加载、本地 RAG 代码索引、快照恢复和轻量 Runtime API。

这个项目的定位是“轻量终端运维 Agent”，不是完整 AIOps 平台，也不是企业级知识库系统。它适合作为学习项目：用较小的 Go 代码库观察一个 Agent 如何接入真实工具、如何做基础安全约束、如何暴露本地 Runtime API，以及如何把运维自动化能力组织成可读、可测试、可扩展的模块。

## 项目描述

MiniOpsAgent is a lightweight terminal-first AI operations agent. It connects an OpenAI-compatible LLM with local workspace tools, command execution, code search, web search, MCP tools, skills, snapshots, and a small runtime API for learning and automation experiments.

中文描述：

MiniOpsAgent 是一个面向运维场景的轻量级终端 AI Agent，支持本地工具调用、命令执行、RAG 检索、MCP 接入、快照恢复和 Runtime API，适合用于学习 AI Agent 与真实运维工具链的集成方式。

## 核心能力

- ReAct Agent 循环：模型可以多轮思考、调用工具、读取结果并生成最终回答。
- OpenAI-compatible LLM 客户端：默认支持 DeepSeek，也预置 OpenAI、GLM、Kimi、Step 等兼容配置。
- 终端 TUI：基于 Bubble Tea，提供全屏对话界面、流式输出、工具调用展示和状态栏。
- 单次命令模式：通过 `--once` 适配脚本、管道和自动化调用。
- 本地工具：读文件、写文件、列目录、glob、grep、执行命令、创建项目。
- 代码索引：本地构建轻量 RAG 索引，支持代码检索和简单关系图。
- Web 工具：支持 Web Search 和网页抓取，并对内网地址做基础拦截。
- Skill 机制：扫描用户级和项目级 `SKILL.md`，按需加载完整能力说明。
- MCP 接入：读取用户级和项目级 MCP 配置，动态注册远端工具。
- Runtime API：提供线程、turn 和 events，用于把本地 Agent 接到其他系统。
- 安全基础：PathGuard、CommandGuard、危险工具审计、turn 前后快照。

## 项目边界

MiniOpsAgent 保持小型项目定位，刻意不包含以下内容：

- 不提供完整 Web 管理后台。
- 不提供告警接入、审批流、RBAC、工单和可视化 Runbook 平台。
- 不直接内置 Kubernetes、Prometheus、Loki 等生产连接器。
- 不做企业文档知识库、权限隔离和大规模向量数据库。
- 不承诺生产级自动修复能力，高风险命令仍需要人工判断。

如果你需要完整 AIOps 平台，可以把 MiniOpsAgent 看成“本地 Agent/Worker 原型”；如果你需要企业知识库，可以把它看成“命令行侧的轻量代码检索能力”。

## 环境要求

- Go 1.26 或更新版本。
- 推荐安装 `rg`，用于更快的本地代码搜索。
- 可选：MCP server 所需运行时，例如 Node.js、Python 或远端 HTTP MCP 服务。
- 可选：`SERPAPI_API_KEY` 或 `SEARXNG_BASE_URL`，用于更稳定的联网搜索。

当前模块路径：

```go
module github.com/SuperWangYU-8088/MiniOpsAgent
```

## 快速开始

```bash
git clone https://github.com/SuperWangYU-8088/MiniOpsAgent.git
cd MiniOpsAgent

go run ./cmd/miniopsagent doctor
go run ./cmd/miniopsagent --once "介绍一下当前项目"
go run ./cmd/miniopsagent --plain
go run ./cmd/miniopsagent
```

常用命令：

```bash
go run ./cmd/miniopsagent version
go run ./cmd/miniopsagent doctor
go run ./cmd/miniopsagent index .
go run ./cmd/miniopsagent search "Agent"
go run ./cmd/miniopsagent graph
go run ./cmd/miniopsagent serve --port 8080
go run ./cmd/miniopsagent snapshot list
```

## 配置方式

MiniOpsAgent 会按以下顺序读取配置：

1. 内置 provider 默认值。
2. 当前目录 `.env`。
3. 用户 home 目录 `.env`。
4. 用户状态目录配置：`~/.miniopsagent/config.json`。
5. 环境变量覆盖。

推荐最小配置：

```bash
export MINIOPS_PROVIDER=deepseek
export MINIOPS_API_KEY=你的模型密钥
export MINIOPS_MODEL=deepseek-chat
```

也可以使用 provider 专属变量：

```bash
export DEEPSEEK_API_KEY=...
export OPENAI_API_KEY=...
export GLM_API_KEY=...
export KIMI_API_KEY=...
export STEP_API_KEY=...
```

联网搜索：

```bash
export SERPAPI_API_KEY=...
export SEARXNG_BASE_URL=http://localhost:8080
```

Runtime API：

```bash
export MINIOPS_RUNTIME_API_KEY=change-me
go run ./cmd/miniopsagent serve --port 8080
```

## 状态目录

用户级状态默认保存在：

```text
~/.miniopsagent/
```

主要内容：

```text
config.json                    模型与 Web 配置
mcp.json                       用户级 MCP 配置
skills.json                    Skill 禁用状态
skills/                        用户级 Skill
rag/code-index.json            本地代码索引
memory/long_term_memory.json   长期记忆
snapshots/                     turn 前后快照
audit/YYYY-MM-DD.jsonl         危险工具审计日志
```

项目级 MCP 和 Skill 可以放在：

```text
.miniopsagent/mcp.json
.miniopsagent/skills/<skill-name>/SKILL.md
```

## 内置工具

Agent 当前可用工具：

| 工具 | 说明 |
| --- | --- |
| `read_file` | 读取工作区内 UTF-8 文本文件 |
| `write_file` | 写入工作区内文件，单次最多 5 MB |
| `list_dir` | 列出目录 |
| `glob_files` | 使用 glob 查找文件 |
| `grep_code` | 精确搜索代码文本，优先使用 `rg` |
| `execute_command` | 在工作区执行命令，Windows 使用 PowerShell，Linux/macOS 使用 `sh -lc` |
| `create_project` | 在工作区创建简单项目目录 |
| `web_search` | 通过 SearXNG、SerpAPI 或 DuckDuckGo fallback 搜索 |
| `web_fetch` | 抓取公开网页并提取正文，阻止本地和内网地址 |
| `save_memory` | 用户明确要求时保存长期记忆 |
| `load_skill` | 加载完整 Skill 内容 |
| `search_code` | 查询本地代码索引 |
| `revert_turn` | 恢复某个快照 |

## Runtime API

启动：

```bash
export MINIOPS_RUNTIME_API_KEY=dev-secret
go run ./cmd/miniopsagent serve --port 8080
```

创建 thread：

```bash
curl -X POST http://127.0.0.1:8080/v1/threads \
  -H "Authorization: Bearer dev-secret"
```

创建 turn：

```bash
curl -X POST http://127.0.0.1:8080/v1/threads/<thread-id>/turns \
  -H "Authorization: Bearer dev-secret" \
  -H "Content-Type: application/json" \
  -d '{"input":"检查当前项目结构","mode":"react"}'
```

查看事件：

```bash
curl http://127.0.0.1:8080/v1/threads/<thread-id>/events \
  -H "Authorization: Bearer dev-secret"
```

## 目录结构

```text
cmd/miniopsagent/        CLI 入口和命令注册
internal/agent/          ReAct、Plan、Team 模式和记忆注入
internal/config/         配置加载、provider 默认值和环境变量覆盖
internal/llm/            OpenAI-compatible Chat/Stream 客户端
internal/tools/          文件、命令、Web、MCP、RAG、Skill、Snapshot 工具
internal/rag/            本地代码索引、词频检索和 Go 符号关系
internal/skill/          Skill 扫描、frontmatter 解析和延迟注入
internal/snapshot/       工作区快照创建与恢复
internal/runtime/        本地 HTTP Runtime API
internal/tui/            Bubble Tea 全屏终端界面
docs/                    中文教程、设计和实现文档
```

## 中文文档

- [入门教程](docs/入门教程.md)
- [使用说明](docs/使用说明.md)
- [架构设计](docs/架构设计.md)
- [实现说明](docs/实现说明.md)

## 安全说明

MiniOpsAgent 是学习型和本地自动化项目，不建议直接作为生产自动修复系统使用。

- 写文件、执行命令、恢复快照和 MCP 工具调用会写审计日志。
- 文件工具会限制在工作区内，避免直接访问工作区外路径。
- 命令工具会拒绝一批明显危险命令，例如 `sudo`、`rm -rf /`、`curl | sh`。
- `web_fetch` 会拦截 `localhost`、内网 IP、loopback、link-local 等非公网目标。
- 快照恢复是基础实现，适合学习和小项目，不等同于 Git 级别完整回滚。

## 后续可扩展方向

- Runtime API 持久化：使用 SQLite 保存 threads、turns 和 events。
- 真实 SSE：让 `/events` 持续推送新事件，而不是只返回已有事件。
- 命令审批：对高风险命令增加人工确认或策略配置。
- 运行日志：增加 run id 和结构化事件日志。
- RAG 升级：接入 embedding provider 和 SQLite/向量数据库。
- 安装脚本：增加 Windows/Linux 打包和服务托管脚本。
