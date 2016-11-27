IP Service
====
基于 https://github.com/teambition/gear  Go Web 框架和 http://www.ipip.net/ IP 数据库实现的 IP 查询服务。

## 运行
### 从源码运行

```bash
go get github.com/zensh/ipservice
cd path_to_ipservice
go run app.go --data=./data/17monipdb.dat
```

### 编译可执行文件并运行

```bash
# 编译成可运行的二进制文件：ipservice
go build -o ipservice app.go

# 未提供参数显示帮助信息
./ipservice
# 指定 IP 数据库
./ipservice --data=./data/17monipdb.dat
# 指定 IP 数据库并指定监听端口
./ipservice --data=./data/17monipdb.dat --port=3000
```

### Docker

## API

### GET /json/:ip

```bash
curl 127.0.0.1:8080/json/8.8.8.8
# 返回 JSON 数据
{"IP":"8.8.8.8","Status":200,"Message":"","Data":{"Country":"GOOGLE","Region":"GOOGLE","City":"N/A","Isp":"N/A"}}
```

### GET /json/:ip?callback=xxx

```bash
# callback=xxxx 返回 JSONP 数据
curl 127.0.0.1:8080/json/8.8.8.8?callback=readIP
# 返回 JSONP 数据
/**/ typeof readIP === "function" && readIP({"IP":"8.8.8.8","Status":200,"Message":"","Data":{"Country":"GOOGLE","Region":"GOOGLE","City":"N/A","Isp":"N/A"}});
```

## Bench

Environment: MacBook Pro, 2.4 GHz Intel Core i5, 8 GB 1600 MHz DDR3

Start service: `./ipservice --data=./data/17monipdb.dat`

Result: **41132.68 req/sec**
```bash
wrk 'http://localhost:8080/json/8.8.8.8' -d 60 -c 100 -t 4
Running 1m test @ http://localhost:8080/json/8.8.8.8
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     3.51ms    6.45ms 199.21ms   95.37%
    Req/Sec    10.34k     1.14k   14.33k    86.12%
  2470564 requests in 1.00m, 558.40MB read
Requests/sec:  41132.68
Transfer/sec:      9.30MB
```
