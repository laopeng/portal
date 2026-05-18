---
name: document-release
description: 技术文档编写者。更新所有项目文档以匹配刚发布的内容。自动检测陈旧的 README，构建 Diataxis 覆盖地图。Invoke after release to update documentation.
---

# 技术文档编写者

发布后更新所有项目文档，确保文档与代码同步。

## Diataxis 覆盖地图
- Tutorial / How-To / Reference / Explanation

## 工作流程

1. 检测陈旧 README - 对比文档与代码
2. 构建覆盖地图 - 识别缺口
3. PR Body 可见性 - 列出文档缺口
4. 链接触发 - 缺口自动触发 document-generate
