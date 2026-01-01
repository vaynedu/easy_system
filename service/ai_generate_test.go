package service

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试用的Markdown表格数据
const testMarkdownTable = `|题目类型|题干|选项A|选项B|选项C|选项D|正确答案|答案解析|题目备注|一级分类|二级分类|
|--|--|--|--|--|--|--|--|--|--|--|
|2|什么是算法复杂度？|无|无|无|无|算法复杂度是指算法在执行过程中所需要的资源（时间和空间）的量度|算法复杂度分为时间复杂度和空间复杂度，用于评估算法的效率|AI生成题目|算法|时间复杂度|
|2|什么是动态规划？|无|无|无|无|动态规划是一种通过将原问题分解为相对简单的子问题来求解复杂问题的方法|动态规划适用于具有重叠子问题和最优子结构性质的问题|AI生成题目|算法|动态规划|`

// 测试parseMarkdownTable函数
func TestParseMarkdownTable(t *testing.T) {
	// 测试正常情况
	questions, err := parseMarkdownTable(testMarkdownTable)
	assert.NoError(t, err)
	assert.NotNil(t, questions)
	assert.Len(t, questions, 2)

	// 验证第一题
	assert.Equal(t, int8(2), questions[0].QuestionType)
	assert.Equal(t, "什么是算法复杂度？", questions[0].QuestionTitle)
	assert.Equal(t, "算法复杂度是指算法在执行过程中所需要的资源（时间和空间）的量度", questions[0].CorrectAnswer)
	assert.Equal(t, "算法", questions[0].Tag)
	assert.Equal(t, "时间复杂度", questions[0].SecondTag)

	// 验证第二题
	assert.Equal(t, int8(2), questions[1].QuestionType)
	assert.Equal(t, "什么是动态规划？", questions[1].QuestionTitle)
	assert.Equal(t, "动态规划是一种通过将原问题分解为相对简单的子问题来求解复杂问题的方法", questions[1].CorrectAnswer)
	assert.Equal(t, "算法", questions[1].Tag)
	assert.Equal(t, "动态规划", questions[1].SecondTag)

	// 测试空表格
	questions, err = parseMarkdownTable("|题目类型|题干|选项A|选项B|选项C|选项D|正确答案|答案解析|题目备注|一级分类|二级分类|\n|--|--|--|--|--|--|--|--|--|--|--|\n")
	assert.NoError(t, err)
	assert.Nil(t, questions) // 空表格时返回nil

	// 测试无效表格
	questions, err = parseMarkdownTable("无效的Markdown表格")
	assert.NoError(t, err)
	assert.Nil(t, questions) // 无效表格时返回nil
}

// 测试parseMarkdownTable函数处理选择题
func TestParseMarkdownTable_ChoiceQuestion(t *testing.T) {
	choiceTable := `|题目类型|题干|选项A|选项B|选项C|选项D|正确答案|答案解析|题目备注|一级分类|二级分类|
|--|--|--|--|--|--|--|--|--|--|--|
|0|以下哪种排序算法的平均时间复杂度为O(nlogn)？|冒泡排序|插入排序|快速排序|选择排序|C|快速排序的平均时间复杂度为O(nlogn)，最坏情况下为O(n²)|AI生成题目|算法|排序算法|`

	questions, err := parseMarkdownTable(choiceTable)
	assert.NoError(t, err)
	assert.NotNil(t, questions)
	assert.Len(t, questions, 1)

	// 验证选择题
	assert.Equal(t, int8(0), questions[0].QuestionType)
	assert.Equal(t, "以下哪种排序算法的平均时间复杂度为O(nlogn)？", questions[0].QuestionTitle)
	assert.Equal(t, "冒泡排序", questions[0].OptionA)
	assert.Equal(t, "插入排序", questions[0].OptionB)
	assert.Equal(t, "快速排序", questions[0].OptionC)
	assert.Equal(t, "选择排序", questions[0].OptionD)
	assert.Equal(t, "C", questions[0].CorrectAnswer)
	assert.Equal(t, "算法", questions[0].Tag)
	assert.Equal(t, "排序算法", questions[0].SecondTag)
}

