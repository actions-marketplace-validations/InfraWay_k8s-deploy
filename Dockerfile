FROM golang:1.15.5-alpine3.12

ENV KUBCTL_VERSION=1.15.1
RUN wget -q https://storage.googleapis.com/kubernetes-release/release/v1.15.1/bin/linux/amd64/kubectl && \
      chmod u+x kubectl && \
      mv kubectl /usr/local/bin/kubectl && \
      kubectl version --client=true

ENV HELM_VERSION=3.3.0
RUN wget -q https://get.helm.sh/helm-v${HELM_VERSION}-linux-amd64.tar.gz && \
      tar -zxf helm-v${HELM_VERSION}-linux-amd64.tar.gz && \
      mv linux-amd64/helm /usr/local/bin/helm && \
      helm version

WORKDIR /usr/src/app
COPY . /usr/src/app/
RUN go mod download && go build .
ENTRYPOINT ["/usr/src/app/deploy"]