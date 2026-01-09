package consts

// PrimaryTag 一级分类结构体
type PrimaryTag struct {
	Name      string   `json:"name"`       // 一级分类名称
	SecondTag []string `json:"second_tag"` // 该一级分类下的二级分类列表
}

// KnowledgeTree 全部知识体系数据（硬编码在此文件，全局可访问）
var KnowledgeTree = []PrimaryTag{
	{
		Name: "算法",
		SecondTag: []string{
			"数组", "链表", "栈", "队列", "哈希表", "堆",
			"二叉树", "二叉搜索树", "平衡树", "线段树", "树状数组", "Trie",
			"图的存储", "图的遍历", "最短路径", "最小生成树", "拓扑排序", "强连通分量",
			"递归与分治", "动态规划", "贪心算法", "回溯法", "滑动窗口", "双指针", "位运算",
			"排序与查找", "数学与数论", "高级/竞赛算法", "并查集", "单调栈/单调队列", "剪枝优化", "前缀和", "差分",
		},
	},
	{
		Name: "系统设计",
		SecondTag: []string{
			"架构设计", "缓存设计", "消息队列", "DDD", "分布式系统", "分布式事务", "分布式锁",
		},
	},
	{
		Name: "数据存储",
		SecondTag: []string{
			"MySQL", "Redis", "MongoDB", "HBase", "Neo4j", "Elasticsearch",
			"分布式文件系统", "云对象存储", "InfluxDB", "Prometheus TSDB", "Hive", "ClickHouse",
		},
	},
	{
		Name: "高频考点",
		SecondTag: []string{
			"计算机网络", "操作系统", "编程语言与框架", "并发与多线程",
			"分布式系统理论", "安全与加密", "DevOps 与部署",
		},
	},
}

// IsValidPrimaryTag 判断一级 tag 是否合法
func IsValidPrimaryTag(tag string) bool {
	for _, pt := range KnowledgeTree {
		if pt.Name == tag {
			return true
		}
	}
	return false
}

// IsValidSecondaryTag 判断二级 tag 是否合法（在整个知识体系中是否存在）
func IsValidSecondaryTag(tag string) bool {
	for _, pt := range KnowledgeTree {
		for _, st := range pt.SecondTag {
			if st == tag {
				return true
			}
		}
	}
	return false
}

// IsSecondaryOfPrimary 判断某二级 tag 是否属于某一级 tag
func IsSecondaryOfPrimary(primary, secondary string) bool {
	for _, pt := range KnowledgeTree {
		if pt.Name == primary {
			for _, st := range pt.SecondTag {
				if st == secondary {
					return true
				}
			}
			return false // 一级分类存在，但二级分类不属于该一级分类
		}
	}
	return false // 一级分类不存在
}
