package grpcsr

import (
	"net"
	"strconv"

	"github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/naming"
)

// NewConsulResolver consul resolver
func NewConsulResolver(address string, service string) naming.Resolver {
	return &consulResolver{
		address: address,
		service: service,
	}
}

type consulResolver struct {
	address string
	service string
}

// Resolve implement
func (r *consulResolver) Resolve(target string) (naming.Watcher, error) {
	config := api.DefaultConfig()
	config.Address = r.address
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &consulWatcher{
		client:  client,
		service: r.service,
		addrs:   map[string]*api.AgentService{},
	}, nil
}

type consulWatcher struct {
	client    *api.Client
	service   string
	addrs     map[string]*api.AgentService
	lastIndex uint64
}

func (w *consulWatcher) Next() ([]*naming.Update, error) {
	for {
		services, metainfo, err := w.client.Health().Service(w.service, "", true, &api.QueryOptions{
			WaitIndex: w.lastIndex, // 同步点，这个调用将一直阻塞，直到有新的更新
		})
		if err != nil {
			logrus.Warn("error retrieving instances from Consul: %v", err)
			continue
		}
		w.lastIndex = metainfo.LastIndex

		addrs := map[string]*api.AgentService{}
		for _, service := range services {
			addrs[net.JoinHostPort(service.Service.Address, strconv.Itoa(service.Service.Port))] = service.Service
		}

		var updates []*naming.Update
		for addr := range w.addrs {
			if service, ok := addrs[addr]; !ok {
				updates = append(updates, &naming.Update{Op: naming.Delete, Addr: addr, Metadata: service})
			}
		}

		for addr := range addrs {
			if service, ok := w.addrs[addr]; !ok {
				updates = append(updates, &naming.Update{Op: naming.Add, Addr: addr, Metadata: service})
			}
		}

		if len(updates) != 0 {
			w.addrs = addrs
			return updates, nil
		}
	}
}

func (w *consulWatcher) Close() {
	// nothing to do
}
