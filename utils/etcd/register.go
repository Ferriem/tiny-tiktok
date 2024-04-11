package etcd

import (
	"context"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type ServiceRegister struct {
	cli     *clientv3.Client
	leaseID clientv3.LeaseID
}

func NewServiceRegister(endpoints []string) (*ServiceRegister, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	srv := &ServiceRegister{
		cli: cli,
	}

	return srv, nil
}

func (s *ServiceRegister) ServiceRegister(serviceName, serviceAddr string, lease int64) error {
	err := s.CreateLease(lease)
	if err != nil {
		return err
	}
	err = s.BindLease(serviceName, serviceAddr)
	if err != nil {
		return err
	}
	err = s.KeepAlive()
	return err
}

func (s *ServiceRegister) CreateLease(lease int64) error {
	resp, err := s.cli.Grant(context.Background(), lease)
	if err != nil {
		return err
	}
	s.leaseID = resp.ID
	return nil
}

func (s *ServiceRegister) BindLease(key, val string) error {
	_, err := s.cli.Put(context.Background(), key, val, clientv3.WithLease(s.leaseID))
	if err != nil {
		return err
	}

	log.Printf("bind lease success %v \n", s.leaseID)
	return nil
}

func (s *ServiceRegister) KeepAlive() error {
	leaseRespChan, err := s.cli.KeepAlive(context.Background(), s.leaseID)
	if err != nil {
		return err
	}

	go func(leaseRespChan <-chan *clientv3.LeaseKeepAliveResponse) {
		for {
			select {
			case resp := <-leaseRespChan:
				log.Println("续约 leaseID:%d", resp.ID)
			}
		}
	}(leaseRespChan)

	return nil
}

func (s *ServiceRegister) Close() error {
	log.Println("close...")
	s.cli.Revoke(context.Background(), s.leaseID)
	return s.cli.Close()
}
