# 部署



### 从github拉取代码

```shell
git clone https://github.com/fandypeng/excel2config
cd excel2config/configs
```



### 配置mongodb数据库

进入excel2config/configs目录修改项目配置，项目使用mongodb作为数据库，先编辑mongo.toml

```toml
[Client]
	name = "mongo"
	proto = "tcp"
	addr = "127.0.0.1:27017"
	authSource = "sheets"
	dsn = "mongodb://admin:123456@127.0.0.1:27017/?authSource=sheets"
```

dsn字段填写mongodb链接的dsn，数据库名称可以自定义，默认情况下是sheets



### 配置登录方式

项目支持邮箱注册、LDAP登录、公司邮箱登录三种登录方式，可以在configs/application.toml里面配置

```toml
[Ldap]
    serverHost = "ldap://ldap.excel2config.com:389" # 服务器域名和端口号
    bindBaseDn = "cn=admin,dc=excel2config,dc=com"  # 登录的起始索引
    bindPassword = "password" 											# 登录的密码
    searchBaseDn = "dc=excel2config,dc=com"					# 搜索的起始索引
    searchBy = "email"															# 搜索属性名
    searchUserName = "fandy@excel2config.com"				# 搜索属性值

[Dingtalk]
    # reference: https://developers.dingtalk.com/document/app/enterprise-internal-application-logon-free
    corpHost = "https://oapi.dingtalk.com"					# 钉钉开放平台域名
    corpId = "ding639831aaa1212c9a"									# 钉钉公司ID
    corpSecret = "v1iZc8pp21dsghgR3VdcDK-845B3Rqq3sjPzEk_lJ0qR4Wasaf5TdC2cO9IgPHUdz" # 钉钉公司密码
    chatId = "122385892"														# 群聊ID，不为空的情况下，群内用户才可以登录系统

```



如果配置了钉钉登录，那么需要在钉钉的控制台加入如下链接来获取登录方式：

```html
http://e2c.fandypeng.com/#/login?code=
```

在钉钉的控制台，动态加入code到url末尾，即可实现钉钉登录



### 配置HTTP和WebSocket服务

HTTP服务支持监听、超时、跨域等配置，可以在configs/application.toml里面配置

```toml
[Server]
    addr = "0.0.0.0:8000"																										# 监听地址和端口
    timeout = "2s"																													# HTTP接口超时时间
    crossDomains = [ "http://localhost:9528", "http://e2c.fandypeng.com" ]  # 跨域开放的域名，配置了新域名一定要添加
[WebSocket]
    addr = "0.0.0.0:8001"
    readTimeout = "5s"
    writeTimeout = "5s"
```



默认情况下，HTTP服务监听的是8000端口，WebSocket监听的是8001端口，如果端口号被占用需要更换的话，可以直接修改配置文件修改监听端口。修改了HTTP服务监听的端口号，需要同时修改前端API的端口号，文件位于：`front/.env.production`

```toml
# just a flag
ENV = 'production'

# base api
VUE_APP_BASE_API = 'http://e2c.fandypeng.com:8000'
```

VUE_APP_BASE_API 配置的是API的域名，所有API请求都是在此域名下。





### 启动服务

执行项目根目录的update.sh，即可编译项目并启动服务

```shell
[root@zlskjsagskjgja excel2config]# ./update.sh
来自 https://github.com/fandypeng/excel2config
Already up-to-date.
restart succeed
nohup: 把输出追加到"nohup.out"
[root@zlskjsagskjgja excel2config]# ps -ef | grep excel2config
root     25705     1  1 11:34 pts/0    00:00:02 ./excel2config -conf ./configs/
root     25719 25621  0 11:36 pts/0    00:00:00 grep --color=auto excel2config
```

启动之后，执行 `ps -ef | grep excel2config` 即可看到进程ID，可以通过nohup.out查看进程输出，排查启动异常问题



### 配置前端服务

切换到根目录下面的front目录，编译前端资源，配置nginx的静态文件目录

