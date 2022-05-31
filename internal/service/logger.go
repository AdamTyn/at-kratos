package service

import (
	"at-kratos/api/logger/v1"
	"at-kratos/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

type LoggerService struct {
	*v1.UnimplementedLoggerServiceServer
	uc  *biz.SignupLogUseCase
	log *log.Helper
}

func NewLoggerService(uc *biz.SignupLogUseCase, l log.Logger) (*LoggerService, func()) {
	uc.NewCollector()
	go uc.Collector.Listen()
	return &LoggerService{uc: uc, log: log.NewHelper(l)}, func() {
		uc.Close()
	}
}

func SucSignupGetReply(rows []*v1.SignupLogDTO, total int64) *v1.SignupGetReply {
	return &v1.SignupGetReply{Code: 100, TotalCount: total, Rows: rows}
}

// 查询登录日志
func (s *LoggerService) GetSignupLogs(ctx context.Context, req *v1.SignupLogGetRequest) (*v1.SignupGetReply, error) {
	q := &biz.SignupLogQuery{
		IP:       req.Ip,
		Platform: req.Platform,
		Way:      req.Way,
	}
	if req.Pq != nil {
		q.Page = req.Pq.Page
		q.PageSize = req.Pq.PageSize
	}
	ret, total, err := s.uc.GetSignupLogs(ctx, q)
	if err != nil {
		return &v1.SignupGetReply{Code: 400, Msg: err.Error()}, nil
	}
	data := make([]*v1.SignupLogDTO, len(ret))
	for k := range ret {
		data[k] = &v1.SignupLogDTO{
			Way:          ret[k].Way,
			CompanyUuid:  ret[k].CompanyUUID,
			OperatorUuid: ret[k].OperatorUUID,
			OperatorName: ret[k].OperatorName,
			Extra:        []byte(ret[k].Extra),
			Ip:           ret[k].IP,
			Platform:     ret[k].Platform,
			CreatedTime:  float64(ret[k].CreatedTime.Unix() * 1e3),
		}
	}
	return SucSignupGetReply(data, total), nil
}

// 添加登录日志
func (s *LoggerService) PutSignupLog(ctx context.Context, req *v1.SignupLogPutRequest) (*v1.SignupLogPutReply, error) {
	if !biz.WayOfSignupLog.Has(req.Way) {
		return &v1.SignupLogPutReply{Code: 400, Msg: "Way 无效"}, nil
	}
	do := biz.SignupLogDO{
		CompanyUUID:  req.CompanyUuid,
		OperatorUUID: req.OperatorUuid,
		OperatorName: req.OperatorName,
		Way:          req.Way,
	}
	if req.Extra != nil {
		do.Extra = string(req.Extra)
	}
	if req.Meta != nil {
		do.CreatedTime = time.Now()
		if req.Meta.Timestamp > 0 {
			do.CreatedTime = time.UnixMilli(int64(req.Meta.Timestamp))
		}
		do.IP = req.Meta.Ip
		do.Platform = req.Meta.Platform
	}
	go func(do *biz.SignupLogDO) {
		err := s.uc.Collector.Put(do, 3*time.Second)
		if err != nil {
			log.Error(err)
			_ = s.uc.Collector.Put(do, 5*time.Second)
		}
	}(&do)
	return &v1.SignupLogPutReply{Code: 100}, nil
}
