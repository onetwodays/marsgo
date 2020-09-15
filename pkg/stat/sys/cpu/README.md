## stat/sys

System Information

## 项目简介

获取Linux平台下的系统信息，包括cpu主频、cpu使用率等

## 使用
直接调用全局函数即可
```cgo
// ReadStat read cpu stat. stat 保存结果值
func ReadStat(stat *Stat) {
	stat.Usage = atomic.LoadUint64(&usage)
}

// GetInfo get cpu info.
func GetInfo() Info {
	return stats.Info()
```
