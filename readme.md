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

### Q&A

Q: Imagine this proxy being deployed in an infrastructure. What would be the security
concerns you would raise?

A: It is necessary to implement a limited number of simultaneous connections to ip addresses, as well as black-list addresses from which there is too much traffic (although it is better to block with ACL on some network filter).

normally all services of a company (i assume this proxy can be used for internal services, not for internet) are located on the ip addresses of some limited range. appearance of ip addresses which are not included into known pool can mean the attempt to direct traffic to the wrong servers

cache service access must be restricted. record swapping can redirect traffic to rogue servers

Q: How would you integrate that solution in a distributed, microservices-oriented and
containerized architecture?

A: I would use a traffic balancer that redirects requests to a pool of microservices, autoscaling based on service load or response time.

Also when using in a pool, you need to use a single cache for all containers, e.g. based on redis.

CICD must be configured with a canary or blue/green update. 

Q: What other improvements do you think would be interesting to add to the project?

A: export metrics for monitoring

Now the upstream server is selected randomly from the pool, we need to implement a more intelligent selection

creation of a request queue (e.g. based on buffered channels) and a pool of workers for each service. in this way it will be easy to regulate and monitor the maximum number of simultaneous requests that each service handles   