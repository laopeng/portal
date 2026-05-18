---
name: scrape
description: 只读页面数据提取。3 条路径：匹配路径 ~200ms、原型路径 ~30s、写入意图拒绝。Invoke to extract data from web pages.
---

# 页面数据提取

只读页面数据提取，支持多种提取路径。

## 核心原则：只读
/scrape 是只读操作。绝不写入或修改页面。

## 三条提取路径

1. 匹配路径 ~200ms - 使用已编码的 browser-skill
2. 原型路径 ~30s - 使用浏览器命令驱动真实浏览器
3. 写入意图拒绝 - 如果意图涉及写入，拒绝执行

## 使用方式
/scrape <意图描述>

## 输出
从页面提取的结构化数据，JSON/CSV/Markdown。