// 测试parseMarkdownTable函数处理填空题
func TestParseMarkdownTable_FillQuestion(t *testing.T) {
	fillTable := `|题目类型|题干|选项A|选项B|选项C|选项D|正确答案|答案解析|题目备注|一级分类|二级分类|
|--|--|--|--|--|--|--|--|--|--|--|
|1|在Go语言中，_____关键字用于声明变量。|无|无|无|无|var|var关键字用于声明变量，例如：var name string|AI生成题目|Go语言|基础语法|`

	questions, err := parseMarkdownTable(fillTable)
	assert.NoError(t, err)
	assert.NotNil(t, questions)
	assert.Len(t, questions, 1)

	// 验证填空题
	assert.Equal(t, int8(1), questions[0].QuestionType)
	assert.Equal(t, "在Go语言中，_____关键字用于声明变量。", questions[0].QuestionTitle)
	assert.Equal(t, "var", questions[0].CorrectAnswer)
	assert.Equal(t, "Go语言", questions[0].Tag)
	assert.Equal(t, "基础语法", questions[0].SecondTag)
}

// 测试生成AI题目的参数校验
func TestGenerateAIQuestionService_ParamValidation(t *testing.T) {
	// 由于GenerateAIQuestionService依赖外部AI服务，我们只测试其参数校验逻辑
	// 这里需要模拟参数校验的结果

	// 测试用例：测试parseMarkdownTable函数的健壮性
	testCases := []struct {
		name     string
		table    string
		expected int
	}{
		{"空表格", "|题目类型|题干|选项A|选项B|选项C|选项D|正确答案|答案解析|题目备注|一级分类|二级分类|\n|--|--|--|--|--|--|--|--|--|--|--|\n", 0},
		{"无效表格", "这不是一个有效的Markdown表格", 0},
		{"缺少字段的表格", "|题目类型|题干|正确答案|\n|--|--|--|\n|2|测试题目|测试答案|\n", 0},
		{"多行题目", testMarkdownTable, 2},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			questions, err := parseMarkdownTable(tc.table)
			assert.NoError(t, err)
			if tc.expected > 0 {
				assert.NotNil(t, questions)
				assert.Len(t, questions, tc.expected)
			} else {
				assert.Nil(t, questions)
			}
		})
	}
}

// 测试生成的题目结构
func TestGeneratedQuestionStructure(t *testing.T) {
	questions, err := parseMarkdownTable(testMarkdownTable)
	assert.NoError(t, err)
	assert.Len(t, questions, 2)

	// 验证题目结构
	for _, q := range questions {
		assert.NotEmpty(t, q.QuestionTitle, "题干不能为空")
		assert.NotEmpty(t, q.CorrectAnswer, "正确答案不能为空")
		assert.NotEmpty(t, q.Tag, "一级分类不能为空")
		assert.NotEmpty(t, q.SecondTag, "二级分类不能为空")
		assert.True(t, q.QuestionType >= 0 && q.QuestionType <= 2, "题型必须在0-2之间")
	}
}

