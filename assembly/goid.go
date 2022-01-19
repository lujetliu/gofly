package main

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

/*
 */

func GetGoid() int64 {
	var (
		buf [64]byte
		n   = runtime.Stack(buf[:], false)
		stk = strings.TrimPrefix(string(buf[:n]))
	)

	idField := strings.Fields(stk, "goroutine ")[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Errorf("can not get goroutine id: %v", err))
	}

	return int64(id)
}
