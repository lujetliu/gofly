package main

import (
	"fmt"
	"github.com/go-playground/validator/v10" // TODO: 源码
)

type RegisterReq struct {
	// gt=0 表示长度必须 > 0, gt=greater than
	Username    string `validate:"gt=0"`
	PasswordNew string `validate:"gt=0"`
	// eqfield 跨字段相等校验
	PasswordRepeat string `validate:"eqfield=PasswordNew"`
	// 合法 email 格式校验
	Email string `validate:"email"`
}

var validate = validator.New()

// 使用校验库 validator 就不需要在每个请求进入业务逻辑之前都编写重复的
// 校验代码了
func (req RegisterReq) valid() error {
	err := validate.Struct(req)
	if err != nil {
		// ...
		return err
	}
	return nil
}

func main() {
	var req = RegisterReq{
		Username:       "Xiaohong",
		PasswordNew:    "phno",
		PasswordRepeat: "ohn",
		Email:          "alex@abc.com",
	}

	err := req.valid()
	if err != nil {
		fmt.Println(err)
	}
}

// 输出:
// Key: 'RegisterReq.PasswordRepeat' Error:Field validation for
// 'PasswordRepeat' failed on the 'eqfield' tag
// TODO: validator 校验库提供的错误信息不够人性化, 针对每种标签定制错误信息

/* 校验器的原理就是使用反射对结构体进行树形遍历 ./reflect.go */
