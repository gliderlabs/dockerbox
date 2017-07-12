FROM docker:stable-dind
CMD ["/usr/local/bin/dockerd", "-H", "tcp://0.0.0.0:2375", "--userns-remap=default"]
