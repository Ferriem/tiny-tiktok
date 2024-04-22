package etcd

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type ServiceDiscovery struct {
	cli        *clientv3.Client
	serviceMap sync.Map
	prefix     string
}

func NewServiceDiscovery(endpoints []string) (*ServiceDiscovery, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	return &ServiceDiscovery{
		cli: cli,
	}, nil
}

func (s *ServiceDiscovery) ServiceDiscovery(prefix string) error {
	s.prefix = prefix
	resp, err := s.cli.Get(context.Background(), s.prefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	for _, kv := range resp.Kvs {
		s.SetService(string(kv.Key), string(kv.Value))
	}
	go s.watcher()

	return nil
}

func (s *ServiceDiscovery) watcher() {
	rch := s.cli.Watch(context.Background(), s.prefix, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case mvccpb.PUT:
				s.SetService(string(ev.Kv.Key), string(ev.Kv.Value))
			case mvccpb.DELETE:
				s.DelService(string(ev.Kv.Key))
			}
		}
	}
}

func (s *ServiceDiscovery) SetService(key, val string) {
	s.serviceMap.Store(key, val)
	log.Println("set key:", key, "val:", val)
}

func (s *ServiceDiscovery) DelService(key string) {
	s.serviceMap.Delete(key)
	log.Println("del key:", key)
}

func (s *ServiceDiscovery) GetService(key string) (string, error) {
	if val, ok := s.serviceMap.Load(key); ok {
		return val.(string), nil
	}
	return "", fmt.Errorf("cannot get serviceAddr")
}

func (s *ServiceDiscovery) Close() error {
	return s.cli.Close()
}
