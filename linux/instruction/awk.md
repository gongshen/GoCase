awk由模式和操作组成：
# 1、模式
### 模式的类别
* /正则表达式/
* 关系运算符表达式
* 模式匹配表达式
	- 匹配：`~ /正则/`
	- 不匹配：`~ /正则/`
* BEGIN语句块、pattern语句块、END语句块

#### 例子
（1）、输出所有以a开头的用户名和ID
```
shell> awk -F ":" 'BEGIN{printf "%-10s%-10s\n","用户名","用户ID"} /^a/{printf "%-10s%-10s\n",$1,$3}' /etc/passwd
```
（2）、寻找以`/bin/bash`登陆的用户
```
shell> awk '/\/bin\/bash{print $0}' /etc/passwd
```
（3）、匹配中的a出现次数为1到2次
```
shell> awk --posix 'a{1,2}k{print $0}' file1
// 或者
shell> awk --re-interval 'a{1,2}k{print $0}' file1
```
（4）、选取第一次出现a到第一次出现b
```
shell> awk '/a/,/b/{print $0}' file1
```
（5）、选取第3行到第6行
```
shell> awk 'NR>=3 && NR<=6{print $0}' file1
```
（6）、找出`192.168.0.0/16`网段内的主机
```
shell> awk --posix '$2~/192\./168\.[0-9]{1,3}\.[0-9]{1,3}/{print $1,$2}' file2
```