package main

import "fmt"

func main() {

	ret1 := SingleNumber([]int{5, 1, 3, 1, 3})
	fmt.Println(ret1)

	ret2 := IsPalindrome(12321)
	fmt.Println(ret2)

	ret3 := IsValid("({}[({})])")
	fmt.Println(ret3)

	ret4 := LongestCommonPrefix([]string{"flower", "flow", "flight"})
	fmt.Println(ret4)

	ret5 := PlusOne([]int{9, 9, 9})
	fmt.Println(ret5)

	arr := []int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}
	ret6 := RemoveDuplicates(arr)
	fmt.Println(ret6, arr)

	ret7 := Merge([][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}})
	fmt.Println(ret7)
}
