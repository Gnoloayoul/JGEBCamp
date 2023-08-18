# k8s集群搭建(ubuntu20)
## 来自群里的大佬

# 集群中所有实例都要操作
```Bash
#修改主机名 然后重启(不修改也行)
vi /etc/hostname
#修改host文件(多节点的情况下要改 所有节点都要改)
vi /etc/hosts
ps:
192.168.10.10  master(主机名)
192.168.10.11  slave(主机名)


#关闭交换分区
swapoff -a
vi /etc/fstab  #注释交换分区挂载点
#/swapfile                                 none            swap    sw              0       0

#修改DNS
sudo vi /etc/resolv.conf
nameserver 223.5.5.5  
nameserver 223.6.6.6

sudo apt-get update

#新云机可以忽略
#docker  卸载
dpkg -l | grep docker
sudo apt-get autoremove docker-ce-*


#安装docker
sudo apt-get install -y docker.io
#修改dokcer 配置文件
cat > /etc/docker/daemon.json <<EOF
{
    "registry-mirrors": [
        "https://reg-mirror.qiniu.com",
        "https://docker.mirrors.ustc.edu.cn",
        "https://dockerhub.azk8s.cn",
        "https://hub-mirror.c.163.com",
        "https://registry.docker-cn.com"
    ],
    "exec-opts": [
        "native.cgroupdriver=systemd"
    ]
}
EOF
#重启docker
systemctl restart docker




#设置iptables
cat <<EOF | sudo tee /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
EOF

#配置生效
sudo sysctl --system

#安装k8s 组件
#Update the apt package index and install packages needed to use the Kubernetes apt repository
sudo apt-get update
sudo apt-get install -y apt-transport-https ca-certificates curl
sudo curl -s https://mirrors.aliyun.com/kubernetes/apt/doc/apt-key.gpg | sudo apt-key add -
sudo tee /etc/apt/sources.list.d/kubernetes.list <<-'EOF'
deb https://mirrors.aliyun.com/kubernetes/apt kubernetes-xenial main
EOF

#更新
sudo apt-get update

#查看可以安装的版本
apt-cache madison kubeadm

#默认按照最新版本
sudo apt-get install -y kubelet kubeadm kubectl
sudo apt-mark hold kubelet kubeadm kubectl

#安装指定版本(正在使用)
apt-get install <<package name>>=<<version>> 
ps：
sudo apt-get install -y kubelet=1.23.5-00 kubeadm=1.23.5-00
apt-get install -y kubectl=1.23.5-00



#配置kubelet的cgroup
mkdir -p /etc/sysconfig
cat > /etc/sysconfig/kubelet <<EOF
KUBELET_CGROUP_ARGS="--cgroup-driver=systemd"
KUBE_PROXY_MODE="ipvs"
EOF

# 重新指定pause的版本
cat <<EOF | sudo tee /etc/containerd/config.toml
version = 2
[plugins."io.containerd.grpc.v1.cri"]
  sandbox_image = "registry.aliyuncs.com/google_containers/pause:3.9"
EOF

# 使这个指定pause版本操作生效
sudo systemctl restart containerd
sudo systemctl enable containerd

```

# 主实例操作

主节点操作：

```Bash
root@master:~/m8# kubeadm version
kubeadm version: &version.Info{Major:"1", Minor:"23", GitVersion:"v1.23.5", GitCommit:"c285e781331a3785a7f436042c65c5641ce8a9e9", GitTreeState:"clean", BuildDate:"2022-03-16T15:57:37Z", GoVersion:"go1.17.8", Compiler:"gc", Platform:"linux/amd64"}
```

