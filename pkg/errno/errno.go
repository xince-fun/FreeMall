package errno

import (
	"fmt"
	"github.com/xince-fun/FreeMall/kitex_gen/leaf"
)

type ErrNo struct {
	ErrCode int64
	ErrMsg  string
}

type Response struct {
	Code    int64       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (e ErrNo) Error() string {
	return fmt.Sprintf("err_code=%d, err_msg=%s", e.ErrCode, e.ErrMsg)
}

// NewErrNo return ErrNo
func NewErrNo(code int64, msg string) ErrNo {
	return ErrNo{
		ErrCode: code,
		ErrMsg:  msg,
	}
}

func (e ErrNo) WithMessage(msg string) ErrNo {
	e.ErrMsg = msg
	return e
}

var (
	ErrDBFailed               = NewErrNo(int64(leaf.Error_DB_ERROR), "db error")
	ErrTagNotFound            = NewErrNo(int64(leaf.Error_BIZ_TAG_NOT_FOUND), "biz tag not found")
	ErrIDTwoSegmentsAreNull   = NewErrNo(int64(leaf.Error_ID_TWO_SEGMENTS_ARE_NULL), "two segments are null")
	ErrSnowflakeTimeException = NewErrNo(int64(leaf.Error_SNOWFLAKE_TIME_EXCEPTION), "snowflake time callback exception")
)
