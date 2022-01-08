package main

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// 从结构上看, 每一个结构体都可以看成是一棵树, 如下:
type Nested struct {
	Email string `validate:"email"`
}

type T struct {
	Age    int `validate:"eq=10"`
	Nested Nested
}

/*
	TODO: 树(数据结构), 深度优先搜索, 广度优先搜索
                  struct T
				  /       \
				 /         \
				/           \
			   /         struct Nested
			  /           \
		    Age            \
                            \
                           Email

	从字段校验的需求来看, 可以采用深度优先搜索或广度优先搜索对这棵结构体树进行
	遍历.
*/

// 递归的深度优先搜索方式的遍历
func Validate(v interface{}) (bool, string) {
	validateResult := true
	errmsg := "success"
	vt := reflect.TypeOf(v)
	vv := reflect.ValueOf(v)
	for i := 0; i < vv.NumField(); i++ {
		fieldVal := vv.Field(i)
		tagContent := vt.Field(i).Tag.Get("validate")
		// TODO: 如果没有 validate tag?
		k := fieldVal.Kind()

		switch k {
		case reflect.Int:
			val := fieldVal.Int()
			tagValStr := strings.Split(tagContent, "=")
			// TODO: struct tag 如何保证 tag 满足规则
			tagVal, _ := strconv.ParseInt(tagValStr[1], 10, 64)
			if val != tagVal {
				errmsg = "Validate int failed, tag is: " +
					strconv.FormatInt(tagVal, 10)
				validateResult = false
			}
		case reflect.String:
			val := fieldVal.String()
			tagValStr := tagContent // TODO: 为何要定义新的变量?
			switch tagValStr {
			case "email":
				nestedResult := validateEmail(val)
				if !nestedResult {
					errmsg = "Validate mail failed, field val is: " + val
					validateResult = false
				}
				// TODO: 其他类型
			}
		case reflect.Struct:
			// 如果有内嵌的 struct, 则深度优先搜索遍历就是一个递归过程
			valInterface := fieldVal.Interface()
			nestedResult, msg := Validate(valInterface)
			if !nestedResult {
				validateResult = false
				errmsg = msg
			}
			// TODO: 其他类型
		}
	}
	return validateResult, errmsg
}

func validateEmail(input string) bool {
	if pass, _ := regexp.MatchString(
		`^([\w.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`, input, // TODO: 正则表达式
	); pass {
		return true
	}
	return false
}

func main() {
	var a = T{Age: 10, Nested: Nested{Email: "abc@abc.com"}}
	validateResult, errmsg := Validate(a)
	fmt.Println(validateResult, errmsg)
}

// 以上简单的支持了 eq=x 和 email 这两个 tag
// TODO: 支持更多的 tag, 功能的完善和容错

/*
   对结构体校验时大量使用了反射, 而 go 的反射在性能上不够理想, 有时甚至影响程
   序的性能, 但需要对结构体进行大量校验的场景往往出现在 web 服务中, 这里并不是
   程序的性能瓶颈所在, 实际的效果还是要从 pprof (TODO) 中做更精确的判断.

   TODO: Parser
   如果确实因为在校验中使用 reflect 而造成了程序的性能瓶颈, 也可以避免反射: 使
   用 go 内置的 Parser 对源代码进行扫描, 然后根据结构体的定义生成校验代码, 可以
   将所有需要校验的结构体放在单独的包内,

*/
