//问题
// 我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，
//是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？

//答：如果上层是直接返回 error 的话，dao 层可以 Wrap 这个 error 抛给上层，这样的在最上层打印错误日志的时候可以输出堆栈信息。
//如果不确定上层的怎么处理 error 还是建议直接返回 error，否则上层再 wrap 一次就 wrap 多次了。

//具体代码如下，拿了个简单的基于 go-kit 搭建的 user 服务进行举例，dao 层 wrap error，service 直接返回 error，最后在 endpoint 层打印错误日志
//服务起起来后(需要先启动 redis 和 mysql)，通过 curl 命令：
//curl -X POST http://localhost:10086/login -H "content-type: application/x-www-form-urlencoded" -d "email=123%40mail.com&password=123"
//查找一个不存在的账户就能返回错误，错误日志如下：
/*
2020/11/30 20:00:03 Login err, request : &{123@mail.com 123}
2020/11/30 20:00:03 err : record not found
SelectByEmail error
github.com/cty898/Go-000/Week02/dao.(*UserDAOImpl).SelectByEmail
	d:/学习资料/Go-000/Week02/dao/user_dao.go:32
github.com/cty898/Go-000/Week02/service.(*UserServiceImpl).Login
	d:/学习资料/Go-000/Week02/service/user_service.go:51
github.com/cty898/Go-000/Week02/endpoint.MakeLoginEndpoint.func1
	d:/学习资料/Go-000/Week02/endpoint/user_endpoint.go:28
github.com/go-kit/kit/transport/http.Server.ServeHTTP
	C:/Users/cty/go/pkg/mod/github.com/go-kit/kit@v0.10.0/transport/http/server.go:121
github.com/gorilla/mux.(*Router).ServeHTTP
	C:/Users/cty/go/pkg/mod/github.com/gorilla/mux@v1.7.4/mux.go:210
net/http.serverHandler.ServeHTTP
	c:/go/src/net/http/server.go:2802
net/http.(*conn).serve
	c:/go/src/net/http/server.go:1890
runtime.goexit
	c:/go/src/runtime/asm_amd64.s:1357
ts=2020-11-30T12:00:03.4582939Z caller=error_handler.go:27 err="SelectByEmail error: record not found"
*/

//具体相关修改如下：
//1、dao/user_dao.go 中第 33 行将 error 进行 wrap
//2、service/user_service.go 中第 63 行直接返回 error
//3、endpoint/user_endpoint.go 中第 29 至 32 行，打印出报错时的 request，采用 %+v 的方式打印出堆栈信息

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/cty898/Go-000/Week02/dao"
	"github.com/cty898/Go-000/Week02/endpoint"
	"github.com/cty898/Go-000/Week02/redis"
	"github.com/cty898/Go-000/Week02/service"
	"github.com/cty898/Go-000/Week02/transport"
)

func main() {

	var (
		// 服务地址和服务名
		servicePort = flag.Int("service.port", 10086, "service port")
	)

	flag.Parse()

	ctx := context.Background()
	errChan := make(chan error)

	err := dao.InitMysql("127.0.0.1", "3306", "root", "123456", "user")
	if err != nil {
		log.Fatal(err)
	}

	err = redis.InitRedis("127.0.0.1", "6379", "")
	if err != nil {
		log.Fatal(err)
	}

	userService := service.MakeUserServiceImpl(&dao.UserDAOImpl{})

	userEndpoints := &endpoint.UserEndpoints{
		endpoint.MakeRegisterEndpoint(userService),
		endpoint.MakeLoginEndpoint(userService),
	}

	r := transport.MakeHttpHandler(ctx, userEndpoints)

	go func() {
		errChan <- http.ListenAndServe(":"+strconv.Itoa(*servicePort), r)
	}()

	go func() {
		// 监控系统信号，等待 ctrl + c 系统信号通知服务关闭
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	error := <-errChan
	log.Println(error)

}
