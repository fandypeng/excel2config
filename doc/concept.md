# Excel2Config

产品名称，将Excel内容转换成配置数据，简称e2c



# e2cdatabus



[项目主页](https://github.com/fandypeng/e2cdatabus) 

如果你希望使用Excel来管理配置，但是部署excel2config项目的服务器又无法访问你的数据库， 那么你可以通过在能读写数据库的项目中集成e2cdatabus来接收Excel的配置数据，由e2cdatabus将配置数据写入数据库。





# LuckySheet



[项目主页](https://github.com/mengshukeji/Luckysheet)

🚀Luckysheet ，一款纯前端类似excel的在线表格，功能强大、配置简单、完全开源，并且可以自定义服务器支持在线协作。



## 配置表



一个sheet代表一张配置表，可以通过多个配置表的组合来描述一个较为复杂的配置数据，sheet的命名必须使用英文，建议是使用小写字母加下划线的命名方式，因为sheet的名称会被作为MySQL的表名，或者redis的key名称



## 配置集合



一个Excel对应一个配置集合，同类型或者同功能模块的配置表可以存放在一个Excel里面



## GridKey



Excel的唯一ID，在Excel编辑页面的链接里面可以查询到，点击Excel列表里的名称可以进入到编辑Excel的页面，页面链接如下：

`http://xxx.xxx.com/#/config/excel/英雄列表/609cda552e144eaecc6b5dfa`

此链接中的 `609cda552e144eaecc6b5dfa` 即为gridKey

