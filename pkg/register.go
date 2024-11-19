package etcd

// key（namespace + service_name + lease） = val
// val service_name
//        method
//     service_name
//        method

type ServiceInstance struct {
	Methods   []string `json:"methods"`
	Endpoints []string `json:"endpoint"`

	DeviceId  string `json:"device_id"`
	NetworkId string `json:"network_id"`
}

// 服务注册仅注册服务（节点）不注册服务方法
// 服务发现仅发现服务（节点）不发现方法

// 服务方法由服务自身 通过 缓存服务组件接口 添加到 redis 中

// 网关接收到请求后 先去校验该服务节点是否存在 再去校验该方法是否存在

// ---------------------------- mode-1 -----------------------------
// device1 gateway NetworkId = 1 InternalNetIp = 192.168.1.100:1000 OuterNetIp = 117.34.9.1:1000
// device1 service1 NetworkId = 1 InternalNetIp = 192.168.1.100:1001 OuterNetIp = 117.34.9.1:1001
// device1 service1 NetworkId = 1 InternalNetIp = 192.168.1.100:1002 OuterNetIp = 117.34.9.1:1002
// device1 service1 NetworkId = 1 InternalNetIp = 192.168.1.100:1003 OuterNetIp = 117.34.9.1:1003

// ---------------------------- mode-1 -----------------------------
// device1 gateway NetworkId = 1 InternalNetIp = 192.168.1.100:1000 OuterNetIp = 117.34.9.1:1000
// device2 service1 NetworkId = 1 InternalNetIp = 192.168.1.101:1001 OuterNetIp = 117.34.9.2:1001
// device3 service1 NetworkId = 1 InternalNetIp = 192.168.1.102:1002 OuterNetIp = 117.34.9.3:1002
// device4 service1 NetworkId = 1 InternalNetIp = 192.168.1.103:1003 OuterNetIp = 117.34.9.4:1002

// ---------------------------- mode-2 -----------------------------
// server1 gateway NetworkId = 1 InternalNetIp = 192.168.1.100:1000 OuterNetIp = 117.34.9.1:1000
// server2 service1 NetworkId = 2 InternalNetIp = 192.168.1.101:1001 OuterNetIp = 117.34.9.1:1001
// server3 service1 NetworkId = 3 InternalNetIp = 192.168.1.102:1002 OuterNetIp = 117.34.9.1:1002
// server4 service1 NetworkId = 4 InternalNetIp = 192.168.1.103:1003 OuterNetIp = 117.34.9.1:1003

type Register struct {
}
