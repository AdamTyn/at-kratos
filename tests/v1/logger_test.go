package v1

import (
	pb "at-kratos/api/logger/v1"
	cli "at-kratos/pkg/logger/v1"
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

var (
	endpoint = "127.0.0.1:9000"
	timeout  = 10 * time.Second
)

func TestGetSignupLogs(t *testing.T) {
	c, cleanup, err0 := cli.NewLoggerClient(endpoint, timeout)
	if err0 != nil {
		t.Fatal(err0)
	}
	defer cleanup()
	resp, err1 := c.GetSignupLogs(context.Background(), &pb.SignupLogGetRequest{
		Way: "pc",
		Pq: &pb.PageQuery{
			Page:     1,
			PageSize: 12,
		},
	})
	if err1 != nil {
		t.Fatal(err0)
	}
	fmt.Println(resp)
	fmt.Println(resp.TotalCount)
	if resp.TotalCount > 0 {
		fmt.Println(resp.Rows[0].Extra)
		var m map[string]interface{}
		err2 := json.Unmarshal(resp.Rows[0].Extra, &m)
		if err2 != nil {
			t.Fatal(err2)
		}
		fmt.Println(m)
	}
}

func TestPutSignupLogs(t *testing.T) {
	c, cleanup, err0 := cli.NewLoggerClient(endpoint, timeout)
	if err0 != nil {
		t.Fatal(err0)
	}
	defer cleanup()
	m := map[string]interface{}{
		"before": map[string]string{"name": "stephen"},
		"after":  map[string]string{"name": "james"},
	}
	cnt, _ := json.Marshal(m)
	resp, err1 := c.PutSignupLog(context.Background(), &pb.SignupLogPutRequest{
		Way:          "pc",
		CompanyUuid:  "test",
		OperatorUuid: "test",
		OperatorName: "test",
		Extra:        cnt,
		Meta: &pb.Meta{
			Ip:        "192.168.0.0",
			Platform:  "ubuntu",
			Timestamp: float64(time.Now().Nanosecond() / 1e6),
		},
	})
	if err1 != nil {
		t.Fatal(err0)
	}
	fmt.Println(resp)
}
