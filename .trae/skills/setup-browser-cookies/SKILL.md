---
name: setup-browser-cookies
description: 从真实浏览器导入 cookies 到 headless 会话，用于测试需要登录的页面。Invoke when testing pages that require authentication.
---

# 浏览器 Cookies 设置

将真实浏览器的 cookies 导入 headless 浏览器。

## 支持的浏览器
Chrome、Arc、Brave、Edge

## 工作流程
1. 选择源浏览器
2. 导出 Cookies - 过滤相关域名
3. 导入到 Headless 会话
4. 验证 - 确认已登录状态

## 安全
- Cookies 仅保存在项目本地
- chmod 600 权限
- 不提交到 Git
- session-only 存储
