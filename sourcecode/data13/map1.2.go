package main

func main() {
	m := map[string]int{
		"a": 1,
	}
	if v, ok := m["a"]; ok { // 判断 key 是否存在。
		println(v, ok)
	}
	println(m["c"])       // 对于不存在的 key，直接返回 \0，不会出错。
	m["b"] = 2            // 新增或修改。
	delete(m, "c")        // 删除。如果 key 不存在，不会出错。
	println(len(m))       // 获取键值对数量。cap 无效。
	for k, v := range m { // 迭代，可仅返回 key。随机顺序返回，每次都不相同。
		println(k, v)
	}
}
