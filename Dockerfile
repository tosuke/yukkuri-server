# syntax = docker/dockerfile:1.3

FROM golang:1.22 AS dev

WORKDIR /work

RUN --mount=type=secret,id=AQTK1_URL \
    curl -o a1tk1-lnx.tar.gz $(cat /run/secrets/AQTK1_URL) \
    && tar -xzvf a1tk1-lnx.tar.gz

COPY . .

RUN go build -o server
CMD ["go", "run", "main.go"]

FROM gcr.io/distroless/base AS prod

WORKDIR /work
COPY --from=dev /work/server /work/server
COPY --from=dev /work/aqtk1-lnx/lib64 /work/aqtk1-lnx/lib64
CMD ["./server"]
