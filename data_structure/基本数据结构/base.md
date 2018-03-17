# 1、**String**
![string](https://github.com/gongshen/GoCase/blob/master/pic/string.png)
# 2、**Make和New**
new返回一个指向已清零的指针，而make返回一个复杂结构。
![make_new](https://github.com/gongshen/GoCase/blob/master/pic/make_new.png)
# 3、**Array**
数组是值类型，赋值和函数传参操作都会复制整个数组（为了避免复制数组，你可以传递一个指向数组的指针，但是数组指针不是数组）
# 4、**Slice**
切片的操作不是复制切片指向的元素，而是创建一个新的切片并复用原来切片的`数组`。**所以一个新的切片修改元素会影响原始切片对应的元素。**
### 切片的复制操作，复制的操作可以由`copy`内置函数替代
```
t:=make([]byte,len(s),(cap(s)+1)*2)
for i:=range s{
	t[i]=s[i]
}
```
![slice](https://github.com/gongshen/GoCase/blob/master/pic/slice.png)
# 5、**Map**
整个hash的存储：

![map](https://github.com/gongshen/GoCase/blob/master/pic/map.png)

注意到Bucket中的key/value存放位置，是将keys放在一起，values放在一起。
**这样对于字节对齐，会节省很多空间。**
