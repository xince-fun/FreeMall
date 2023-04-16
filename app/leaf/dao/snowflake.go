package dao

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

type SnowflakeIdGenRepoImpl struct {
	cli *clientv3.Client
}

func NewSnowflakeIdGenRepoImpl(cli *clientv3.Client) *SnowflakeIdGenRepoImpl {
	return &SnowflakeIdGenRepoImpl{
		cli: cli,
	}
}

func (s *SnowflakeIdGenRepoImpl) GetPrefixKey(ctx context.Context, prefix string) (*clientv3.GetResponse, error) {
	return s.cli.Get(ctx, prefix, clientv3.WithPrefix())
}

// CreateKeyWithOptLock 事务乐观锁创建
func (s *SnowflakeIdGenRepoImpl) CreateKeyWithOptLock(ctx context.Context, key string, val string) bool {
	_, err := s.cli.Txn(ctx).
		If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).
		Then(clientv3.OpPut(key, val)).
		Commit()
	if err != nil {
		return false
	}
	return true
}

func (s *SnowflakeIdGenRepoImpl) CreateTemporaryKey(ctx context.Context, key string, val string) bool {

	session, err := concurrency.NewSession(s.cli, concurrency.WithTTL(5))
	if err != nil {
		return false
	}

	_, err = s.cli.Put(ctx, key, val, clientv3.WithLease(session.Lease()))
	if err != nil {
		return false
	}

	return true
}

func (s *SnowflakeIdGenRepoImpl) CreateOrUpdateKey(ctx context.Context, key string, val string) bool {

	_, err := s.cli.Put(ctx, key, val)
	if err != nil {
		return false
	}

	return true
}

func (s *SnowflakeIdGenRepoImpl) GetKey(ctx context.Context, key string) (*clientv3.GetResponse, error) {

	return s.cli.Get(ctx, key)
}
