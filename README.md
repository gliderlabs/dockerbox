# dockerbox

Container that runs a Docker daemon configured for running user code.

Currently it runs Docker in Docker with configuration that increases container
isolation. It also adds extra iptables rules and makes it easy to add new IPs
to block via config file.

The architecture is modular so new components can be added to augment the
Docker daemon.

PLEASE CONTRIBUTE by adding any configuration I've missed that will further
isolate/secure containers run by this Docker daemon.

## Run in Docker

```
$ docker run -d -p 12375:2375 --privileged gliderlabs/dockerbox
$ DOCKER_HOST=tcp://127.0.0.1:12375 docker ps
```

## Run on Kubernetes

Should be run as a Daemon Set but feel free to run however. Working manifest
in `run`:

```
$ kubectl apply -f run/dockerbox.yaml
```
Now a headless service is available to use, typically via DNS. A container
running in Kubernetes with a Docker client can do:
```
$ DOCKER_HOST=tcp://dockerbox.default.svc.cluster.local:2375 docker ps
```
