# Centos安装Docker
```
1. 下载docker-ce的repo

curl https://download.docker.com/linux/centos/docker-ce.repo -o /etc/yum.repos.d/docker-ce.repo

2. 安装依赖（这是相比centos7的关键步骤）

yum install https://download.docker.com/linux/fedora/30/x86_64/stable/Packages/containerd.io-1.2.6-3.3.fc30.x86_64.rpm

3. 安装docker-ce

yum install docker-ce

4. 启动docker

  systemctl start docker
  
5. 设置开机自启动

systemctl enable docker
```
# Centos安装docker-compose
```
sudo curl -L "https://github.com/docker/compose/releases/download/1.26.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose

sudo chmod +x /usr/local/bin/docker-compose

sudo ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose

docker-compose --version
```