```shell
[root@zlskjsagskjgja excel2config]# cd front/
[root@zlskjsagskjgja front]# npm install
npm WARN Excel2Config@1.0.0 No repository field.
npm WARN optional SKIPPING OPTIONAL DEPENDENCY: fsevents@2.1.3 (node_modules/fsevents):
npm WARN notsup SKIPPING OPTIONAL DEPENDENCY: Unsupported platform for fsevents@2.1.3: wanted {"os":"darwin","arch":"any"} (current: {"os":"linux","arch":"x64"})
npm WARN optional SKIPPING OPTIONAL DEPENDENCY: fsevents@1.2.13 (node_modules/jest-haste-map/node_modules/fsevents):
npm WARN notsup SKIPPING OPTIONAL DEPENDENCY: Unsupported platform for fsevents@1.2.13: wanted {"os":"darwin","arch":"any"} (current: {"os":"linux","arch":"x64"})
npm WARN optional SKIPPING OPTIONAL DEPENDENCY: fsevents@1.2.13 (node_modules/watchpack-chokidar2/node_modules/fsevents):
npm WARN notsup SKIPPING OPTIONAL DEPENDENCY: Unsupported platform for fsevents@1.2.13: wanted {"os":"darwin","arch":"any"} (current: {"os":"linux","arch":"x64"})
npm WARN optional SKIPPING OPTIONAL DEPENDENCY: fsevents@1.2.13 (node_modules/webpack-dev-server/node_modules/fsevents):
npm WARN notsup SKIPPING OPTIONAL DEPENDENCY: Unsupported platform for fsevents@1.2.13: wanted {"os":"darwin","arch":"any"} (current: {"os":"linux","arch":"x64"})

audited 1726 packages in 21.745s

68 packages are looking for funding
  run `npm fund` for details

found 4914 vulnerabilities (4 low, 201 moderate, 4709 high)
  run `npm audit fix` to fix them, or `npm audit` for details
[root@zlskjsagskjgja front]# ./update.sh
Already up-to-date.

> Excel2Config@1.0.0 build:prod /var/go/src/excel2config/front
> vue-cli-service build


⠼  Building for production...

 WARNING  Compiled with 2 warnings                                                                                                                                                                              上午11:54:09

 warning

asset size limit: The following asset(s) exceed the recommended size limit (244 KiB).
This can impact web performance.
Assets:
  static/js/chunk-25197187.0798c8cc.js (909 KiB)
  static/js/chunk-elementUI.cb459a4a.js (653 KiB)
  static/js/chunk-libs.79ca039e.js (374 KiB)

 warning

entrypoint size limit: The following entrypoint(s) combined asset size exceeds the recommended limit (244 KiB). This can impact web performance.
Entrypoints:
  app (1.25 MiB)
      static/css/chunk-elementUI.6188a60f.css
      static/js/chunk-elementUI.cb459a4a.js
      static/css/chunk-libs.bf952545.css
      static/js/chunk-libs.79ca039e.js
      static/css/app.f86da1aa.css
      static/js/app.862b3555.js


  File                                      Size             Gzipped

  dist/static/js/chunk-25197187.0798c8cc    909.01 KiB       296.84 KiB
  .js
  dist/static/js/chunk-elementUI.cb459a4    653.13 KiB       159.99 KiB
  a.js
  dist/static/js/chunk-libs.79ca039e.js     373.56 KiB       130.07 KiB
  dist/static/js/app.862b3555.js            39.71 KiB        13.88 KiB
  dist/static/js/chunk-5667abd6.8867bce6    38.62 KiB        8.33 KiB
  .js
  dist/static/js/chunk-3d3f5d7a.5839c992    5.68 KiB         1.97 KiB
  .js
  dist/static/js/chunk-39939de9.754682aa    4.15 KiB         1.44 KiB
  .js
  dist/static/js/chunk-a88334d6.2473b18a    3.76 KiB         1.48 KiB
  .js
  dist/static/js/chunk-7a968e41.cea9a7b8    2.96 KiB         0.97 KiB
  .js
  dist/static/js/chunk-29dc5dd0.fb2ae0a6    1.87 KiB         0.96 KiB
  .js
  dist/static/js/chunk-76cea4de.5ea47175    1.74 KiB         0.73 KiB
  .js
  dist/static/js/chunk-2d0c8bf7.350506e2    0.41 KiB         0.31 KiB
  .js
  dist/static/js/chunk-2d0e4e1f.d2235087    0.41 KiB         0.31 KiB
  .js
  dist/static/js/chunk-2d226cab.334c23ce    0.39 KiB         0.29 KiB
  .js
  dist/static/js/chunk-2d229205.f9fd3364    0.37 KiB         0.29 KiB
  .js
  dist/static/js/chunk-2d0cfaef.4803ac5d    0.35 KiB         0.27 KiB
  .js
  dist/static/js/chunk-2d0e944c.e6094cf7    0.35 KiB         0.28 KiB
  .js
  dist/static/js/chunk-2d2104c6.484ff941    0.35 KiB         0.27 KiB
  .js
  dist/static/css/chunk-elementUI.6188a6    204.32 KiB       32.75 KiB
  0f.css
  dist/static/css/app.f86da1aa.css          8.80 KiB         2.30 KiB
  dist/static/css/chunk-76cea4de.6e6258d    4.47 KiB         0.82 KiB
  4.css
  dist/static/css/chunk-libs.bf952545.cs    3.21 KiB         1.22 KiB
  s
  dist/static/css/chunk-3d3f5d7a.69552a0    1.75 KiB         0.58 KiB
  b.css
  dist/static/css/chunk-a88334d6.af479e7    1.57 KiB         0.62 KiB
  b.css
  dist/static/css/chunk-39939de9.133e08c    1.54 KiB         0.61 KiB
  6.css
  dist/static/css/chunk-5667abd6.49a0301    0.47 KiB         0.24 KiB
  a.css
  dist/static/css/chunk-29dc5dd0.fc50c43    0.35 KiB         0.24 KiB
  c.css
  dist/static/css/chunk-7a968e41.5cd9884    0.04 KiB         0.06 KiB
  a.css

  Images and other types of assets omitted.

 DONE  Build complete. The dist directory is ready to be deployed.
 INFO  Check out deployment instructions at https://cli.vuejs.org/guide/deployment.html



   ╭────────────────────────────────────────────────────────────────╮
   │                                                                │
   │      New major version of npm available! 6.14.11 → 7.13.0      │
   │   Changelog: https://github.com/npm/cli/releases/tag/v7.13.0   │
   │               Run npm install -g npm to update!                │
   │                                                                │
   ╰────────────────────────────────────────────────────────────────╯

```

编译好前端资源后，需要配置nginx

```nginx
server {
    listen       80;
    server_name  e2c.fandypeng.com; # 填入你要解析的域名
    root         /var/go/src/excel2config/front/dist;

    error_page 404 /404.html;
        location = /40x.html {
    }

    error_page 500 502 503 504 /50x.html;
        location = /50x.html {
    }
}
```



重启nginx

```shell
> service nginx restart
```



此时打开浏览器访问你配置的域名，即可正常访问系统。





