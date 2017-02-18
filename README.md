IP Service
====
基于 https://github.com/teambition/gear  Go Web 框架和 http://www.ipip.net/ IP 数据库实现的 IP 查询服务。

## 运行

### 从 docker 官方仓库抓取 image 运行
```sh
docker run --rm -p 8080:8080 zensh/ipservice
```

### 从源码运行

```bash
go get github.com/zensh/ipservice
cd path_to_ipservice
go run main.go --data ./data/17monipdb.dat
```

### 编译可执行文件并运行

```bash
# 编译成可运行的二进制文件：ipservice
go build -o ipservice main.go

# 未提供参数显示帮助信息
./ipservice
# 指定 IP 数据库
./ipservice -data ./data/17monipdb.dat
# 指定 IP 数据库并指定监听端口
./ipservice -data ./data/17monipdb.dat -port 3000
```

### Docker (15.01 MB)

```sh
make docker
```

Try it:
```sh
make run
curl 127.0.0.1:8080/json/8.8.8.8
```

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
