# 0.我是干什么用的
### 运维实用的工具包，目前包含了多个功能，但可能不是很成熟，属于凑合能用的阶段
### 优点: 单二进制文件，无需任何其他依赖,适合在纯内网、不好解决依赖的环境中使用
### 功能列表
- 低配版的ansible
  - shell模块
  - script模块
  - copy模块
- 端口检查工具
  - TCP 端口检查
  - UDP 端口检查
- 批量免密工具
  - 批量免密，完全兼容ansible配置文件
  - 批量取消免密
  - 批量修改密码
- ssl证书工具
  - 证书信息的检查
  - 自签证书生成
  - 基于aliyun的DNS免费SSL证书生成(3个月有效期，可生成泛域名)

### 1. 低配版的ansible
使用方法:
```shell
gotools ansible [HOST-PATTERN] -m [MODULE_NAME] -a [ARGS]
# 其中[HOST-PATTERN]既可以兼容ansible的hosts文件，也支持IP地址、地址段的方式
# -m 指定模块名称
# -a 指定参数
# 具体查看--help

# example1:
# 已做免密
gotools ansible 192.168.1.1 -m shell -a w                    
# 未作免密的情况下，支持使用--password认证
gotools ansible 192.168.1.1 -m shell -a w --password {PASS}  

# example2:
# 当[HOST-PATTERN]为组名时,默认读取/etc/ansible/hosts下的ansible配置文件
gotools ansible group1 -m shell -a w
# 自定义配置文件路径
gotools ansible group1 -m shell -a w --config-file=/path/to/file 

# 配置文件的2种结构,结构1可以指定user和pass
[group1]
ansible_host=172.168.101.71 ansible_port=22 ansible_user=root ansible_ssh_pass=317210
ansible_host=172.168.101.72 ansible_port=22 ansible_user=root ansible_ssh_pass=317210
# 结构2需要已经做过免密
[group2]
172.168.101.71
```
### 2. 端口检查工具
端口检查工具，是为了解决网络策略开通的异常问题,用于测试网络策略或防火墙策略
需要server端和client端
```shell
# 支持TCP和UDP 2种协议,其中UDP必须同时使用server和client
# TCP如果目前端口已经打开,可以仅使用client
gotools port --server --ports=80,443,8080-8099 --protocol tcp
gotools port --client --ports=80,443,8080-8099 --host=127.0.0.1 --protocol udp
```
### 3. 批量免密工具
```shell
# 可以使用IP地址或地址段指定主机,也完全兼容ansible的配置文件,读取配置文件中的组名,默认使用/etc/ansible/hosts
gotools sshkey [HOST-PATTERN] -p={PASS} -u={USER}
gotools sshkey [HOST-PATTERN] -p={PASS} -u={USER} --config-file=/path/to/file 
# example1:
# 批量新增免密
gotools sshkey 192.168.1.1-10 -p={PASS}
# example2:
gotools sshkey group1 -p={PASS}
# example3:
# 批量删除免密
gotools sshkey group1 --delete
# example4:
# 该功能需要先做免密
# 批量修改密码,如果没有指定-p参数,将使用12位随机密码,并保存在当前目录下new_passwd.txt文件中
gotools sshkey group1 --chpasswd
gotools sshkey group1 --chpasswd -p={PASS}
```
### 4. ssl证书工具
#### 4.1 证书检查
```shell
# 支持2种证书检查方式
# 1. 通过域名检查
gotools ssl check -d=baidu.com
# 2. 通过证书文件检查
gotools ssl check -f=cert.pem
```

#### 4.2 自签证书生成
```shell
# -d 可以指定一个或者多个域名或者ip
# -y 指定时间,不指定默认为10年
gotools ssl privite -d=zsops.cn
gotools ssl privite -d=domain1.cn -d=domain2.cn -y=10
```

#### 4.3 ACME证书生成
```shell
# -d 可以指定一个或者多个域名或者ip
# -a 指定阿里云的AK,需要dns相关的权限
# -s 指定阿里云的SK,需要dns相关的权限
gotools ssl acme -d=zsops.cn -a={AliAK} -s={AliSK}
```