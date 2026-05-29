# admin-web 项目说明 (for Codex)

## 项目定位

此前端项目为 `admin-api` 的管理后台。

## 技术栈

- Vue 3 / React（根据实际情况写）
- 状态管理库：...
- UI 组件库：...
- 请求工具：axios，封装在 `src/utils/http.ts`

## 接口调试说明

- 后端 API 基础地址配置在 `.env.development` 文件的 `VITE_API_BASE_URL` 变量中。
- 所有接口定义在 `src/api/` 目录下，按模块划分。
- 接口代理配置在 `vite.config.ts` 的 `server.proxy` 中。

## 工作约定

- 修改代码前，先阅读相关文件。
- 仅修改前端代码，不要触碰后端逻辑。
- 如果接口请求失败，停下分析并向我报告，不要尝试修改后端代码。
- 改完后，运行 `npm run typecheck` 和 `npm run lint` 验证。

## 常用验证命令

- `npm run dev`: 启动开发服务器
- `npm run lint`: 代码检查
- `npm run typecheck`: 类型检查

