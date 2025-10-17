package main

/* 最长公共前缀
编写一个函数来查找字符串数组中的最长公共前缀。
如果不存在公共前缀，返回空字符串 ""。
*/

func LongestCommonPrefix(strs []string) string {
	curMaxPrefix := strs[0]
	for _, str := range strs[1:] {
		len1 := len(curMaxPrefix)
		len2 := len(str)
		minLen := min(len1, len2)
		if minLen == 0 {
			return ""
		}
		index := 0
		for index < minLen {
			if curMaxPrefix[index] != str[index] {
				break
			}
			index++
		}
		curMaxPrefix = curMaxPrefix[:index]
	}

	return curMaxPrefix
}
