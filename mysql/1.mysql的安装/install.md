# 在Windows上通过压缩包安装
1. [下载](https://dev.mysql.com/downloads/mysql/5.6.html#downloads)mysql的压缩包。
2. 添加环境变量
3. 修改my-default.ini更名为my.ini
```ini
[mysqld]
loose-default-character-set = utf8
character-set-server = utf8
basedir = D:\Program Files\mysql-5.6.23-win32
datadir = D:\Program Files\mysql-5.6.23-win32\data

[client]
#设置客户端字符集  
loose-default-character-set = utf8

[WinMySQLadmin]
Server = D:\Program Files\mysql-5.6.23-win32\bin\mysqld.exe
```
4. 管理员身份打开cmd，输入mysqld -install安装mysql
5. 启动服务。net start mysql

# 在linux上安装MySQL
1. 检查系统版本
```shell
$ cat /etc/redhat -release
```
2. 安装mysql
```shell
$ yum install mysql
$ yum install mysql-devel
```
3. 安装mariadb
```shell
$ yum install mariadb-server mariadb 
````
mariadb数据库的相关命令是：
```shell
systemctl start mariadb  #启动MariaDB
systemctl stop mariadb  #停止MariaDB
systemctl restart mariadb  #重启MariaDB
systemctl enable mariadb  #设置开机启动
```
```shell
$ mysql -u root -p
```
4. 官网下载安装mysql
```shell
$ wget http://dev.mysql.com/get/mysql-community-release-el7-5.noarch.rpm
$ rpm -ivh mysql-community-release-el7-5.noarch.rpm
$ yum install mysql-community-server
$ service mysqld restart
```
5. 设置密码
```shell
mysql> set password for 'root'@'localhost' =PASSWORD('password');
```