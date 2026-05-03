FROM golang:1.25 AS builder

COPY . /teacup
WORKDIR /teacup/mayhem

RUN go install github.com/dvyukov/go-fuzz/go-fuzz@latest github.com/dvyukov/go-fuzz/go-fuzz-build@latest
RUN GOFLAGS=-mod=mod go get github.com/dvyukov/go-fuzz/go-fuzz-dep
RUN GOFLAGS=-mod=mod go get github.com/AdaLogics/go-fuzz-headers
RUN apt update && apt install -y clang

RUN go-fuzz-build -libfuzzer -o fuzzteacup.a && \
    clang -fsanitize=fuzzer fuzzteacup.a -o fuzzteacup.libfuzzer

FROM debian:bookworm-slim
COPY --from=builder /teacup/mayhem/fuzzteacup.libfuzzer /
