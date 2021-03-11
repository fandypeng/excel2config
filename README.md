# excel2config
Build configs for your app use online excel , then export to sql、json、lua file and any other files, excel is really friendly to developer、designer、product manager. 

# 简介
作为一个游戏行业从业五年的老鸟，我深切感受到了游戏行业配置管理系统的重要性。目前业内并没有一个统一的管理游戏配置的工具，从我目前接触到的公司和项目情况来看，目前市面上常见的配置管理有如下几种：
- 编写后台系统，使用表单CRUD(适用于业务变化频率低，维护难度大)
- 使用Excel管理配置，使用脚本解析之后导入数据库(管理较为灵活，可视化较差，无法多人协作编辑，Excel文件容易出现冲突)
- 使用json、lua等脚本语言编写(管理较为灵活，可视化程度低，编辑需要一定专业技能)

从实际的使用场景角度来看，配置不仅仅是开发同学来使用，还有产品、运营、策划等同学来使用，所以json、lua脚本的方案首先排除掉，后台表单只适用于变更不频繁的业务系统。由于Excel是策划、运营等同学的必备技能，所以使用Excel管理配置是最有选择，但是这个方案依然有许多问题，比如无法多人协作编辑、容易出现修改冲突、无法比对修改内容等。

Excel2Config系统使用在线Excel来支持多人在线协作编辑表格，不会出现文件冲突的问题，可以在导出配置的时候查看到本次修改的内容。

Excel2Config支持LDAP登录，部署一套系统可以支持项目层面的数据隔离，支持同一个公司的多个项目和人员的管理。


# 项目结构
本项目基于B站开源的 [kratos](https://github.com/go-kratos/kratos) 框架编写，在框架基础上引入websocket服务，并增加mongodb存储Excel数据。


# 体验地址
[点我马上体验](http://e2c.fandypeng.com)

示例项目帐号：demo@163.com，密码：demo-project

![恰饭](http://cdn.fandypeng.com/donate.jpg)
