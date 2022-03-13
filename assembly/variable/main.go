package main

import (
	"fmt"
	"unsafe"
	pkg "var/pkg"
)

func main() {
	fmt.Println(pkg.Id)
	fmt.Println(pkg.Name)
	fmt.Println(pkg.Num[1])
	fmt.Println(pkg.Helloworld)
	fmt.Println(string(pkg.SliceData))
	fmt.Println(unsafe.Sizeof(pkg.SliceData))
	fmt.Println(unsafe.Sizeof(pkg.Name))
}
