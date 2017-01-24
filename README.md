#beego+angularjs 实现的基于mesos marathon的docker管理平台
##docker 网络配置
docker network create --driver=bridge -o com.docker.network.bridge.name=dockerbr0 --aux-address "DefaultGatewayIPv4=10.200.70.1" --gateway=10.200.70.172 --subnet 10.200.70.0/24 dockerbr0
通过--aux-address "DefaultGatewayIPv4=10.200.70.1" 设置网关，gateway仍然指向当前host然后host机器上ifconfig配置了dockerbr0网桥
##实现了部分
ip 池管理，机房管理，业务管理，容器扩容
##未完成部分
权限控制，发布
## 效果如下：
Dashboard:
![Dashboard](https://raw.githubusercontent.com/xianyouQ/go-dockermgr/master/introduction/dashboard.png)
容器展示：
![容器管理](https://raw.githubusercontent.com/xianyouQ/go-dockermgr/master/introduction/container.png)
上图所建容器对应marathon的截图：
![上图对应marathon状态的截图](https://raw.githubusercontent.com/xianyouQ/go-dockermgr/master/introduction/marathon.png)
机房管理：
![机房管理](https://raw.githubusercontent.com/xianyouQ/go-dockermgr/master/introduction/idcConf.png)
发布管理(目前只做了一部分)：
![发布管理](https://raw.githubusercontent.com/xianyouQ/go-dockermgr/master/introduction/release.png)
业务管理：
![业务管理](https://raw.githubusercontent.com/xianyouQ/go-dockermgr/master/introduction/yewu.png)
