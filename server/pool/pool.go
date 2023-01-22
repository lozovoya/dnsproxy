package pool

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"sync"
)

type ConnectionPoolInterface interface {
	CreatePool(hosts []string) error
	CreateConnection(host string) (*Connection, error)
	ReleaseConnection(id int)
	GetConnection() (*Connection, int, error)
	RefreshConnection(id int) error
}

type Connection struct {
	Host string
	Id   int
	Conn net.Conn
}

type ConnectionPool struct {
	mu       sync.Mutex
	pool     map[int]*Connection
	idleConn map[int]int
}

func NewConnectionPool() ConnectionPoolInterface {
	return &ConnectionPool{
		mu:       sync.Mutex{},
		pool:     make(map[int]*Connection),
		idleConn: make(map[int]int),
	}
}

func (p *ConnectionPool) CreatePool(hosts []string) error {
	for i, h := range hosts {
		connection, err := p.CreateConnection(h)
		if err != nil {
			log.Printf("Error connection, host %v, error %v", h, err)
			continue
		}
		connection.Id = i + 1
		p.mu.Lock()
		p.pool[connection.Id] = connection
		p.idleConn[connection.Id] = 0
		p.mu.Unlock()
	}
	if len(p.pool) == 0 {
		return fmt.Errorf("No available connections")
	}
	log.Printf("Connection pool: %+v", p.pool)
	log.Printf("idle connections: %+v", p.idleConn)
	return nil
}

func (p *ConnectionPool) CreateConnection(host string) (*Connection, error) {
	config := tls.Config{Certificates: nil}
	c, err := tls.Dial("tcp", host, &config)
	if err != nil {
		return nil, err
	}
	var connection = &Connection{
		Host: host,
		Conn: c,
	}
	return connection, nil
}

func (p *ConnectionPool) ReleaseConnection(id int) {
	p.mu.Lock()
	p.idleConn[id] = 0
	p.mu.Unlock()
	log.Printf("Connection id %d is available", id)
}

func (p *ConnectionPool) GetConnection() (*Connection, int, error) {
	connId, err := p.getAvailableConnectionID()
	if err != nil {
		return nil, 0, err
	}
	conn, ok := p.pool[connId]
	if !ok {
		return nil, 0, fmt.Errorf("wrong connection id")
	}
	return conn, connId, nil
}

func (p *ConnectionPool) getAvailableConnectionID() (id int, err error) {
	if len(p.idleConn) == 0 {
		return 0, fmt.Errorf("no available connections")
	}
	var connId int
	p.mu.Lock()
	defer p.mu.Unlock()
	for c, _ := range p.idleConn {
		connId = c
		delete(p.idleConn, c)
		break
	}
	log.Printf("idle connections: %+v", p.idleConn)
	return connId, nil
}

func (p *ConnectionPool) RefreshConnection(id int) error {
	conn, err := p.CreateConnection(p.pool[id].Host)
	if err != nil {
		return err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.pool[id].Conn = conn.Conn
	return nil
}
