package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/vaynedu/exam_system/config"
	"github.com/vaynedu/exam_system/consts"
	"github.com/vaynedu/exam_system/dao"
	"github.com/vaynedu/exam_system/model"
	"github.com/vaynedu/exam_system/third_part"
)

// GenerateAIQuestionRequest 定义接收参数的结构体
type GenerateAIQuestionRequest struct {
	QuestionType int    `json:"question_type"`
	Tag          string `json:"tag"`
	SecondTag    string `json:"second_tag"`
	Count        int    `json:"count"`
	Requirements string `json:"requirements"`
}

func ValidateGenerateAIQuestionRequest(req *GenerateAIQuestionRequest) error {
	// 1. 标签校验
	if err := validateTagRelation(req.Tag, req.SecondTag); err != nil {
		return fmt.Errorf("无效的标签关系：%s-%s", req.Tag, req.SecondTag)
	}

	// 2. 题型校验
	if req.QuestionType != consts.QuestionTypeShortAnswer {
		return fmt.Errorf("当前仅支持生成问答题")
	}

	// 3. 题目数量校验
	if req.Count <= 0 {
		return fmt.Errorf("无效的题目数量：%d", req.Count)
	}
	if req.Count > 10 {
		return fmt.Errorf("题目数量不能超过10:%d", req.Count)
	}

	// 4.题目描述不能很长
	if len(req.Requirements) > 500 {
		return fmt.Errorf("题目描述不能超过500个字符")
	}

	return nil
}

// GenerateAIQuestionService 生成AI题目服务
func GenerateAIQuestionService(ctx context.Context, questionType int, tag, secondTag string, count int, requirements string) ([]*model.ExamQuestion, error) {

	// 这里可以创建让AI回答输出的模板
	// 比如:  按照此格式Excel表头：题目类型、题干、选项A、选项B、选项C、选项D、正确答案、答案解析、题目备注、一级分类、二级分类; 其中题型取值：0=选择题、1=填空题、2=问答题
	excelDesc := "按照此格式Excel表头:题目类型、题干、选项A、选项B、选项C、选项D、正确答案、答案解析、题目备注、一级分类、二级分类; 其中题型取值:0=选择题、1=填空题、2=问答题"
	questionTypeDesc := fmt.Sprintf("生成题型是%s", consts.GetQuestionTypeName(questionType))
	questionNumDesc := fmt.Sprintf("生成题目数量是%d", count)
	QuestionRemarkDesc := "题目备注：来源、难度、考察点"
	tagDesc := fmt.Sprintf("其中一级分类是%s,二级分类是%s", tag, secondTag)
	requirementsDesc := fmt.Sprintf("题目描述是%s;%s;%s;%s;%s;%s", requirements, excelDesc, questionTypeDesc, questionNumDesc, QuestionRemarkDesc, tagDesc)

	// 调用第三方AI接口
	generatedQuestionContent, err := third_part.NewDouBaoAiService().GetAiGenerateQuestion(ctx, requirementsDesc)
	if err != nil {
		return nil, fmt.Errorf("调用AI接口失败:%w", err)
	}

	// 将返回的markdown内容转化成模型
	questions, err := parseMarkdownTable(generatedQuestionContent)
	if err != nil {
		return nil, fmt.Errorf("解析AI生成题目失败：%w", err)
	}

	// 5. 批量插入数据库
	if len(questions) > 0 {
		if err = dao.NewQuestionDao(config.DB).CreateQuestionsInBatches(questions, 100); err != nil {
			return nil, fmt.Errorf("保存AI生成题目失败：%w", err)
		}
	}

	return questions, nil
}

func parseMarkdownTable(table string) ([]*model.ExamQuestion, error) {
	lines := strings.Split(table, "\n")
	var questions []*model.ExamQuestion

	// 跳过表头行和分隔行
	for i := 2; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		// 分割行内容
		fields := strings.Split(line, "|")
		if len(fields) < 11 {
			continue // 跳过无效行
		}

		// 提取字段（去除首尾空格）
		trimFields := make([]string, 0, 11)
		for j := 1; j < len(fields)-1; j++ { // 跳过首尾空字段
			trimFields = append(trimFields, strings.TrimSpace(fields[j]))
		}

		if len(trimFields) != 11 {
			fmt.Println("字段数量不匹配", trimFields, len(trimFields), fields, len(fields))
			continue // 跳过字段数量不匹配的行
		}

		// 解析题目类型
		qType := int8(0)
		if t, err := strconv.ParseInt(trimFields[0], 10, 8); err == nil {
			qType = int8(t)
		}

		// 创建问题对象
		// |题目类型|题干|选项A|选项B|选项C|选项D|正确答案|答案解析|题目备注|一级分类|二级分类|
		q := &model.ExamQuestion{
			QuestionType:   qType,
			QuestionTitle:  trimFields[1],
			OptionA:        trimFields[2],
			OptionB:        trimFields[3],
			OptionC:        trimFields[4],
			OptionD:        trimFields[5],
			CorrectAnswer:  trimFields[6],
			AnswerAnalysis: trimFields[7],
			QuestionRemark: trimFields[8],
			Tag:            trimFields[9],
			SecondTag:      trimFields[10],
			UploadType:     consts.QuestionImportTypeAiDouBao, // 默认为AI生成
		}

		// 针对tag检查, 这里要使用log代替
		if !IsValidPrimaryTag(q.Tag) {
			fmt.Printf("无效的标签关系：%s-%s", q.Tag, q.SecondTag)
			continue
		}
		if !IsSecondaryOfPrimary(q.Tag, q.SecondTag) {
			fmt.Printf("无效的标签关系：%s-%s", q.Tag, q.SecondTag)
			continue
		}
		if err := validateTagRelation(q.Tag, q.SecondTag); err != nil {
			fmt.Printf("无效的标签关系：%s-%s", q.Tag, q.SecondTag)
			continue
		}

		questions = append(questions, q)
	}

	return questions, nil
}
