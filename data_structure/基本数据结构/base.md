# **String**
![string](https://github.com/gongshen/GoCase/blob/master/pic/string.png)
# **Make和New**
new返回一个指向已清零的指针，而make返回一个复杂结构。
![make_new](https://github.com/gongshen/GoCase/blob/master/pic/make_new.png)
# **Slice**
![slice](https://github.com/gongshen/GoCase/blob/master/pic/slice.png)
# **Map**
整个hash的存储：

![map](https://github.com/gongshen/GoCase/blob/master/pic/map.png)

注意到Bucket中的key/value存放位置，是将keys放在一起，values放在一起。
**这样对于字节对齐，会节省很多空间。**
