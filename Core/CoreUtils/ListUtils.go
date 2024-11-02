package CoreUtils

// AllIntElementsInList 检查列表 A 的所有元素是否都在列表 B 中
// A = [1,2,3]; B = [1,2,3,4]
// AllElementsInList(A, B) = true
func AllIntElementsInList(A, B []int) bool {
	// 创建一个 map 来存储列表 B 的元素
	elementMap := make(map[int]struct{})
	for _, b := range B {
		elementMap[b] = struct{}{}
	}

	// 遍历列表 A，检查每个元素是否在 map 中
	for _, a := range A {
		if _, exists := elementMap[a]; !exists {
			return false // 如果 A 中的某个元素不在 B 中，返回 false
		}
	}

	return true // 所有元素都在 B 中，返回 true
}
func AllStringElementsInList(A, B []string) bool {
	// 创建一个 map 来存储列表 B 的元素
	elementMap := make(map[string]struct{})
	for _, b := range B {
		elementMap[b] = struct{}{}
	}

	// 遍历列表 A，检查每个元素是否在 map 中
	for _, a := range A {
		if _, exists := elementMap[a]; !exists {
			return false // 如果 A 中的某个元素不在 B 中，返回 false
		}
	}

	return true // 所有元素都在 B 中，返回 true
}
