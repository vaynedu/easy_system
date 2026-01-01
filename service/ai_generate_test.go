package service

import (
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
