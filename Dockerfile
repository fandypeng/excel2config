# faces docker file

FROM golang as build

WORKDIR /var/go/src/

COPY . .

#ENV GO111MODULE off

RUN go build -o face_server cmd/main.go

FROM alpine:latest

MAINTAINER fandy fandypeng@163.com

COPY --from=build /var/go/src/face_server .
COPY --from=build /var/go/src/configs ./configs/

RUN chmod +x face_server

RUN  mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

EXPOSE 8000

CMD ["./face_server", "-conf", "./configs"]


