Service implements proxying tcp dns requests in DNS-over-TLS. Multi-threading is implemented, the service can handle multiple requests simultaneously.

### Upstsream connection pool

The service uses a pool of DNS servers to redirect incoming requests. The pool is implemented as in-memory storage. Connection to the servers is performed at the moment when the service starts. In case of connection loss, reconnection is performed before using the session.  

### DNS Cache

To minimize the use of upstream servers and accelerate the service, a cache of dns records is implemented. This version uses an in-memory cache. Application of the cache is implemented as an interface, so it can be easily changed to redis or something similar.  For parsing dns messages and making changes in the headers used an external module, I hope the purpose of this task was not to implement parsing. Disabling the cache eliminates the need to use these modules.

### Configuration

the service uses the following parameters:

hosts: list of upstream dns servers
port: local port on which the service receives dns requests

The service is configured using the config.yaml file
configuration example: 
```azure
Hosts:
  - 1.1.1.1:853
  - 8.8.4.4:853
Port: 853
```
The configuration file must be placed in the same folder where the service executable is located (e.g. using persistent volume when running in docker or kubernetes).
### Service running

Running service in docker:

```azure
docker build --tag=dnsproxy .
docker run -it -p 853:853 -v "$(pwd)"/config.yaml:/dnsproxy/config.yaml dnsproxy
```  
