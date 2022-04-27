# Pathwae - universal reverse proxy for Containers

Pathwae is a **universal** reverse proxy that works for Docker (and docker-compose) and Podman (and podman-compose). The configuration is sent as environment variable.

It supports TLS/SSL, auto-redirect, and it provides a nice UI (and API) to monitor your container stack.

- no need to specify with certificate to use
- you don't habe certificates? Pathwae creates temporary certificates for you

It's simple to use. **Very simple!**

> Before to test, if you want to use "Podman" and/or "podman-compose", please ensure that a standard user can open ports < 1024. This is a simple command to type:
```
sudo systcl -w net.ipv4.ip_unprivileged_port_start=0
```

You want to serve a "Ghost" container? Ghost listens on port 2368, so let's proxify!

```yaml
version: "3"
services:
  blog:
    image: ghost
    environment:
      url: http://blog.site.localhost

  # the pathwae method!
  proxy:
    image: pathwae/pathwae
    environment:
      CONFIG: |
        blog.site.localhost:
          to: http://blog:2368
    ports:
    - 80:80
    - 8080:8080
```

And then visit http://blog.site.localhost - and to see UI of Pathwae to monitor http://localhost:8080

Oh, you want a TLS certificates and make "http" redirected to "https"?:

```yaml
version: "3"
services:
  blog:
    image: ghost
    environment:
      # set scheme to https, that's a ghost prerequisist
      url: https://blog.site.localhost

  # the pathwae method! Add "force_ssl" to "true"
  proxy:
    image: pathwae/pathwae
    environment:
      CONFIG: |
        blog.site.localhost:
          to: http://blog:2368
          force_ssl: true
    ports:
    - 80:80
    - 443:443
    - 8080:8080
```

OK, you created your certificates with `mkcert` like this:
```bash
mkdir -p certs
mkcert -install -cert-file certs/foo.pem -key-file certs/foo.key blog.site.localhost
```

Just mount the `certs` directory in `/certs` and leave Pathwae make the needed task to retrieve the right certificate to use (yes, this is OK with several certificates and backends):

```yaml
version: "3"
services:
  blog:
    image: ghost
    environment:
      # set scheme to https, that's a ghost prerequisist
      url: https://blog.site.localhost

  # the pathwae method! Add "force_ssl" to "true"
  proxy:
    image: pathwae/pathwae
    volumes:
    - ./certs:/certs:ro,z
    environment:
      CONFIG: |
        blog.site.localhost:
          to: http://blog:2368
          force_ssl: true
    ports:
    - 80:80
    - 443:443
    - 8080:8080
```

(You probably need to restart your browser if you tried the previous example)

That's all!

# But... there is Traefik, right?

Traefik is a powerful, famous, complete and probably "better" revers proxy **for production** environment. But...

- It will never work with rootless containers (Podman) because it want to use Docker API to read labels. Pathwae doesn't use labels, API...
- It is a bit more complexe to use, specifically with certificates (you need to write a configuration and mount it to the container + modify the starting command)
- It doesn't create temporary certificates
- And others stuffs...

> Don't think that we are saying here that Traefik is a bad reverse-proxy. It's the opposite of what we think!
> The only thing we defend is that it is mostly adapted to production environments (and especially as an ingress-controller within Kubernetes) but that it cannot fully fulfill the role of a simple reverse-proxy to develop on a working machine. And in particular... with Podman

Also, Traefik's high configuration capacity becomes a disadvantage on a development project in a container stack: there are many labels to add, many configurations to think about. Of course, this should only be done once, but it is sometimes the source of misunderstandings and problems that are hard to find.

Pathwae is not intended to replace Traefik wherever you use it, but to serve as a working base in a local environment without spending time on configuration.

What you certainly want is just to say that such and such address points to such and such container (and port), that's Pathwae!