```Bash
root@node:~# ip a
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host 
       valid_lft forever preferred_lft forever
2: ens33: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc fq_codel state UP group default qlen 1000
    link/ether 00:0c:29:c3:b7:60 brd ff:ff:ff:ff:ff:ff
    altname enp2s1
    inet 172.248.190.2/24 brd 172.16.131.255 scope global dynamic noprefixroute ens33
       valid_lft 1158sec preferred_lft 1158sec
    inet6 fe80::5108:e43d:94d6:9936/64 scope link noprefixroute 
       valid_lft forever preferred_lft forever
3: docker0: <NO-CARRIER,BROADCAST,MULTICAST,UP> mtu 1500 qdisc noqueue state DOWN group default 
    link/ether 02:42:91:ef:a3:4c brd ff:ff:ff:ff:ff:ff
    inet 172.17.0.1/16 brd 172.17.255.255 scope global docker0
       valid_lft forever preferred_lft forever

# 给cri解禁[20230817]
sudo vim /etc/containerd/config.toml
## 参考文档 https://blog.csdn.net/weixin_51546928/article/details/131176783
sudo systemctl restart containerd

kubeadm init \
  --apiserver-advertise-address=172.248.190.2 \
  --image-repository registry.aliyuncs.com/google_containers \
  --kubernetes-version=v1.23.5 \
  --service-cidr=10.96.0.0/12 \
  --pod-network-cidr=10.244.0.0/16
  
  
  
--apiserver-advertise-address=172.248.190.2 \    #master ip
--image-repository registry.aliyuncs.com/google_containers \ # 这个是镜像地址，由于国外地址无法访问，故使用的阿里云仓库
--kubernetes-version=v1.23.5 \    # 这个参数是下载的k8s软件版本号  kubectl version
--service-cidr=10.96.0.0/12 \     这个参数后的IP地址直接就套用10.96.0.0/12 ,以后安装时也套用即可，不要更改
--pod-network-cidr=10.244.0.0/16 # k8s内部的pod节点之间网络可以使用的IP段，不能和service-cidr写一样，如果不知道怎么配，就先用这个10.244.0.0/16



mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config



#去污
root@master:~# kubectl describe nodes master  |grep Taints
Taints:             node-role.kubernetes.io/control-plane:NoSchedule

#Taints 内容版本不一样可能由差异 其实就是 输出结果 的 '键' + '-'
root@master:~# kubectl taint nodes master node-role.kubernetes.io/control-plane-
node/master untainted

#获取加入集群的命令
root@master:~/calico# kubeadm token create --print-join-command
kubeadm join xxx.xxx.xxx.xxx:6443 --token aam13h.i33mphjl1yin1j50 --discovery-token-ca-cert-hash sha256:89dc638ee26a77fda1cdf863b0379f810c0afbc19432efa4e1ed284513f4f9fb 

```



# 从实例操作
```Bash
#获取加入集群的命令
kubeadm join xxx.xxx.xxx.xxx:6443 --token aam13h.i33mphjl1yin1j50 --discovery-token-ca-cert-hash sha256:89dc638ee26a77fda1cdf863b0379f810c0afbc19432efa4e1ed284513f4f9fb 

```

# 主实例操作
安装网络插件
```Bash
root@master:~# kubectl get nodes
NAME     STATUS   ROLES                  AGE     VERSION
master   NoReady    control-plane,master   1h   v1.23.5
slave    NoReady    <none>                 1h   v1.23.5


#下载calico.yaml 文件
curl https://docs.projectcalico.org/v3.18/manifests/calico.yaml -O

#修改
vi calico.yaml 找到下面的注释内容(行数可能不一样)
3672             - name: CALICO_IPV4POOL_CIDR
3673               value: "10.244.0.0/16"


#执行
kubectl create -f calico.yaml

#等待调度以及镜像拉取(等到都是Running)
kubectl get pod -A | grep calico
kube-system   calico-kube-controllers-6cfb54c7bb-9qn9j   1/1     Running   0          11m
kube-system   calico-node-98hhf                          1/1     Running   0          11m
kube-system   calico-node-xwtrf                          1/1     Running   0          11m


root@master:~# kubectl get nodes
NAME     STATUS   ROLES                  AGE     VERSION
master   Ready    control-plane,master   1h   v1.23.5
slave    Ready    <none>                 1h   v1.23.5

```

```Bash

#docker 替换成  containerd

#1.停止服务
systemctl stop kubelet
systemctl stop docker
systemctl stop containerd
#2.创建 containerd 配置文件
sudo mkdir -p /etc/containerd
containerd config default | sudo tee /etc/containerd/config.toml
#3.更新配置
vi /etc/containerd/config.toml
sed -i s#k8s.gcr.io/pause:3.5#registry.aliyuncs.com/google_containers/pause:3.5#g /etc/containerd/config.toml
sed -i s#'SystemdCgroup = false'#'SystemdCgroup = true'#g /etc/containerd/config.toml
#4. 修改kubelet  启动参数
vi /etc/systemd/system/kubelet.service.d/10-kubeadm.conf
Environment="KUBELET_EXTRA_ARGS=--container-runtime=remote --container-runtime-endpoint=unix:///run/containerd/containerd.sock --pod-infra-container-image=registry.aliyuncs.com/google_containers/pause:3.5"
#5.重启服务
systemctl daemon-reload
systemctl restart containerd
systemctl restart kubelet
#6. 设置 crictl to set correct endpoint
cat <<EOF | sudo tee /etc/crictl.yaml
runtime-endpoint: unix:///run/containerd/containerd.sock
EOF


#7.添加国内源：
修改/etc/containerd/config.toml 

 66     [plugins."io.containerd.grpc.v1.cri".cni]
 67       bin_dir = "/opt/cni/bin"
 68       conf_dir = "/etc/cni/net.d"
 69       conf_template = ""
 70       max_conf_num = 1

# 添加以下片段
[plugins."io.containerd.grpc.v1.cri".registry.mirrors."docker.io"]
        endpoint = ["https://mirror.ccs.tencentyun.com","http://hub-mirror.c.163.com","http://registry.docker-cn.com","http://docker.mirrors.ustc.edu.cn"]

```
