package main

/* 有效的括号
给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串 s ，判断字符串是否有效。
有效字符串需满足：
左括号必须用相同类型的右括号闭合。
左括号必须以正确的顺序闭合。
每个右括号都有一个对应的相同类型的左括号。
*/

func IsValid(s string) bool {
	stack := []byte{}

	flags := map[byte]byte{
		')': '(',
		']': '[',
		'}': '{',
	}

	for _, v := range s {
		str := byte(v)
		flag, ok := flags[str]
		if !ok {
			stack = append(stack, str)
		} else {
			lFlag := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if lFlag != flag {
				return false
			}
		}
	}

	return len(stack) == 0
}
