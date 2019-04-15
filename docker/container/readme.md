容器指令
生成并运行一个容器实例
```zsh
$ docker container run hello-world
```
终止一个容器的运行
```zsh
$ docker container kill [containID]
```
列出所有容器（包括停止运行的）
```zsh
$ docker container ls --all
```
删除容器文件
```zsh
$ docker container rm [containID]
```
运行容器指令
这时你就已经在容器中，可以按Ctrl+c停止容器，Ctrl+d退出容器
```zsh
$ docker container run \
-p 8000:3000 \
-it koa-demo:0.1.1 \
--volume "$PWD/":/var/www/html \
--name wordpress \
--env MYSQL_ROOT_PASSWORD=123456 \
--env MYSQL_DATABASE=wordpress \
--link wordpressdb:mysql \
/bin/bash
```
```zsh
- -p：将容器的3000端口映射到本机的8000端口
- -it：容器的shell映射到当前的shell，当前的输入会传到容器中
- /bin/bash：容器启动后第一个执行的指令，`会覆盖Dockerfile中的CMD指令`
- --rm：如果加入这个flag，那么在容器停止运行后自动删除容器文件
- --volume：将当前目录映射到Apache对外访问的默认目录并设置数据卷
- --name：容器名
- --env：设置mysql的环境变量
- -d：容器后台运行
- --link：该容器连接到wordpressdb容器，别名是mysql
```
每次运行docker container run都会生成一个container文件。
```zsh
$ docker container start [containID]
```
查看运行中容器的输出
```zsh
$ docker container logs [containID]
```
进入正在运行的容器中
```zsh
$ docker container exec -it [containID] /bin/bash
```
将正在运行的容器中的文件拷贝到本机
```zsh
$ docker container cp [containID]:[/path/to/file] .
```
.dockerignore文件
创建image镜像时忽略的文件
```zsh
$ touch .dockerfile
$ vim .dockerfile
.git
node_modules
xxxxx.log
```
Dockerfile：新建image镜像
拉取node的8.4版本的image文件，命名为othername
```zsh
FROM node:8.4 as othername
```
设置维护人
```zsh
maintainer [username] "[useremail]"
```
将当前目录除了.dockerignore打包进/app目录下
```zsh
COPY . /app
```
指定接下来的工作目录为/app
```zsh
WORKDIR	/app
```
在/app目录下，运行指令安装依赖，这些都将打包进image中
```zsh
RUN npm install --registry=https://registry.npm.taobao.org
```
容器将3000端口暴露出来
```zsh
EXPOSE 3000
```
设置容器运行后执行的指令
```zsh
CMD ["node demos/01.js"]
```
设置容器运行后执行的指令；如果你在docker run的时候，后面增加的参数会替换掉dockerfile中的CMD指令。
使用ENDPOINT的话，docker run后的的参数会加在ENDPOINT后面。
```zsh
ENDPOINT ["node demos/01.js"]
```
将image文件发布到网上
登陆
```zsh
$ docker login
```
标注用户名和版本
```zsh
$ docker image tag koa-demo:0.1.1 [username]/[repository]:[tag]
```
发布
```zsh
$ docker image push [username]/[repository]:[tag]
```
Multi-Stage Builds功能