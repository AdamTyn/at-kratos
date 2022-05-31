package entity

import (
	"gitee.com/chunanyong/zorm"
	"time"
)

const TableSignupLog = "signup_log"

type SignupLog struct {
	zorm.EntityStruct
	ID           int64     `column:"id"`
	Way          string    `column:"way"`
	CompanyUUID  string    `column:"company_uuid"`
	OperatorUUID string    `column:"operator_uuid"`
	OperatorName string    `column:"operator_name"`
	IP           string    `column:"ip"`
	Platform     string    `column:"platform"`
	Extra        string    `column:"extra"`
	CreatedTime  time.Time `column:"created_time"`
}

func (entity *SignupLog) GetTableName() string {
	return TableSignupLog
}

func (entity *SignupLog) GetPKColumnName() string {
	return "id"
}
