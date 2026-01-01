package consts

// 题目录入方式，默认0=手动 1=excel表格 2=豆包AI 3=阿里AI 4=云雾AI
const (
	QuestionImportTypeManual = iota
	QuestionImportTypeExcel
	QuestionImportTypeAiDouBao
	QuestionImportTypeAiAli
	QuestionImportTypeAiYunWu
)

// 题目类型 0=选择题，1=填空题，2=问答题
const (
	QuestionTypeChoice = iota
	QuestionTypeFillInTheBlank
	QuestionTypeShortAnswer
)

func CheckQuestionType(questionType int) bool {
	switch questionType {
	case QuestionTypeChoice, QuestionTypeFillInTheBlank, QuestionTypeShortAnswer:
		return true
	}
	return false
}

func GetQuestionTypeName(questionType int) string {
	switch questionType {
	case QuestionTypeChoice:
		return "选择题"
	case QuestionTypeFillInTheBlank:
		return "填空题"
	case QuestionTypeShortAnswer:
		return "简答题"
	}
	return "未知"
}
