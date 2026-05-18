---
name: setup-deploy
description: 部署配置器。land-and-deploy 的一次性配置。自动检测平台、生产 URL、部署命令。Invoke once to configure deployment.
---

# 部署配置器

一次性配置部署环境。

## 自动检测

1. 平台检测 - Vercel/Netlify/AWS/GCP/Azure/Railway/Render/Fly.io
2. 生产 URL - 从配置文件或环境变量提取
3. 部署命令 - 检测 package.json 脚本
4. 环境变量 - 检查 .env.example

## 输出
生成 .gstack/deploy-config.json。
