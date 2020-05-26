package binding

import (
	"reflect"
	"sync"

	"gopkg.in/go-playground/validator.v9"
)

type defaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

var _ StructValidator = &defaultValidator{}

//-----------------------对接口的实现
func (v *defaultValidator) ValidateStruct(obj interface{}) error {
	if kindOfData(obj) == reflect.Struct {
		v.lazyinit()
		if err := v.validate.Struct(obj); err != nil {
			return err
		}
	}
	return nil
}

//先注册验证器,注册验证器时会初始化内部的validate器
func (v *defaultValidator) RegisterValidation(key string, fn validator.Func) error {
	v.lazyinit() //
	return v.validate.RegisterValidation(key, fn)
}

//-----------------------对接口的实现

func (v *defaultValidator) lazyinit() {
	v.once.Do(func() {
		v.validate = validator.New()
	})
}

//返回入参的结构类型
func kindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()
	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}
