```shell
$ mkdir mynginx
$ cd mynginx
$ touch Dockerfile
```
内容是：
```txt
FROM nginx
RUN echo '<h1>Hello Docker!</h1>' > /usr/share/nginx/html/index.html
```
# `FROM`：指定基础镜像
如果你以`scratch(表示空白镜像)`为开始的话，那么你接下来的指令将作为镜像的第一层。
```txt
FROM scratch
```
# `RUN`：执行命令行命令，有两种格式
- shell格式：
```txt
RUN echo '<h1>Hello Docker!</h1>' > /usr/share/nginx/html/index.html
```

- exec格式：
```txt
RUN ["可执行文件","参数1","参数2"]
```
## 优化Dockfile
错误写法：
```txt
FROM debain:jessie

RUN apt-get update
RUN apt-get install -y gcc libc6-dev make
RUN wget -O redis.tar.gz "http://download.redis.io/releases/redis-3.2.5.tar.gz"
RUN mkdir -p /usr/src/redis
RUN tar -xzf redis.tar.gz -C /usr/src/redis --strip-components=1
RUN make -C /usr/src/redis
RUN make -C /usr/src/redis install
```
正确写法：
```txt
FROM debain:jessie

#编译、安装redis可执行文件
RUN buildDeps='gcc libc6-dev make' \
	&& apt-get update \
	&& apt-get install $buildDeps \
	&& wget -O redis.tar.gz "http://download.redis.io/releases/redis-3.2.5.tar.gz" \
	&& mkdir -p /usr/src/redis \
	&& tar -xzf redis.tar.gz -C /usr/src/redis --strip-components=1 \
	&& make -C /usr/src/redis \
	&& make -C /usr/src/redis install \
	&& rm -rf /var/lib/apt/lists/* \				#清理缓存
	&& rm redis.tar.gz \							#清理下载文件
	&& rm -r /usr/src/redis \						
	&& apt-get purge -y --auto-remove $buildDeps	#清理为了编译构建所需软件
```
## 构建镜像
```shell
$ docker build -t nginx:v3 .
Sending build context to Docker daemon  2.048kB
Step 1/2 : FROM nginx
 ---> e4e6d42c70b3
Step 2/2 : RUN echo '<h1>Hello Docker!</h1>' > /usr/share/nginx/html/index.html
 ---> Running in 48ed6a61a80e
 ---> 74fe963dff36
Removing intermediate container 48ed6a61a80e
Successfully built 74fe963dff36
Successfully tagged nginx:v3
```
可以看到：显示RUN启动一个容器`48ed6a61a80e`，执行要求的指令，最后提交给这一层`74fe963dff36`，最后删掉所用的容器。

## 其他构建方法
### 1.直接用Git repo构建
```shell
$ docker build https://github.com/twang2218/gitlab-ce-zh.git#:8.14
```

### 2.给定的tar压缩包构建
Docker引擎会自动下载，并解压缩，作为上下文，构建镜像
```shell
$ docker build http://server/context.tar.gz
```

### 3.从标准输入中读取Dockerfile进行构建
```shell
$ docker build - < Dockerfile
```
或
```shell
$ cat Dockerfile | docker build -
```

### 4.从标准输入中读取上下文压缩包进行构建
```shell
$ docker build - < context.tar.gz
```

# `COPY`：复制文件
格式：
- COPY <源路径>...<目标路径>
- COPY ["<源路径1>",..."<目标路径>"]

```dockerfile
 COPY package.json /usr/src/app
#可以是通配符
 COPY hom* /mydir/
 COPY hom?.txt /mydir/
```
**注意**：源文件的各种元数据（权限，变更时间）都会保留。

# `ADD`：更高级的复制文件（需要自动解压缩时用）
有个功能非常有用：如果源路径是压缩文件，它会自动解压到目标路径去。
```dockerfile
FROM scratch
ADD ubuntu-xenial-core-cloudimg-amd64-root.tar.gz /
```

# `CMD`：指定容器启动程序和参数
格式：
- shell格式：CMD <命令>
- exec格式：CMD ["可执行文件","参数1","参数2"...]
- 参数列表格式：CMD ["参数1","参数2"...]，在指定了`ENTRYPOINT`指令后，用`CMD`指定具体的参数。

