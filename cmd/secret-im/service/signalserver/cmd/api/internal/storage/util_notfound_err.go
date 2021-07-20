package storage
// 数据未找到错误
import (
	"github.com/go-redis/redis"
	"github.com/gocassa/gocassa"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
)

func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	if err == redis.Nil {
		return true
	}
	if err==sqlx.ErrNotFound{
		return true
	}



	_, ok := err.(gocassa.RowNotFoundError)
	if ok {
		return true
	}
	return false
}

