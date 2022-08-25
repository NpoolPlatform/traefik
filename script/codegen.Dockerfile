FROM golang:1.17

# ARG USER=$USER
# ARG UID=$UID
# ARG GID=$GID
# RUN useradd -m ${USER} --uid=${UID} && echo "${USER}:" chpasswd
# USER ${UID}:${GID}

USER root

ARG ALL_PROXY
ENV all_proxy=$ALL_PROXY
ENV GOPROXY=https://proxy.golang.com.cn,direct

ARG KUBE_VERSION

RUN echo $all_proxy
RUN all_proxy=$ALL_PROXY go get k8s.io/code-generator@$KUBE_VERSION; exit 0
RUN all_proxy=$ALL_PROXY go get k8s.io/apimachinery@$KUBE_VERSION; exit 0
RUN all_proxy=$ALL_PROXY go get k8s.io/code-generator/cmd/deepcopy-gen@$KUBE_VERSION; exit 0
RUN all_proxy=$ALL_PROXY go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.6.2; exit 0

# COPY go.mod .
# COPY go.sum .
# RUN GO111MODULE=on GOPROXY=https://goproxy.cn,direct go mod download
COPY . /go/src/github.com/traefik/traefik
COPY . /go/src/github.com/traefik/traefik/v2
RUN ls /go/src/github.com/traefik/traefik/script

RUN mkdir -p $GOPATH/src/k8s.io/{code-generator,apimachinery}
RUN cp -R $GOPATH/pkg/mod/k8s.io/code-generator@$KUBE_VERSION $GOPATH/src/k8s.io/code-generator
RUN cp -R $GOPATH/pkg/mod/k8s.io/apimachinery@$KUBE_VERSION $GOPATH/src/k8s.io/apimachinery
RUN chmod +x $GOPATH/src/k8s.io/code-generator/generate-groups.sh

RUN ls /go/src/github.com/traefik/traefik/script
WORKDIR $GOPATH/src/k8s.io/code-generator