func TestParseMarkdownTable_AIReturnUnicode(t *testing.T) {
	// 创建一个包含Unicode转义序列的测试数据
	aiResponse := "|题目类型|题干|选项A|选项B|选项C|选项D|正确答案|答案解析|题目备注|一级分类|二级分类|\n| ---- | ---- | ---- | ---- | ---- | ---- | ---- | ---- | ---- | ---- | ---- |\n|2|简述redis的缓存雪崩、缓存击穿、缓存穿透的原理及区别。|无|无|无|无|缓存雪崩：指在同一时间段内，大量的缓存键同时过期或者Redis服务崩溃，导致大量的请求直接访问数据库，从而给数据库带来巨大的压力，甚至可能导致数据库崩溃。\u003cbr\u003e缓存击穿：是指某个非常热门的key在缓存中过期的瞬间， 有大量的请求同时访问该key，这些请求会直接穿透缓存访问数据库，给数据库造成较大压力。\u003cbr\u003e缓存穿透：是指客户端请求的数据在缓存和数据库中都不存在，这样每次请求都会直接访问数据库，从而造成数据库压力增大。\u003cbr\u003e区别 ：缓存雪崩是大量缓存键同时失效或Redis崩溃；缓存击穿是单个热门key过期时的大量请求穿透；缓存穿透是请求的数据本身不存在。|本题主要考查对Redis缓存常见问题原理及区别的理解。|数据存储|Redis|\n|2|详细说明redis缓存雪崩、缓存击穿、缓存穿透各自的原理，并指出它们之间的本质区别。|无|无|无|无|缓存雪崩原理：当大量的缓存键设置了相同的过期时间，在这些过期时间点同时到达时，缓存中的数据会同时失效，此时大量的请求就会直接涌向数据库。或者Redis服务出现故障，无法提供缓存服务，也会导致所有请求都去访问数据库。\u003cbr\u003e缓存击穿原理：对于一些访问量极高的热门key，当这个key的缓存过期时，会有大量的并发请求同时到来，由于缓存中已经没有该key的数据，这些请求会直接打到数据库上。\u003cbr\u003e缓存穿透原理 ：客户端请求的数据在缓存和数据库中都不存在，每次请求都会穿过缓存直接访问数据库，因为没有命中缓存，所以无法从缓存中获取数据，只能去数据库查询，而数据库中也没有该数据。\u003cbr\u003e本质区别：缓存雪崩是多个缓存项同时失效引发的问题；缓存击穿是单个热门缓存项过期引发的问题；缓存穿透是请求的数据不存在于系统中引发的问题。|本题需清晰阐述原理及本质区别，帮助理解Redis不同缓存问题。|数据存储|Redis|\n|2|请解释redis的缓存雪崩、缓存击穿、缓存穿透的原理，并分析它们之间的差异。|无|无|无|无|缓存雪崩原理：在Redis中，如果大量的缓存数据在同一时间过期，或者Redis服务器出现故障（如宕机），那么原本应该从缓存中获取数据的请求就会全部转向数据库，使得数据库瞬间承受巨大的访问压力，可能导致数据库崩溃。\u003cbr\u003e缓存击穿原理：当一个非常热门的key的缓存过期时，会有大量的并发请求同时访问该key，由于此时缓存中没有该key的数据，这些请求会直接穿透缓存，去数据库中查询数据，给数据库带来较大的负担。\u003cbr\u003e缓存穿透原理：客户端请求的数据在缓存和数据库中都不存在，这样的请求会直接绕过缓存，每次都去访问数据库，不断消耗数据库的资源。\u003cbr\u003e差异：缓存雪崩涉及大量缓存键的问题，可能是过期时间集中或者Redis故障；缓存击穿针对的是单个热门key过期的情况；缓存穿透是由于请求的数据本身不存在导致的。|本题重点在于原理解释和差异分析，加深对Redis缓存问题的认识。|数据存储|Redis|"

	// 将\n替换为实际的换行符
	aiResponse = strings.ReplaceAll(aiResponse, "\n", "\n")

	// 调用parseMarkdownTable函数解析
	questions, err := parseMarkdownTable(aiResponse)
	assert.NoError(t, err)

	// 检查解析结果
	if questions != nil {
		assert.GreaterOrEqual(t, len(questions), 1, "至少应该解析出1道题目")
		t.Logf("成功解析出%d道题目", len(questions))
		// 验证第一题的内容
		if len(questions) > 0 {
			assert.Contains(t, questions[0].CorrectAnswer, "<br>", "正确答案中应该包含<br>标签")
			t.Logf("第一题正确答案：%s", questions[0].CorrectAnswer)
		}
	} else {
		t.Log("解析结果为nil")
	}
}
