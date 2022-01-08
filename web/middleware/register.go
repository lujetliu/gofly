package main

import "errors"

/* 请求校验 */

// 假设数据已经通过某个开源库绑定到了具体的结构体上

type RegisterReq struct {
	Username       string `json:"username"`
	PasswordNew    string `json:"password_new"`
	PasswordRepeat string `json:"password_repeat"`
	Email          string `json:"email"`
}

// 箭头型代码
func register(req RegisterReq) error {
	if len(req.Username) > 0 {
		if len(req.PasswordNew) > 0 && len(req.PasswordRepeat) > 0 {
			if req.PasswordNew == req.PasswordRepeat {
				if emailFormatValid(req.Email) {
					createUser()
					return nil
				} else {
					return errors.New("invalid email")
				}
			} else {
				return errors.New("password and reinput must be the same")
			}
		} else {
			return errors.New("password and password reinput must be longer than 0")
		}
	} else {
		return errors.New("length of username cannot be 0")
	}
}

func emailFormatValid(email string) bool {
	// ...
	return false
}

func createUser() {
	// ...
}

// 以上的箭头型代码结构传递给阅读者的消息是: 各个分支有同样的重要性;使用 <重构>
// 一书中提供的卫语句方案(Guard Clauses), 以卫语句取代条件表达式告诉阅读者:
// 这种情况很罕见, 如果它真的发生了, 请做一些必要的处理工作, 然后退出.
/*
	TODO: <重构>
	240(Consolidate Conditional Expression, 合并表达式)
	250(Replace Nested Conditional with Guard Clauses, 以卫语句取代嵌套条件表达式)

	做法:
	- 对于每个检查, 放进一个卫语句
		==> 卫语句要不就从函数中返回, 要不就抛出一个异常.
	- 每次将条件检查替换成卫语句后, 编译并测试
		==> 如果所有卫语句都导致相同结果, 则将这些卫语句合并为一个条件表达式,
		并将这个条件表达式提炼成为一个独立函数.
*/

// 使用卫语句重构 register() 函数
func register1(req RegisterReq) error {
	if len(req.Username) == 0 {
		return errors.New("length of username cannot be 0")
	}

	if len(req.PasswordNew) == 0 || len(req.PasswordRepeat) == 0 {
		return errors.New("password and password reinput must be longer than 0")
	}

	if req.PasswordNew != req.PasswordRepeat {
		return errors.New("password and reinput must be the same")
	}

	if !emailFormatValid(req.Email) {
		return errors.New("invalid email")
	}

	createUser()
	return nil
}

// 重构后代码逻辑和结构更清晰了, 但为了避免为每一个 HTTP 请求都写这么一套差不多
// 的校验函数, 可以选择使用校验器; 从设计的角度讲, 一定需要为每个请求都声明一个
// 结构体, http 请求的校验场景都可以通过请求校验器完成工作, ./validator.go
