# vm-manager

一个轻量级的单机qemu kvm虚拟机管理器

## 使用说明

请将config-tpl.toml修改后保存为config.toml放在程序的工作目录下  
请自行修改vm-manager-systemd-tpl.service文件后放置于systemd的service目录

### 关于使用ovs-dpdk

请自行配置好dpdk，并且划分好hugepage  
hugepage不足会启不了虚拟机  
hugepage需要挂载在/dev/hugepages  
ovs的bridge现在写死名字为br-dpdk  
连接使用dpdkvhostuserclient方式   
socket全部位于/run/openvswitch/下

ovs-dpdk配置示例(ubuntu 20.04)

```bash
sudo apt-get install openvswitch-switch-dpdk
sudo update-alternatives --set ovs-vswitchd /usr/lib/openvswitch-switch-dpdk/ovs-vswitchd-dpdk
sudo ovs-vsctl set Open_vSwitch . "other_config:dpdk-init=true"
sudo service openvswitch-switch restart
sudo ovs-vsctl add-br br-dpdk -- set bridge br-dpdk datapath_type=netdev
# 配置对外连接
sudo ovs-vsctl add-port br-dpdk dpdk0 -- set Interface dpdk0 type=dpdk  "options:dpdk-devargs=${OVSDEV_PCIID}" 

```


