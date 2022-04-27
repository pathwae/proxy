# Basic example

We want to serve a "Ghost" blog website that listens on port 2368. We use Tumpike to set http://ghost.test.localhost to point on "http://ghost:2368".

We don't use SSL here. That's the simpliest way to do.

If you're using Podman, ensure that you allow rootless container to bind localhost:80 with this comand:

```bash
sudo sysctl -w net.ipv4.ip_unprivileged_port_start=0
```

Then:

```bash
# using docker-compose
docker-compose up

# using podman-compose
podman-compose up
```

Then go to http://ghost.test.localhost - and enjoy !
