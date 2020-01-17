package grpcsr

import (
	"fmt"
	"net"
	"time"

	"github.com/hashicorp/consul/api"
)

// NewConsulRegister create a new consul register
func NewConsulRegister(addr string, svc string, port int, tag []string, meta map[string]string) (*ConsulRegister, error) {
	config := api.DefaultConfig()
	config.Address = addr
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	ip := localIP()
	deregisterCriticalServiceAfter := time.Duration(1) * time.Minute
	interval := time.Duration(10) * time.Second
	reg := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%v-%v-%v", svc, ip, port), // 服务节点的名称
		Name:    fmt.Sprintf("grpc.health.v1.%v", svc),  // 服务名称
		Tags:    tag,                                    // tag，可以为空
		Port:    port,                                   // 服务端口
		Address: ip,                                     // 服务 IP
		Meta:    meta,
		Check: &api.AgentServiceCheck{ // 健康检查
			Interval:                       interval.String(),                       // 健康检查间隔
			GRPC:                           fmt.Sprintf("%v:%v/%v", ip, port, svc),  // grpc 支持，执行健康检查的地址，service 会传到 Health.Check 函数中
			DeregisterCriticalServiceAfter: deregisterCriticalServiceAfter.String(), // 注销时间，相当于过期时间
		},
	}

	r := &ConsulRegister{
		client:   client,
		register: reg,
	}
	return r, nil
}

// ConsulRegister consul service register
type ConsulRegister struct {
	client   *api.Client
	register *api.AgentServiceRegistration
}

// Register register service
func (r *ConsulRegister) Register() error {
	agent := r.client.Agent()
	return agent.ServiceRegister(r.register)
}

func (r *ConsulRegister) Deregister() {
	agent := r.client.Agent()
	agent.ServiceDeregister(r.register.ID)
}

func localIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
