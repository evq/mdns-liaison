# mdns-liaison

mdns-liaison is a shitty hack of a unicast DNS -> multicast mDNS A record request proxy.

It was designed to meet a use case where the ability to do mDNS hostname 
lookups were desired inside docker containers without running a full Avahi stack.

## Usage

```
docker build -t mdns-liaison .
docker run -d --net=host mdns-liaison
```

Add `nameserver 172.17.0.1` to your docker host `/etc/resolv.conf`, Docker
should add this entry into container `/etc/resolv.conf` files automatically.

NOTE: By default this listens on all interfaces.
