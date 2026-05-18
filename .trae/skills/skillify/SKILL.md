---
name: skillify
description: 将最近成功的 scrape 原型编码为永久 browser-skill。11 步流程，3 份锁定合约。Invoke after a successful scrape.
---

# Skill 编码

将成功的 scrape 原型编码为可复用的 browser-skill。

## 11 步流程
1. 分析成功 scrape 数据流
2. 提取关键选择器
3. 设计鲁棒的数据提取
4. 处理边界情况
5. 编码为 browser-skill 模板
6. 验证在测试数据上
7. 生成技能描述
8. 添加出处信息
9. 写入技能文件
10. 更新技能索引
11. 清理临时文件

## 3 份锁定合约
1. 出处守卫 - 记录 URL、日期、结构假设
2. 合成输入切片 - 保存成功输入作为回归测试
3. 原子写入 - 一次写入，验证完整性
