package render

// JSON common json struct.
// 1.如果一个域不是以大写字母开头的，那么转换成json的时候，这个域是被忽略的。看Field name都是大写
// 2.如果没有使用json:"name"tag，那么输出的json字段名和域名是一样的。
// 3.总结一下，json:"name"格式串是用来指导json.Marshal/Unmarshal，在进行json串和golang对象之间转换的时候映射字段名使用的
//
/*
type Person struct {
    Name string    `json:"name"`
    Age  int       `json:"age" valid:"1-100"`
}
func validateStruct(s interface{}) bool {
  v := reflect.ValueOf(s)

  for i := 0; i < v.NumField(); i++ {
    fieldTag    := v.Type().Field(i).Tag.Get("valid")
    fieldName   := v.Type().Field(i).Name
    fieldType   := v.Field(i).Type()
    fieldValue  := v.Field(i).Interface()

    if fieldTag == "" || fieldTag == "-" {
        continue
    }

    if fieldName == "Age" && fieldType.String() == "int" {
        val := fieldValue.(int)

        tmp := strings.Split(fieldTag, "-")
        var min, max int
        min, _ = strconv.Atoi(tmp[0])
        max, _ = strconv.Atoi(tmp[1])
        if val >= min && val <= max {
            return true
        } else {
            return false
        }
    }
  }
  return true
}
*/
