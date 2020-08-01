# Centos OS
#报错
```
error while loading shared libraries: libcrypto.so.6

error while loading shared libraries: libpcre.so.0

等等类似的...
```
# 解决办法
```
ldd --d /usr/local/nginx/sbin/nginx	查看依赖

yum search libssh 	搜索缺少的依赖

yum install libssh	安装缺少的依赖

yum provides libcrypto.so.6	 查询哪个rpm包,包含 这个lib库

yum install packagename(包名)

find / -name libcrypto (查找文件)

ln -sv /usr/local/lib/libcrypto.so.10 /lib64/libcrypto.so.6 (创建软连接)

```
# Kali OS

```
locate libpcre (查看本地环境是否有libpcre)

apt-cache search libpcre		（检索 包含libpcre包）

需要注意具体软连接的路径

        linux-vdso.so.1 (0x00007ffc032e0000)
        libpthread.so.0 => /lib/x86_64-linux-gnu/libpthread.so.0 (0x00007f2d33cc2000)
        libdl.so.2 => /lib/x86_64-linux-gnu/libdl.so.2 (0x00007f2d33cbd000)
        libpcre.so.1 => not found
        libc.so.6 => /lib/x86_64-linux-gnu/libc.so.6 (0x00007f2d33afa000)
        /lib64/ld-linux-x86-64.so.2 (0x00007f2d33cf9000)


ln -sv /usr/lib/x86_64-linux-gnu/libpcre.so.3.13.3 /lib/x86_64-linux-gnu/libpcre.so.1
```
