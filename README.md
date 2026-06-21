# contracts

`contracts` 定义 FoundationX 跨域稳定契约——端口（接口）、事件协议和 DTO。它是域间通信的唯一合法通道，确保数据域、分析域、决策域和执行域之间的接口稳定、可演进。

## 定位

contracts 遵循 xlib-standard 的治理协议，但不是标准源、不是 generator、不是模板仓库。contracts 不拥有传输实现，不绑定具体通信协议。

## 版本

- 当前发布版本：`v0.4.7`

## 核心内容

- **端口接口**: MarketDataProvider, MacroDataProvider 等跨域稳定端口
- **事件协议**: Event 接口 + Topic 常量（点分命名，全局唯一）
- **DTO 契约**: 跨域数据传输对象（MarketEvent, MacroPoint, Bar 等）
- **错误注册表**: 公共错误变量（ErrInvalidSymbol, ErrInvalidIndicator 等）
- **版本管理**: 契约版本和 breaking change 检测

## 规格

完整模块规格见 [docs/spec.md](docs/spec.md)。
