FROM gitpod/workspace-full

# update go
RUN  rm -rf /usr/local/go && tar -C /usr/local -xzf go1.20.2.linux-amd64.tar.gz && export PATH=$PATH:/usr/local/go/bin

# install go watcher
RUN go install github.com/mitranim/gow@latest
