<h1 align="center">easyimage_go</h1>

此仓库是一个基于`gin`框架的图床服务，只提供最简单给自己用的图床

另有`picgo`[插件](https://github.com/buyfakett/picgo-plugin-easyimage-go),可以直接上传图片到图床

[前端](https://github.com/buyfakett/easyimage_go_web)

## 技术栈

- [Go](https://golang.org/)

- [gin](https://github.com/gin-gonic/gin)

## 快速启动

### Docker compose

```yaml
services:
    easyimage_go:
        image: buyfakett/easyimage_go
        container_name: easyimage_go
        network_mode: host
        restart: always
        volumes:
            - ./config/config.yaml:/app/config.yaml:ro
            - ./i:/app/i
        command: --config=/app/config.yaml
```

### 配置文件

```
server:
  port: 8080                    # 服务端口
  domain: http://localhost:8080 # 服务域名(用于拼接图片url)
  token: 123456                 # 鉴权token(用于鉴权)
image:
  uri: /i                       # 图片存储路径(相对路径)
  webp_quality: 100             # webp压缩质量(0-100)
```

## 项目目录

```tree
.
├── Dockerfile                  # Dockerfile
├── biz                         # 业务代码
│     ├── handler               # 服务逻辑
│     ├── mw                    # 中间件
│     └── router                # 路由
├── build.sh                    # 编译脚本
├── config                      # 配置文件
│     ├── config.yaml           # 配置文件(可以覆盖默认配置)
│     └── default.yaml          # 默认配置文件(服务端这里定义的默认配置)
├── docs                        # swagger文档
├── main.go                     # 启动文件
├── static                      # 静态文件(前端编译结果，必须要index.html)
└── utils                       # 工具包
```

## 开发

### 启动

如果需要指定配置文件，可以使用以下命令

```bash
go run . --config=config/config.yaml
```

### 自动化

目前使用`github actions`自动化,开发环境每个`commit`会自动编译docker镜像,打v1.0.0的标签的时候会编译docker镜像和二进制文件到`release`下
