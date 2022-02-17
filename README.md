# 学生选排课系统


本项目基于Singo开发
https://github.com/Gourouting/singo


## 特色

本项目基于Singo做了以下修改：

1. 加入了Snowflake算法以生成分布式ID。
2. 将Session存入redis中，以支持分布式会话共享。
3. 
4. 
5. 
6. 
7. 
8. 


本项目已经预先创建了一系列文件夹划分出下列模块:

1. api文件夹就是MVC框架的controller，负责协调各部件完成任务
2. model文件夹负责数据库表实体
3. service是MVC框架的Service层，负责处理业务逻辑
4. serializer储存通用的json模型，把model得到的数据库模型转换成api需要的json对象
5. cache负责redis缓存相关的代码
6. auth权限控制文件夹
7. common一些通用工具、错误状态码、常量等
8. conf放一些静态存放的配置文件
9. db存放初始化数据库文件
10. vo存放页面输入模型
11. middleware存放gin相关中间件

## Godotenv

项目在启动的时候依赖以下环境变量，但是在也可以在项目根目录创建.env文件设置环境变量便于使用(建议开发环境使用)

```shell
MYSQL_DSN="db_user:db_password@/db_name?charset=utf8&parseTime=True&loc=Local" # Mysql连接地址
REDIS_ADDR="127.0.0.1:6379" # Redis端口和地址
REDIS_PW="" # Redis连接密码
REDIS_DB="" # Redis库从0到10
SESSION_SECRET="setOnProducation" # Seesion密钥，必须设置而且不要泄露
GIN_MODE="debug"
```


## 运行

```shell
go run main.go
```

项目运行后启动在3000端口（可以修改，参考gin文档)