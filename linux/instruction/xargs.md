xargs是给其他命令传递参数的一个过滤器，通过xargs的处理，原来的换行和空白会被空格取代。

### 用法：
```
[root@study ~]# cat test.txt
a b c 
d e f
[root@study ~]# cat test.txt | xargs
a b c d e f
```
删除除了a之外的所有文件
```
[root@study ~]# ls | grep -v a |xargs rm -f
```