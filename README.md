# Excel2Config
[![Go Report Card](https://goreportcard.com/badge/github.com/fandypeng/excel2config)](https://goreportcard.com/report/github.com/fandypeng/excel2config)
[![Release](https://img.shields.io/github/v/release/fandypeng/excel2config.svg?style=flat-square)](https://github.com/fandypeng/excel2config)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)


# 简介

Excel2Config是一个基于在线Excel的配置管理工具，他具有多人在线协作、发布前差异比对、发布历史记录查看和回滚、正式环境和测试环境隔离等特性。

因为Excel灵活、编辑友好的特点，您可以将Excel2Config应用于任何软件开发的场景，尤其适合用于游戏项目，与产品、策划、运营共同维护配置数据。

[体验地址](http://e2c.fandypeng.com)

测试帐号：demo@163.com

测试密码：demo-project



# 配置管理工具对比

目前管理项目的配置数据大体有如下几种方法：

1. 离线Excel编辑，使用PHP、Python等脚本读取Excel数据导入数据库
2. 后台编写表单操作数据，保存到数据库
3. 使用JSON、lua、xml等代码编辑配置，应用服务读取解析



Excel2Config对比其他解决方案：

| 对比项/对比产品 | Excel2Config | 离线Excel | 网页表单 | JSON文本 |
| --------------- | ------------ | --------- | -------- | -------- |
| 可读性          | 优           | 优        | 优       | 差       |
| 可扩展性        | 优           | 一般      | 差       | 优       |
| 发布效率        | 优           | 差        | 优       | 一般     |
| 部署效率        | 优           | 差        | 差       | 差       |
| 协作效率        | 优           | 差        | 差       | 差       |


# 详细文档

请参阅[详细文档](doc/readme.md)