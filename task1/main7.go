package main

import "sort"

/* 合并区间
以数组 intervals 表示若干个区间的集合，其中单个区间为 intervals[i] = [starti, endi] 。请你合并所有重叠的区间，并返回 一个不重叠的区间数组，该数组需恰好覆盖输入中的所有区间 。
*/

func Merge(intervals [][]int) [][]int {
	length := len(intervals)
	if length == 0 {
		return [][]int{}
	} else if length == 1 {
		return intervals
	}

	// 排序
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	// 计算
	cur := [][]int{intervals[0]}
	curIndex := 0
	for index := 1; index < length; index++ {
		if cur[curIndex][1] < intervals[index][0] {
			cur = append(cur, intervals[index])
			curIndex++
		} else {
			cur[curIndex][1] = intervals[index][1]
		}
	}

	return cur
}
