package pkg

// var Id = 9527 // 定义 int  型的包级变量
var Id int

/* 字符串 */
// var Name string = "gopher"
// duplicated definition of symbol var/pkg.Name( ./string_amd64.s) TODO:???
var Name string

/*
	> go tool compile -S pkg.go
	go.cuinfo.packagename. SDWARFINFO dupok size=0
			0x0000 70 6b 67                                         pkg
	go.string."gopher" SRODATA dupok size=6
			0x0000 67 6f 70 68 65 72                                gopher
	"".Name SDATA size=16
			0x0000 00 00 00 00 00 00 00 00 06 00 00 00 00 00 00 00  ................
			rel 0+8 t=1 go.string."gopher"+0
*/

/* 数组类型变量 */
var Num [2]int = [2]int{0, 0}
