FROM golang:1.19

# ARG USER=$USER
# ARG UID=$UID
# ARG GID=$GID
# RUN useradd -m ${USER} --uid=${UID} && echo "${USER}:" chpasswd
# USER ${UID}:${GID}

USER root

ARG ALL_PROXY
ENV all_proxy=$ALL_PROXY
ENV GOPROXY=https://goproxy.cn,direct

ARG KUBE_VERSION


RUN all_proxy=$ALL_PROXY go install k8s.io/code-generator/cmd/defaulter-gen@$KUBE_VERSION
RUN all_proxy=$ALL_PROXY go install k8s.io/code-generator/cmd/client-gen@$KUBE_VERSION
RUN all_proxy=$ALL_PROXY go install k8s.io/code-generator/cmd/lister-gen@$KUBE_VERSION
RUN all_proxy=$ALL_PROXY go install k8s.io/code-generator/cmd/informer-gen@$KUBE_VERSION
RUN all_proxy=$ALL_PROXY go install k8s.io/code-generator/cmd/deepcopy-gen@$KUBE_VERSION
RUN all_proxy=$ALL_PROXY go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.6.2

COPY . /go/src/github.com/traefik/traefik
# COPY . /go/src/github.com/traefik/traefik/v2

RUN mkdir -p $GOPATH/src/k8s.io/code-generator
RUN cp -R $GOPATH/pkg/mod/k8s.io/code-generator@$KUBE_VERSION/* $GOPATH/src/k8s.io/code-generator/
RUN chmod +x $GOPATH/src/k8s.io/code-generator/generate-groups.sh

WORKDIR $GOPATH/src/k8s.io/code-generator
