package database

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/IKolyas/otus-highload/config"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Connection struct {
	MasterPool   *pgxpool.Pool
	ReplicaPools []*pgxpool.Pool
	ReplicaIndex int
	ReplicaMutex sync.Mutex
}

var (
	DB Connection
)

// Инициализация пулов соединений
func (p *Connection) Load(c config.Config) error {
	var err error

	masterRaw := "postgres://" + c.DBConfig.User + ":" + c.DBConfig.Password + "@" + c.DBConfig.Host + ":" + c.DBConfig.Port + "/" + c.DBConfig.Dbname

	// Подключение к мастеру
	masterPool, err := pgxpool.Connect(context.Background(), masterRaw)
	if err != nil {
		return err
	}

	p.MasterPool = masterPool

	var replicaUrls []string

	if c.DBConfig.Replicas == 0 {
		replicaUrls = append(replicaUrls, masterRaw)
		pool, err := pgxpool.Connect(context.Background(), replicaUrls[0])
		if err != nil {
			return err
		}
		p.ReplicaPools = append(p.ReplicaPools, pool)
		return nil
	}

	// Подключение к репликам
	for i := range c.DBConfig.Replicas {
		num := strconv.Itoa(i + 1)
		replicaUrls = append(replicaUrls, "postgres://"+c.DBConfig.User+":"+c.DBConfig.Password+"@"+c.DBConfig.Host+"-"+num+":"+c.DBConfig.Port+"/"+c.DBConfig.Dbname)
	}

	for _, url := range replicaUrls {
		pool, err := pgxpool.Connect(context.Background(), url)
		if err != nil {
			return err
		}
		p.ReplicaPools = append(p.ReplicaPools, pool)
	}

	return nil
}

// Получение соединения с репликой (round-robin)
func (p *Connection) getReplicaPool() *pgxpool.Pool {
	p.ReplicaMutex.Lock()
	defer p.ReplicaMutex.Unlock()

	fmt.Println(p.ReplicaIndex)
	pool := p.ReplicaPools[p.ReplicaIndex]
	p.ReplicaIndex = (p.ReplicaIndex + 1) % len(p.ReplicaPools)
	return pool
}

// Выполнение запроса на чтение
func (p *Connection) QueryFromReplica() *pgxpool.Pool {
	return p.getReplicaPool()
}

// Выполнение запроса на запись
func (p *Connection) QueryToMaster() *pgxpool.Pool {
	return p.MasterPool
}
