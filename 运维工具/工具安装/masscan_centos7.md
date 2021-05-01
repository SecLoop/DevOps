yum update -y
yum install clang -y
yum install make gcc -y
vim ~/.bash_profile
source ~/.bash_profile
```
# .bash_profile

# Get the aliases and functions
if [ -f ~/.bashrc ]; then
	. ~/.bashrc
fi

# User specific environment and startup programs

PATH=$PATH:$HOME/bin:/usr/bin/masscan/bin

export PATH
```
