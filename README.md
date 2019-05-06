
# microservices
使用micro微服务框架的一些例子，包含微服务开发中必备的鉴权，熔断，监控，链路追踪等
> 本项目使用consul的K/V存储来模拟一个配置中心。在每一个使用到lib/token包的服务都会在main函数里传递consul地址和JWT私钥的加载路径。在srv/user服务中同样，使用consul K/V 来存储mysql数据库的配置。在启动服务之前，需要现在consul的K/V中设置好这些配置。

![image](https://github.com/Allenxuxu/microservices/raw/master/.screenshots/jwt-config.png)
![image](https://github.com/Allenxuxu/microservices/raw/master/.screenshots/mysql-config.png)

## 主要内容
- 扩展micro的API网关功能
  - JWT鉴权
  - 熔断
  - prometheus监控
  - 链路追踪
- gin开发微服务service
  - 熔断
  - prometheus监控
  - 链路追踪
- 内部服务采用grpc通信
## 详细讲解
[个人博客](https://blog.mogutou.xyz/blog/cate/xuxu/Go-Micor)
## 使用到的其他软件
- consul (服务发现,K/V配置)
- prometheus (监控)
- jaeger (链路追踪)
- hystrix-dashboard (hystrix熔断仪表盘)
- mysql 

### docker启动参考命令
- consul
  > docker run --name consul -d -p 8500:8500/tcp consul agent -server -ui -bootstrap-expect=1 -client=0.0.0.0
- prometheus
  > docker run --name prometheus  -d -p 0.0.0.0:9090:9090 -v /tmp/conf.yml:/etc/prometheus/prometheus.yml   prom/prometheus
- jaeger
  > docker run -d --name jaeger \
  -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
  -p 5775:5775/udp \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 14268:14268 \
  -p 9411:9411 \
  jaegertracing/all-in-one:1.6
- hystrix-dashboard
  > docker run --name hystrix-dashboard -d -p 8081:9002 mlabouardy/hystrix-dashboard:latest

  > hystrix数据监控
    http://localhost:8030/hystrix.stream

- mysql
  > docker run --name mysql -e  MYSQL_ROOT_PASSWORD=123 -d -p 3306:3306 mysql

---
prometheus 简单参考配置文件
```
global:
  scrape_interval: 15s
  scrape_timeout: 10s
  evaluation_interval: 15s
alerting:
  alertmanagers:
  - static_configs:
    - targets: []
    scheme: http
    timeout: 10s
scrape_configs:
- job_name: APIGW
  honor_timestamps: true
  scrape_interval: 15s
  scrape_timeout: 10s
  metrics_path: /metrics
  scheme: http
  static_configs:
  - targets:
    - 10.104.34.106:8080   #10.104.34.106为本机ip， 本机127.0.0.1在容器中无法访问到
```
## 快速体验
- 使用docker 启动consul jaeger hystrix-dashboard (上面有参考命令，复制粘贴执行即可)

- 打开浏览器，进入 http://localhost:8500，进入K/V存储设置JWT私钥配置(参考上面的截图)

- jaeger UI http://localhost:16686

- hystrix-dashboard UI http://localhost:8081/hystrix, 输入 http://{ip}:81/hystrix.stream , 此处ip为本机ip，因为hystrix-dashboard是容器启动的，无法直接访问本机127.0.0.1

- 启动API网关： 
  ```
  cd micro && make run
  ```
- 启动user API服务： 
  ```
  cd api/user &&  make run
  ```
- 启动hello 服务： 
  ```
  cd srv/hello && make run
  ```
- 浏览器访问 http://127.0.0.1:8080/user/test ，或者使用其他工具 GET 127.0.0.1:8080/user/test

![image](https://github.com/Allenxuxu/microservices/raw/master/.screenshots/consul-service.png)
