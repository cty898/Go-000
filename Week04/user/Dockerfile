FROM golang:latest
WORKDIR /root/Go-000/Week04/user
COPY / /root/Go-000/Week04/user
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go build -o user
EXPOSE 10086
ENTRYPOINT ./user -mysql.addr $mysqlAddr -redis.addr $redisAddr