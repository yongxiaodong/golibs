package base

import (
	"fmt"
	"testing"
)

func TestObfuscate(t *testing.T) {
	words := []interface{}{
		"会飞的猪",
		"个",
		"18911110101",
		"1891非1110101",
		"测试",
		"你会飞",
		"",
	}
	for _, word := range words {
		cate := Obfuscate(word.(string))
		fmt.Printf("%s === %s \n", word, cate)
	}

}
