# w2fau2f

Wireguard 2fa u2f

HTTP server in go.
2 factor authentication using fido/u2f security keys.
Successful login opens firewall rules for the connecting ip.

State is stored in files.

## CAUTION: NOT AUDITED

Neither this project, nor dependencies(?) have been audited, please use with caution.

Also, see [#todo](#todo)

## Usage

The idea is to run `w2fau2f` on your vpn server node.
Preferrable used with [wireguard](https://wireguard.com) but should probably work ok with openvpn or similar.

w2fau2f assumes your vpn setup is as follows:

> client -> internet -> vpn node -> internal network

and network wise:

The client is probably connect to internet using NAT
The vpn node has at least two interfaces,
the interface that is connect to the internal network,
the wireguard interface.

The wireguard tunnel is exposed to the internet through port-forwarding, a third internet connected interface etc.

the wireguard interface uses a /24 network which must not overlap with the internal network.

A client connects to the vpn, which is probably authenticated using the pub/priv key.
But to get through the vpn node to the rest of your network, two important iptables rules should be configured.

- allow forwarding for that ip
- enable `NAT` for that ip to the internal network


This is where w2fau2f comes into play.
You can use the "PostUp" wireguard hook to trigger the client to open a browser tab for the w2fau2f running on the vpn node.

the client must then login using their u2f key and if successful w2fau2f adds iptables rules to let through traffic from the client's vpn-ip to the internal network.


```
$ ./w2fau2f --help
Usage of ./w2fau2f:
  -app-id string
    appID (the host) (default "localhost")
  -port int
    the port to listen on (default 3000)
```

u2f/fido requires tls, w2fau2f needs a cert and a key to work.

To generate a self signed certificate and key:

```
$ openssl req -x509 -newkey rsa:1024 -keyout tls.key -nodes -out tls.crt  -days 3650
```

## TODO

This project is not released, and can't be used yet.

A few issues that's in the way for a 0.1.0 release:

- [ ] Configurable internal net. (current is hard coded to 192.168.123.0/24)
- [ ] Handle the u2f counter properly.
- [ ] Check connectivity, or introduce a timeout or ..? To remove the forward/nat rules from iptables.
- [ ] Only allow registration once for each client.
- [ ] Allow backup keys. (multiple registrations)
- [ ] Package html/js into the binary.
- [ ] firefox support

Nice to haves includes

- [ ] Refactor html/css/javascript, maybe use vue.js, typescript etc.
- [ ] Cleanup logging and error messages.
- [ ] Various cleanups across the code base.


## Licence

MIT

## Contribute

Feel free to fork and send PR's but this is just a PoC / project for me to learn and play with golang...

