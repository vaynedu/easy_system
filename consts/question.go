package consts

// 题目录入方式，默认0=手动 1=excel表格 2=豆包AI 3=阿里AI 4=云雾AI
const (
	QuestionImportTypeManual = iota
	QuestionImportTypeExcel
	QuestionImportTypeAiDouBao
	QuestionImportTypeAiAli
	QuestionImportTypeAiYunWu
)
