#报错
```
error while loading shared libraries: libcrypto.so.6

error while loading shared libraries: libpcre.so.0

等等类似的...
```
#解决办法
```
ldd --d /usr/local/nginx/sbin/nginx	查看依赖

yum search libssh 	搜索缺少的依赖

yum install libssh	安装缺少的依赖

yum provides libcrypto.so.6	 查询哪个rpm包,包含 这个lib库

yum install packagename(包名)

find / -name libcrypto (查找文件)

ln -sv /usr/local/lib/libcrypto.so.10 /lib64/libcrypto.so.6 (创建软连接)

```
