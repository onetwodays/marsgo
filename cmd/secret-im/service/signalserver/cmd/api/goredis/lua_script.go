package goredis



import (
	"errors"
	"io/ioutil"

	"github.com/go-redis/redis"
)

// Lua脚本
type LuaScript struct {
	sha1   string
	client *redis.Client
}

// 创建Lua脚本
func NewLuaScript(client *redis.Client, resource string) (*LuaScript, error) {
	data, err := ioutil.ReadFile(resource)
	if err != nil {
		return nil, err
	}
	l := LuaScript{
		client: client,
	}
	l.sha1, err = l.loadScript(string(data))
	if err != nil {
		return nil, err
	}
	return &l, nil
}

// 执行脚本
func (l *LuaScript) Exec(keys []string, args ...interface{}) (interface{}, error) {
	cmd := l.client.EvalSha(l.sha1, keys, args...)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	return cmd.Val(), nil
}

// 加载脚本
func (l *LuaScript) loadScript(script string) (string, error) {
	var sha1 string
	fn := func(node *redis.Client) error {
		cmd := node.ScriptLoad(script)
		if cmd.Err() != nil {
			return cmd.Err()
		}
		if len(sha1) > 0 && sha1 != cmd.Val() {
			return errors.New("inconsistency checksum")
		}
		sha1 = cmd.Val()
		return nil
	}
	if err := fn(l.client); err != nil {
		return "", err
	}
	return sha1, nil
}

