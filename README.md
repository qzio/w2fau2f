# w2fau2f

Wireguard 2fa u2f

HTTP server in go.
2 factor authentication using fido/u2f security keys.
Successful login opens firewall rules for the connecting ip.

State is stored in files.

## CAUTION: NOT AUDITED

Neither this project, nor dependencies(?) have been audited, please use with caution.

Also, see [#todo](#todo)

## Dependencies

- [u2f](https://github.com/tstranex/u2f) (mit license)
- [go-iptables](https://github.com/coreos/go-iptables/) (apache 2.0 license)

## Usage

The idea is to run `w2fau2f` on your (linux) vpn server node.
Preferable used with [wireguard](https://wireguard.com) but should probably work ok with openvpn or similar.

w2fau2f assumes your setup is as follows:

> client -> internet -> vpn node -> internal network

Network wise:

The client is probably connect to internet using NAT
The vpn node has at least two interfaces,
the interface that is connect to the internal network,
the wireguard interface.

The wireguard tunnel is exposed to the internet through port-forwarding, a third internet connected interface etc.

The wireguard interface uses a /24 network which must not overlap with the internal network.

A client connects to the vpn, which is probably authenticated using the pub/priv key.
To get traffic through the vpn node to the rest of your network, two important iptables rules should be configured.

- Allow forwarding for that ip.
- Enable `NAT` for that ip to the internal network.

This is where w2fau2f comes into play.
You can use the "PostUp" wireguard hook to trigger the client to open a browser tab for `w2fau2f` running on the vpn node.

The client must then login using their u2f key and if successful w2fau2f adds iptables rules to forward traffic from the client's vpn-ip to the internal network using NAT.


```
$ ./w2fau2f --help
Usage of ./w2fau2f:
  -app-id string
    appID (the host) (default "localhost")
  -port int
    the port to listen on (default 3000)
```

U2F/Fido requires tls, w2fau2f needs a cert and a key to work.
`w2fau2f` will look for the key and cert in `./tls.key` and `./tls.crt`.

To generate a self signed certificate and key:

```
$ openssl req -x509 -newkey rsa:1024 -keyout tls.key -nodes -out tls.crt  -days 3650
```

## TODO

This project is not released, and can't really be used yet.

A few issues that's in the way for a 0.1.0 release:

- [ ] Configurable internal net. (It is currently hard coded to 192.168.123.0/24.)
- [ ] Handle the u2f counter properly.
- [ ] Check connectivity, or introduce a timeout or ..? To remove the forward/nat rules from iptables.
- [ ] Rework registration and only allow one registration for each client.
- [ ] Allow multiple registrations, (backup key) where the second one must be authenticated.
- [ ] Package html/js into the binary.
- [ ] Ensure firefox support.
- [ ] Move dependency management to go modules.
- [ ] Refactor html/css/javascript, maybe use vue.js, typescript etc.
- [ ] Cleanup logging and error messages.
- [ ] Various cleanups across the code base.


## Licence

MIT

## Contribute

Feel free to fork and send PR's but this is just a PoC / project for me to learn and play with golang...