```dockerfile
CMD echo $HOME
# 其实就等于
CMD ["sh","-c","echo $HOME"]
```
那么如果你要将nginx程序在前台运行；
```dockerfile
CMD ["nginx","-g","daemon off;"]
```

# `ENTYRPOINT`：入口点

### 使用一：让镜像变成像命令一样使用
得知自己当前公网ip的镜像
```dockerfile
FROM ubuntu:14.04
RUN apt-get update \
	&& apt-get install -y curl \
	&& rm -rf /var/lib/apt/list/* 
	CMD ["curl","-s","http://ip.cn"]
```
使用`docker build -t myip .`来构建镜像，如果需要查询ip：
```shell
$ docker run myip
当前 IP：116.224.202.178 来自：上海市 电信
```
如果你想希望加入`-i`参数：

**错误做法**
```shell
#因为加了参数-i的话，运行时会默认替换掉CMD的默认值
$ docker run myip -i
```
**正确做法**
```dockerfile
FROM ubuntu

RUN apt-get update \
	&& apt-get install -y curl \
	&& rm -rf /var/lib/apt/list/*
	ENTRYPOINT ["curl","-s","http://ip.cn"]
```
这样就可以使用`docker run myip:v3 -i`了。
```shell
$ docker run myip -i
$ curl -s http://ip.cn -i
HTTP/1.1 200 OK
Server: nginx/1.11.9
Date: Sun, 23 Jul 2017 05:43:03 GMT
Content-Type: text/html; charset=UTF-8
Transfer-Encoding: chunked
Connection: keep-alive

当前 IP：116.224.202.178 来自：上海市 电信
```
因为当存在`ENTRYPOINT`时，`CMD`内容会被作为参数传给`ENTRYPOINT`，而`-i`就是新的`CMD`，会作为参数传给`curl`。

# `ENV`：设置环境变量
```dockfile
ENV VERSION=1.0 DEBUG=on \
	NAME="Happy Feet"
```
# `ARG`：构建参数
在Dockerfile中定义参数和默认值，该默认值可以在`build`中用`--build-arg <参数名>=<值>`替换

# `VOLUME`：定义匿名卷
容器运行时应该尽量保持容器存储层不发生写操作，对于数据库类需要保存动态数据的应用，其数据库文件应该保存于卷(volume)中，为了防止运行时用户忘记将动态文件所保存目录挂载为卷，在 Dockerfile 中，我们可以事先指定某些目录挂载为匿名卷，这样在运行时如果用户不指定挂载，其应用也可以正常运行，不会向容器存储层写入大量数据。
```dockerfile
VOLUME /data
```
构建时覆盖了这个设置，将mydata挂载到匿名卷中。
```shell
$ docker run -d -v mydata:/data xxxx
```

# `EXPOSE`：声明端口
`EXPOSE`只是声明打算用什么端口，一个用处是在运行时使用随机端口映射，也就是`docker run -P`，会自动映射`EXPOSE`的端口。

# `WORKDIR`：工作目录
指定了工作目录后，以后各层的当前目录就被改为指定的目录。

