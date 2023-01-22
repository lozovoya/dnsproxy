package app

import (
	"bufio"
	"context"
	"dnproxier/server/cache"
	"dnproxier/server/pool"
	"fmt"
	"golang.org/x/net/dns/dnsmessage"
	"log"
	"net"
	"strconv"
)

type App struct {
	pool       pool.ConnectionPoolInterface
	dnsCache   cache.DNSCacheInterface
	listenPort string
}

func NewApp(listenPort string, hosts []string) (*App, error) {
	connPool := pool.NewConnectionPool()
	err := connPool.CreatePool(hosts)
	if err != nil {
		log.Printf("Creating pool error: %v", err)
		return nil, err
	}
	dnsCache := cache.NewDNSCache()
	return &App{pool: connPool, dnsCache: dnsCache, listenPort: listenPort}, nil
}

func (a *App) ListenAndServe() {
	listener, err := net.Listen("tcp4", a.listenPort)
	if err != nil {
		log.Println(err)
		return
	}
	host, port, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("Listening on: %s, port: %s\n", host, port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go a.handle(conn)
	}
}

func (a *App) handle(inConn net.Conn) {
	ctx := context.Background()
	log.Printf("Start processing request %s", inConn.RemoteAddr().String())
	defer inConn.Close()
	buffer, err := a.readFromConnection(inConn)
	if err != nil {
		log.Printf("Incoming connection error: %v", err)
		return
	}
	name, id, err := a.parseMessage(buffer)
	if err != nil {
		log.Printf("parse error: %v", err)
		return
	}
	var response = make([]byte, 0)
	response, err = a.dnsCache.GetFromCache(ctx, name)
	if err != nil {
		outConn, outConnId, err := a.pool.GetConnection()
		if err != nil {
			log.Printf("no available connections, error: %v", err)
			return
		}
		defer a.pool.ReleaseConnection(outConnId)
		log.Printf("Got connection %s", outConn.Conn.RemoteAddr().String())
		_, err = outConn.Conn.Write(buffer)
		if err != nil {
			log.Printf("Upstream connection error: %v", err)
			log.Printf("Refreshing connection")
			err = a.pool.RefreshConnection(outConnId)
			if err != nil {
				log.Printf("Refreshing error: %v", err)
				return
			}
			log.Printf("Connection refreshed")
			_, err = outConn.Conn.Write(buffer)
			if err != nil {
				log.Printf("Again the error, canceling the process %v", err)
				return
			}
		}
		response, err = a.readFromConnection(outConn.Conn)
		if err != nil {
			log.Printf("Upstream connection error: %v", err)
			return
		}
		err = a.dnsCache.AddToCache(ctx, name, response)
		if err != nil {
			log.Printf("Cache adding error: %v", err)
		}
	} else {
		response, err = a.setResponseID(response, id)

	}
	_, err = inConn.Write(response)
	if err != nil {
		log.Printf("Incoming connection error: %v", err)
		return
	}
	log.Printf("Finish processing request %s", inConn.RemoteAddr().String())
	list, _ := a.dnsCache.ListAllRecords(ctx)
	log.Printf("list of cache records: %+v", list)
	return
}

func (a *App) readFromConnection(conn net.Conn) ([]byte, error) {
	r := bufio.NewReader(conn)
	n, err := r.Peek(2)
	if err != nil {
		return nil, err
	}
	length, err := strconv.ParseInt(fmt.Sprintf("%b%b", n[0], n[1]), 2, 64)
	if err != nil {
		return nil, err
	}
	buffer := make([]byte, 0)
	for i := 0; i < int(length+2); i++ {
		var b byte
		b, err = r.ReadByte()
		if err != nil {
			return nil, err
		}
		buffer = append(buffer, b)
	}
	return buffer, nil
}

func (a *App) parseMessage(buffer []byte) (string, uint16, error) {
	message := dnsmessage.Message{}
	err := message.Unpack(buffer[2:])
	if err != nil {
		return "", 0, err
	}
	return message.Questions[0].Name.String(), message.ID, nil
}

func (a *App) setResponseID(buffer []byte, id uint16) ([]byte, error) {
	message := dnsmessage.Message{}
	err := message.Unpack(buffer[2:])
	if err != nil {
		return nil, err
	}
	message.ID = id
	updatedMessage, err := message.Pack()
	if err != nil {
		return nil, err
	}
	response := buffer[:2]
	response = append(response, updatedMessage...)
	log.Println(response)
	return response, nil
}
