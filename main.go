package main

import (
	"fmt"
	"meta_node/task1"
)

func main() {
	ret1 := task1.SingleNumber([]int{5, 1, 3, 1, 3})
	fmt.Println(ret1)

	ret2 := task1.IsPalindrome(12321)
	fmt.Println(ret2)

	ret3 := task1.IsValid("({}[({})])")
	fmt.Println(ret3)

	ret4 := task1.LongestCommonPrefix([]string{"flower", "flow", "flight"})
	fmt.Println(ret4)

	ret5 := task1.PlusOne([]int{9, 9, 9})
	fmt.Println(ret5)

	arr := []int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}
	ret6 := task1.RemoveDuplicates(arr)
	fmt.Println(ret6, arr)

	ret7 := task1.Merge([][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}})
	fmt.Println(ret7)

	fmt.Println("Hello World.")
}
