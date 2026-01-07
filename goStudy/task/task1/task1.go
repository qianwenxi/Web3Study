package main

import (
	"fmt"
	"sort"
	"strconv"
)

// 控制流程-只出现一次的数字
func findOnceNumber(numbers ...int) {
	fmt.Println("=== 控制流程练习-找出只出现一次的数字 ===")
	fmt.Printf("录入数字为：%d\n", numbers)
	var wNumberTimesMap = make(map[int]int)
	for _, num := range numbers {
		wNumberTimesMap[num] = wNumberTimesMap[num] + 1
	}
	// 遍历映射
	var wOnceNumberStr string = ""
	for key, value := range wNumberTimesMap {
		if value == 1 {
			wOnceNumberStr = wOnceNumberStr + strconv.Itoa(key) + ","
		}
	}
	if len(wOnceNumberStr) > 0 {
		wOnceNumberStr = wOnceNumberStr[:len(wOnceNumberStr)-1]
		fmt.Printf("录入数字中，只出现一次的数字有：%s\n", wOnceNumberStr)
		return
	}
	fmt.Println("录入数字中，没有只出现一次的数字！")
}

// 回文数：判断一个整数是否是回文数
func IsPalindrome(pParam uint64) {
	fmt.Println("=== 控制流程练习-判断一个整数是否是回文数 ===")
	fmt.Printf("录入数字为：%d\n", pParam)
	// 边界条件：负数、末尾为0的非零数直接返回false
	if pParam%10 == 0 && pParam != 0 {
		fmt.Println("录入数字非回文数")
	}
	// 反转数字
	var reversedHalf uint64 = 0
	for pParam > reversedHalf {
		// 取出x的最后一位，加入反转后的数字
		reversedHalf = reversedHalf*10 + pParam%10
		// 去掉x的最后一位
		pParam = pParam / 10
	}
	// 两种情况判断：
	// 1. 数字长度为偶数：reversedHalf == x（如1221 → x=12, reversedHalf=12）
	// 2. 数字长度为奇数：reversedHalf/10 == x（如12321 → x=12, reversedHalf=123 → 123/10=12）
	if pParam == reversedHalf || pParam == reversedHalf/10 {
		fmt.Println("录入数字为回文数")
		return
	}
	fmt.Println("录入数字非回文数")
}

// 有效的括号
func isValidBrakets(s string) bool {
	// 奇数长度直接返回 false（括号无法成对）
	if len(s)%2 != 0 {
		return false
	}

	// 右括号对应左括号的映射表
	bracketMap := map[rune]rune{
		')': '(',
		'}': '{',
		']': '[',
	}

	stack := []rune{} // 用切片模拟栈，存储左括号

	for _, char := range s {
		// 若为左括号（map 的值），入栈
		if _, isRight := bracketMap[char]; !isRight {
			// 若为左括号，存入stack
			stack = append(stack, char)
		} else {
			// 若为右括号：检查栈是否为空 或 栈顶左括号不匹配
			if len(stack) == 0 || stack[len(stack)-1] != bracketMap[char] {
				return false
			}
			// 匹配成功，栈顶出栈（切片缩容）
			stack = stack[:len(stack)-1]
		}
	}

	// 所有左括号必须全部匹配闭合（栈为空）
	return len(stack) == 0
}

// 最长公共前缀
func longestCommonPrefix(strs []string) string {
	// 边界条件1：数组为空，返回空字符串（按提示 strs.length >=1，可省略，但兼容极端情况）
	if len(strs) == 0 {
		return ""
	}

	// 以第一个字符串为基准，遍历其每个字符位置
	for i := 0; i < len(strs[0]); i++ {
		// 用基准字符串的第 i 个字符，对比其他所有字符串的第 i 个字符
		for j := 1; j < len(strs); j++ {
			// 终止条件：
			// 1. 某个字符串的长度已达 i（该字符串比基准短，无法继续匹配）
			// 2. 某个字符串的第 i 个字符与基准不匹配
			if i >= len(strs[j]) || strs[j][i] != strs[0][i] {
				// 返回基准字符串的前 i 个字符（0~i-1 是公共前缀）
				return strs[0][:i]
			}
		}
	}

	// 所有字符串都匹配到基准字符串的末尾，说明基准字符串就是最长公共前缀
	return strs[0]
}

// +1
func plusOne(digits []int) []int {
	wLen := len(digits)
	for i := wLen - 1; i >= 0; i-- {
		digits[i]++
		digits[i] %= 10
		if digits[i] != 0 {
			return digits
		}
	}
	return append([]int{1}, digits...)
}

