# contracts

`contracts` 定义 FoundationX 跨域稳定契约——端口（接口）、事件协议和 DTO。它是域间通信的唯一合法通道，确保数据域、分析域、决策域和执行域之间的接口稳定、可演进。

## 定位

contracts 遵循 [xlib-standard](https://github.com/ZoneCNH/xlib-standard) 的治理协议，但不是标准源、不是 generator、不是模板仓库。contracts 不拥有传输实现，不绑定具体通信协议。

## 版本

- 当前发布版本：`v0.4.8`

## 验证

- 发布前执行 `GOWORK=off make docs-check`
- 发布前执行 `GOWORK=off make dependency-check`
- 发布前执行 `GOWORK=off make standard-impact-check`
- 发布前执行 `GOWORK=off make release-check`

## 证据与下游

- 交付说明使用 `DONE with evidence:` 作为证据前缀。
- 发布清单与校验文件包含 `release/manifest/latest.json` 和 `release/manifest/latest.json.sha256`。
- 标准影响产物包含 `release/standard-impact/latest.md`。
- 依赖更新入口包含 `renovate.json` 和 `.github/dependabot.yml`。
- 下游同步信号使用 `downstream_sync_required`，相关约束见 `docs/downstream-sync-policy.md`。
- 分流决策字段 `downstream_release_decision`（只允许 `required` / `not_required`）。
- 仓库规则决策字段 `repository_rules_release_decision`（只允许 `audit_required` / `not_required`）。
- 回归约束包含 `FUZZ_SMOKE_TIME`。
- 主要下游仓库引用 `kernel`。

## 核心内容

- **端口接口**: MarketDataProvider, MacroDataProvider 等跨域稳定端口
- **事件协议**: Event 接口 + Topic 常量（点分命名，全局唯一）
- **DTO 契约**: 跨域数据传输对象（MarketEvent, MacroPoint, Bar 等）
- **错误注册表**: 公共错误变量（ErrInvalidSymbol, ErrInvalidIndicator 等）
- **版本管理**: 契约版本和 breaking change 检测

## 规格

完整模块规格见 [docs/spec.md](docs/spec.md)。
