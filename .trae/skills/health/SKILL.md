---
name: health
description: 技能健康面板。检查所有 gstack 技能的状态和健康度。Invoke to check the status of all installed skills.
---

# 技能健康面板

检查所有已安装 gstack 技能的健康状态。

## 检查项目

1. 技能存在性 - 目录和 SKILL.md 是否存在
2. 技能完整性 - name/description/正文
3. 技能状态 - 是否需要更新
4. 依赖检查 - 前置技能是否存在

## 输出
生成 .gstack/skill-health.md，包含完整的技能健康报告。

## 使用方式
/health