# `USER`：指定当前用户
```dockerfile
RUN groupadd -r redis && useradd -r -g redis redis
USER redis
RUN [ "redis-server" ]
```
如果以`root`执行的脚本，在执行期间希望改变身份，不要使用su或者sudo，建议使用[gosu](https://github.com/tianon/gosu)
```dockerfile
# 建立 redis 用户，并使用 gosu 换另一个用户执行命令
RUN groupadd -r redis && useradd -r -g redis redis
# 下载 gosu
RUN wget -O /usr/local/bin/gosu "https://github.com/tianon/gosu/releases/download/1.7/gosu-amd64" \
    && chmod +x /usr/local/bin/gosu \
    && gosu nobody true
# 设置 CMD，并以另外的用户执行
CMD [ "exec", "gosu", "redis", "redis-server" ]
```

# `HEALTHCHECK`：检查容器健康状况
格式：
- HEALTHCHECK [选项] CMD <命令>
- HEALTHCHECK NONE ：如果基础镜像有健康检查，可以评比健康检查

`HEALTHCHECK`支持下列选项：
- `--interval=<间隔>`：两次检查的间隔时间，默认30秒
- `--timeout=<时间>`：超过这个时间，就是不健康，默认30秒
- `--retries=<次数>`：重试几次后，就是不健康，默认3次

可以用`curl`帮助判断
```dockerfile
FROM nginx 

RUN apt-get update \
	&& apt-get insatll -y curl \
	&& rm -rf /var/lib/apt/list/* 
	HEALTHCHECK --interval=5s --timeout=3s \
		CMD curl -fs http://localhost/ || exit 1
```
构建镜像
```shell
$ docker build -t myweb:v1 .
```
启动容器
```shell
docker run --name web -d -p 80:80 myweb:v1
```
然后我们可以通过`docker ps`查看状态：
```shell
$ docker ps
CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS                            PORTS               NAMES
03e28eb00bd0        myweb:v1            "nginx -g 'daemon off"   3 seconds ago       Up 2 seconds (health: starting)   80/tcp, 443/tcp     web
```
为了排障，健康命令的输出（包括stdout，stderr）都会存储在健康状态里，用`docker inspect`查看。
```shell
$ docker inspect --format '{{json .State.Health}}' web | python -m json.tool
{
    "FailingStreak": 0,
    "Log": [
        {
            "End": "2016-11-25T14:35:37.940957051Z",
            "ExitCode": 0,
            "Output": "<!DOCTYPE html>\n<html>\n<head>\n<title>Welcome to nginx!</title>\n<style>\n    body {\n        width: 35em;\n        margin: 0 auto;\n        font-family: Tahoma, Verdana, Arial, sans-serif;\n    }\n</style>\n</head>\n<body>\n<h1>Welcome to nginx!</h1>\n<p>If you see this page, the nginx web server is successfully installed and\nworking. Further configuration is required.</p>\n\n<p>For online documentation and support please refer to\n<a href=\"http://nginx.org/\">nginx.org</a>.<br/>\nCommercial support is available at\n<a href=\"http://nginx.com/\">nginx.com</a>.</p>\n\n<p><em>Thank you for using nginx.</em></p>\n</body>\n</html>\n",
            "Start": "2016-11-25T14:35:37.780192565Z"
        }
    ],
    "Status": "healthy"
}
```

# `ONBUILD`：为他人做嫁衣裳
`ONBUILD`是一个特殊的指令，它后面跟的是其它指令，比如`RUN`,`COPY`等，而这些指令，在当前镜像构建时并不会被执行。只有当以当前镜像为基础镜像，去构建下一级镜像的时候才会被执行。
请看下面的例子：
假设我们制作Node.js所写应用的镜像。我们知道Node.js使用`npm`包进行管理，所有的以来，启动，配置都放在`package.json`中。我们首先进行`npm`包进行管理，所有的以来，启动，配置都放在`package`里，我们首先要`npm insatll`，然后`npm start`：
```dockerfile
FROM node.slim
RUN mkdir /app
WORKDIR /app
COPY ./package.json ./app
RUN ["npm","install"]
COPY . ./app/
CMD ["npm","start"]
```
但是这样维护起来比较麻烦，我们可以构建一个基础镜像，然后继承他：
```dockerfile
FROM node:slim
RUN mkdir /app
WORKDIR /app
CMD ["npm","start"]
```
假设这个基础镜像名为`my-node`，子项目的dockerfile就是：
```dockerfile
FROM my-node
COPY ./package.json /app
RUN ["npm","insatll"]
COPY . /app/
```
但是如果`npm insatll`需要加一些参数怎么办，可以使用`ONBUILD`解决：
```dockerfile
FROM node:slim
RUN mkdir /app
WORKDIR /app
ONBUILD CPOY ./package.json /app
ONBUILD RUN ["npm","install"]
ONUBILD CPOY . /app/
RUN ["npm","start"]
```
那么各个子项目的`Dockerfile`就变成：
```dockerfile
FROM my-node
```