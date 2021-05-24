```
kali linux 2021.1安装parallels tools踩坑记录

1、加载parallels tools镜像
2、chmod -R 777 整个目录
3、解压tar -zxvf prl_mod.tar.gz
4、修改配置文件
	1、./kmods/prl_fs/SharedFolders/Guest/Linux/prl_fs/inode.c
		开头添加
			#define segment_eq(a, b) (b)
			#define USER_DS 1
	2、./kmods/prl_fs_freeze/Snapshot/Guest/Linux/prl_freeze/prl_fs_freeze.c
		开头添加
			#include <linux/blkdev.h>

	3、./kmods/prl_fs/SharedFolders/Guest/Linux/prl_fs/Makefile
	   ./prl_vid/Video/Guest/Linux/kmod/Makefile
		开头添加
			KBUILD_EXTRA_SYMBOLS := /usr/lib/parallels-tools/kmods/prl_tg/Toolgate/Guest/Linux/prl_tg/Module.symvers

5、重打包
	在kmods目录重打包
		rm prl_mod.tar.gz
		tar -zcvf prl_mod.tar.gz .  dkms.conf Makefile.kmods

参考:https://blog.csdn.net/qq_39563369/article/details/115960130
```