// 引用类型练习-指定有序数组，返回去除重复项后的元素个数
func lengthOfArrNonDuplicates(nums []int) int {
	// 功能前置条件：输入数组必须严格有序
	// 边界条件：数组为空时返回0
	if len(nums) == 0 {
		return 0
	}
	slowIndex := 0
	for fastIndex := 1; fastIndex < len(nums); fastIndex++ {
		if nums[fastIndex] != nums[slowIndex] {
			slowIndex++
			nums[slowIndex] = nums[fastIndex]
		}
		// 重复元素不做处理直接读取下一个元素
	}
	return slowIndex + 1

}

// 合并区间
func mergeIntervals(intervals [][]int) [][]int {
	// 边界条件
	if len(intervals) == 0 {
		return nil
	}

	// 1.按照区间的起始位置升序排列
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})
	// 2.初始化合并后的区间切片，先加入第一个区间
	merged := [][]int{intervals[0]}
	// 遍历剩余期间数据，逐个比较合并或添加
	for _, curr := range intervals[1:] {
		// 获取最近合并后区间数据
		latestMergedIntervals := merged[len(merged)-1]
		//判断当前区间数据是否和上一次处理过的数据重叠
		if latestMergedIntervals[1] >= curr[0] {
			if latestMergedIntervals[1] < curr[1] {
				latestMergedIntervals[1] = curr[1]
			}
		} else {
			merged = append(merged, curr)
		}
	}
	return merged
}

// 已知和找加数
func findAddend(numArr []int, sum int) []int {
	// 定义map  key：元素值 value: 下标
	numMap := make(map[int]int)

	for idx, num := range numArr {
		// 计算需要的互补值
		anotherAddend := sum - num
		// 检查已读取数据中是否已经存在互补值
		if anotherIdx, exists := numMap[anotherAddend]; exists {
			return []int{anotherIdx, idx}
		}
		numMap[num] = idx
	}

	return nil
}

func main() {
	// findOnceNumber(1, 2, 3, 5, 6, 3, 2, 1)
	// IsPalindrome(123545321)
	// fmt.Printf("字符串({[]}[])是有效括号吗：%t\n", isValidBrakets("({[]}])"))

	// fmt.Println("=== 字符串练习-最长公共前缀 ===")
	// testCases4LongestComPrefix := [][]string{
	// 	{"我是谁", "我是谁在", "我是谁在他", "我是谁在你好吗"},
	// 	{"我是谁", "我是", "我"},
	// 	{"我是谁", "我是", "我是他吗", "我是你吗"},
	// }
	// for _, strArr := range testCases4LongestComPrefix {
	// 	result := longestCommonPrefix(strArr)
	// 	fmt.Printf("输入：%v → 输出：%v\n", strArr, result)
	// }

	// fmt.Println("=== 基本值类型练习-加一 ===")
	// testCases4PlusOne := [][]int{
	// 	{1, 2, 3},    // 输出 [1,2,4]
	// 	{4, 3, 2, 1}, // 输出 [4,3,2,2]
	// 	{9},          // 输出 [1,0]
	// 	{9, 9, 9},    // 输出 [1,0,0,0]
	// 	{5, 6, 9, 9}, // 输出 [5,7,0,0]
	// }
	// for _, digits := range testCases4PlusOne {
	// 	result := plusOne(append([]int(nil), digits...)) // 复制切片避免修改原数组
	// 	fmt.Printf("输入：%v → 输出：%v\n", digits, result)
	// }

	// fmt.Println("=== 引用类型练习-指定有序数组，返回去除重复项后的元素个数 ===")
	// testCases4LengthOfArrNonDuplicates := [][]int{
	// 	{1, 2, 2, 3},
	// 	{4, 4, 4, 2, 1},
	// 	{9},
	// 	{},
	// 	{5, 9, 9},
	// }
	// for _, nums := range testCases4LengthOfArrNonDuplicates {
	// 	result := lengthOfArrNonDuplicates(append([]int(nil), nums...)) // 复制切片避免修改原数组
	// 	fmt.Printf("输入：%v → 输出：%d\n", nums, result)
	// }

	// fmt.Println("=== 引用类型练习-合并区间 ===")
	// testCases4MergeIntervals := [][]int{
	// 	{1, 2},
	// 	{4, 5},
	// 	{2, 6},
	// 	{9, 15},
	// }
	// result := mergeIntervals(testCases4MergeIntervals) // 复制切片避免修改原数组
	// fmt.Printf("输入：%v → 输出：%d\n", testCases4MergeIntervals, result)

	fmt.Println("=== 基础练习-已知和找加数 ===")
	nums1 := []int{1, 5, 4, 8, 7, 6, 3, 9, 2}
	target1 := 10
	result := findAddend(nums1, target1) // 复制切片避免修改原数组
	fmt.Printf("输入数组：%v 和为：%d → 输出：%v\n", nums1, target1, result)
}
