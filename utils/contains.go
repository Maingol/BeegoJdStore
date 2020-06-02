package utils

// 自定义一个函数用于方便判断一个切片中是否包含某个元素
func Contains(s []string, subElement string) bool {
	flag := false
	for _, v := range s {
		if v == subElement {
			flag = true
			break
		}
	}
	return flag
}
