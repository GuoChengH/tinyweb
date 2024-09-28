package tinyweb

import (
	"testing"

	"github.com/test-go/testify/assert"
)

// 定义测试用例的结构体
type TestCase struct {
	name     string    // 用例名称
	putPath  string    // 输入到Put的路径
	getPath  string    // 输入到Get的路径
	expected *treeNode // 期望的返回值
}

// 测试前缀树的Put和Get方法
func TestTreeNode(t *testing.T) {
	root := &treeNode{name: "/", children: []*treeNode{}}

	// 定义多个测试用例
	testCases := []TestCase{
		{
			name:     "Simple path insert and get",
			putPath:  "/user/get/:id",
			getPath:  "/user/get/123",
			expected: &treeNode{name: ":id"}, // 期望得到 id 节点
		},
		{
			name:     "Another path insert and get",
			putPath:  "/user/update/:id",
			getPath:  "/user/update/456",
			expected: &treeNode{name: ":id"}, // 期望得到 id 节点
		},
		{
			name:     "Wildcard match",
			putPath:  "/files/*",
			getPath:  "/files/picture.png",
			expected: &treeNode{name: "*"}, // 期望得到通配符 * 节点
		},
		{
			name:     "Double wildcard match",
			putPath:  "/api/**",
			getPath:  "/api/v1/users/list",
			expected: &treeNode{name: "**"}, // 期望得到双星通配符 ** 节点
		},
	}

	// 遍历每个测试用例并执行
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			root.Put(tc.putPath)
			node := root.Get(tc.getPath)

			// 使用 testify 断言
			assert.NotNil(t, node, "节点应当不为空")
			assert.Equal(t, tc.expected.name, node.name, "期望的节点与返回的节点不一致")
		})
	}
}
