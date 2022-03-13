package main

import "fmt"
import pkg "var/pkg"

func main() {
	fmt.Println(pkg.Id)
	fmt.Println(pkg.Name)
	fmt.Println(pkg.Num[1])
}
