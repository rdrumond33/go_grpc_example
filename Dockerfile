FROM golang:1.17.8

RUN apt-get update && apt-get install --no-install-recommends -y protobuf-compiler=3.12.4-1 build-essential=12.9 \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* \
    && go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1 \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0 \
    && wget -nv https://github.com/ktr0731/evans/releases/download/v0.10.2/evans_linux_amd64.tar.gz \
    && tar -xzvf evans_linux_amd64.tar.gz \
    && mv evans ../bin \
    && rm evans_linux_amd64.tar.gz

WORKDIR /app

COPY go.* ./

RUN go install ./...

COPY . ./

CMD [ "tail", "-f", "/dev/null" ]
