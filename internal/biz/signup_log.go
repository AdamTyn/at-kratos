package biz

import (
	"at-kratos/internal/data/entity"
	"context"
	"errors"
	"fmt"
	"gitee.com/chunanyong/zorm"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/jinzhu/copier"
	"strings"
	"time"
)

type SignupLogDO struct {
	Way          string    `json:"way"`
	CompanyUUID  string    `json:"company_uuid"`
	OperatorUUID string    `json:"operator_uuid"`
	OperatorName string    `json:"operator_name"`
	IP           string    `json:"ip"`
	Platform     string    `json:"platform"`
	Extra        string    `json:"extra"`
	CreatedTime  time.Time `json:"created_time"`
}

type SignupLogQuery struct {
	Way      string `json:"way"`
	IP       string `json:"ip"`
	Platform string `json:"platform"`
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
}

func (q *SignupLogQuery) Build(tx *zorm.Finder) *zorm.Finder {
	if q != nil {
		keys := make([]string, 0)
		values := make([]string, 0)
		if q.Way != "" {
			keys = append(keys, "way=?")
			values = append(values, q.Way)
		}
		if q.IP != "" {
			keys = append(keys, "ip=?")
			values = append(values, q.IP)
		}
		if q.Platform != "" {
			keys = append(keys, "platform=?")
			values = append(values, q.Platform)
		}
		if len(keys) > 0 {
			where := fmt.Sprintf("WHERE %s", strings.Join(keys, " AND "))
			keys = nil
			tx.Append(where, values)
		}
	}
	return tx
}

type SignupLogRepo interface {
	Get(context.Context, *SignupLogQuery) ([]*entity.SignupLog, int64, error)
	BatchPut(context.Context, []*entity.SignupLog) error
}

type SignupLogUseCase struct {
	repo      SignupLogRepo
	log       *log.Helper
	Collector *SignupLogCollector
}

func NewSignupLogUseCase(repo SignupLogRepo, logger log.Logger) *SignupLogUseCase {
	return &SignupLogUseCase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (uc *SignupLogUseCase) GetSignupLogs(ctx context.Context, q *SignupLogQuery) ([]*SignupLogDO, int64, error) {
	if q.Page < 1 || q.Page > 2000 {
		q.Page = 1
	}
	if q.PageSize < 1 || q.PageSize > 2000 {
		q.PageSize = 15
	}
	rows, total, _ := uc.repo.Get(ctx, q)
	if rows != nil {
		var ret []*SignupLogDO
		_ = copier.Copy(&ret, rows)
		return ret, total, nil
	}
	return nil, 0, errors.New("not found")
}

func (uc *SignupLogUseCase) NewCollector() {
	c := &SignupLogCollector{repo: uc.repo, log: uc.log}
	c.data = make([]*entity.SignupLog, 0, PerSignupLog)
	c.tunnel = make(chan entity.SignupLog, PerSignupLog)
	uc.Collector = c
}

func (uc *SignupLogUseCase) Close() {
	if uc.Collector != nil {
		uc.Collector.Close()
	}
}

type SignupLogCollector struct {
	closed bool
	data   []*entity.SignupLog
	tunnel chan entity.SignupLog
	repo   SignupLogRepo
	log    *log.Helper
}

func (c *SignupLogCollector) Listen() {
	d := time.Duration(2)
	ticker := time.NewTicker(d * time.Second)
	defer ticker.Stop()
	for {
		select {
		case ent := <-c.tunnel:
			if len(c.data) == PerSignupLog-1 {
				c.data = append(c.data, &ent)
				c.flush()
			} else {
				if len(c.data) >= PerSignupLog {
					c.flush()
				}
				c.data = append(c.data, &ent)
			}
		case <-ticker.C:
			c.flush()
		default:
			if c.closed {
				return
			}
		}
	}
}

func (c *SignupLogCollector) flush() {
	if len(c.data) < 1 {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := c.repo.BatchPut(ctx, c.data)
	if err != nil {
		go c.log.Error(err)
	}
	c.data = c.data[0:0:PerSignupLog]
}

func (c *SignupLogCollector) Put(do *SignupLogDO, t time.Duration) error {
	if c.closed {
		return errClosed
	}
	var ent entity.SignupLog
	err := copier.Copy(&ent, do)
	if err != nil {
		return err
	}
	timer := time.NewTimer(t)
	defer timer.Stop()
	select {
	case c.tunnel <- ent:
		return nil
	case <-timer.C:
		return errTimeout
	}
}

func (c *SignupLogCollector) Close() {
	if !c.closed {
		c.closed = true
		close(c.tunnel)
	}
}
