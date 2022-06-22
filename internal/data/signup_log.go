package data

import (
	"at-kratos/internal/biz"
	"at-kratos/internal/data/entity"
	"at-kratos/internal/pkg/util"
	"context"
	"gitee.com/chunanyong/zorm"
	"github.com/go-kratos/kratos/v2/log"
	"strings"
)

type signupLogRepo struct {
	data *Data
	log  *log.Helper
}

// NewSignupLogRepo .
func NewSignupLogRepo(data *Data, logger log.Logger) biz.SignupLogRepo {
	return &signupLogRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r signupLogRepo) Get(ctx context.Context, q *biz.SignupLogQuery) ([]*entity.SignupLog, int64, error) {
	rows := make([]*entity.SignupLog, 0)
	finder := zorm.NewSelectFinder(entity.TableSignupLog)
	q.Build(finder)
	page := zorm.NewPage()
	page.PageNo = int(q.Page)
	page.PageSize = int(q.PageSize)
	err := zorm.Query(ctx, finder, &rows, page)
	total := int64(page.TotalCount)
	if err != nil {
		r.log.Errorf("[signupLogRepo::Get] err=%s", err)
		return nil, total, err
	}
	return rows, total, nil
}

func (r signupLogRepo) BatchPut(ctx context.Context, signupLogs []*entity.SignupLog) error {
	if len(signupLogs) < 1 {
		return nil
	}
	rows := make([]zorm.IEntityStruct, len(signupLogs))
	for k := range signupLogs {
		r.makeValid(signupLogs[k])
		rows[k] = signupLogs[k]
	}
	_, err := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
		return zorm.InsertSlice(ctx, rows)
	})
	if err != nil {
		r.log.Errorf("[signupLogRepo::BatchPut] err=%s", err)
	}
	return err
}

func (r signupLogRepo) makeValid(en *entity.SignupLog) {
	if en.OperatorName != "" {
		en.OperatorName = strings.Replace(en.OperatorName, "'", "''", -1)
	}
	if en.Extra == "" {
		en.Extra = "{}"
	} else {
		en.Extra = strings.Replace(en.Extra, "'", "''", -1)
	}
	if !util.IsIPv4(en.IP) && !util.IsIPv6(en.IP) {
		en.IP = "0.0.0.0"
	}
}
