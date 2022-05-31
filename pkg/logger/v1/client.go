package v1

import (
	v1 "at-kratos/api/logger/v1"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	transGRPC "github.com/go-kratos/kratos/v2/transport/grpc"
	"google.golang.org/grpc"
	"sync"
	"time"
)

type LoggerClient struct {
	conn *grpc.ClientConn
	cp   *sync.Pool
}

func NewLoggerClient(endpoint string, timeout time.Duration) (*LoggerClient, func(), error) {
	conn, err := transGRPC.DialInsecure(
		context.Background(),
		transGRPC.WithEndpoint(endpoint),
		transGRPC.WithTimeout(timeout),
		transGRPC.WithMiddleware(
			recovery.Recovery(),
			logging.Client(log.DefaultLogger),
		),
	)
	if err != nil {
		return nil, nil, err
	}
	s := &LoggerClient{
		conn: conn,
		cp: &sync.Pool{
			New: func() interface{} {
				return v1.NewLoggerServiceClient(conn)
			},
		},
	}
	return s, s.Close, nil
}

func (c *LoggerClient) Close() {
	_ = c.conn.Close()
}

// 查询登录日志
func (c *LoggerClient) GetSignupLogs(ctx context.Context, req *v1.SignupLogGetRequest) (*v1.SignupGetReply, error) {
	p := c.cp.Get()
	defer c.cp.Put(p)
	client := p.(v1.LoggerServiceClient)
	return client.GetSignupLogs(ctx, req)
}

// 添加登录日志
func (c *LoggerClient) PutSignupLog(ctx context.Context, req *v1.SignupLogPutRequest) (*v1.SignupLogPutReply, error) {
	p := c.cp.Get()
	defer c.cp.Put(p)
	client := p.(v1.LoggerServiceClient)
	return client.PutSignupLog(ctx, req)
}
