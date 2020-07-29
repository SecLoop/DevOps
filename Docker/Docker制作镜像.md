```
制作docker镜像并提交docker镜像到 dockerhub

docker search  搜索镜像

docker pull 到本地

docker cp source 2949a8b25628:/home/		source 为服务器目录  2949a8b25628 为运行的docker id /home/ 为docker内目录

docker commit 2949a8b25628 dedecms:5.7-sp1	保存镜像

docker login 登录账号


docker tag dedecms:5.7sp1  <username>/gelab:dedecms-5.7sp1	修改标签

docker push <username>/gelab:dedecms-5.7sp1			push 到docker hub
```
