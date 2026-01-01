package third_part

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/vaynedu/exam_system/utils"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
)

type DouBaoAiService struct {
	ApiKey string // 从os获取，防止泄露
	Model  string // 模型名称, 可以调用换模型调用
}

func NewDouBaoAiService() *DouBaoAiService {
	return &DouBaoAiService{
		ApiKey: os.Getenv("ARK_API_KEY"),
		Model:  "doubao-1-5-pro-32k-250115", // 一定要和火山官网的模型名称一致
	}
}

func (d *DouBaoAiService) GetAiGenerateQuestion(ctx context.Context, questionDesc string) (string, error) {
	client := arkruntime.NewClientWithApiKey(d.ApiKey)
	req := model.CreateChatCompletionRequest{
		Model: d.Model,
		Messages: []*model.ChatCompletionMessage{
			&model.ChatCompletionMessage{
				Role: "user",
				Content: &model.ChatCompletionMessageContent{
					StringValue: volcengine.String(questionDesc),
				},
				Name: nil,
			},
		},
	}
	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 || resp.Choices[0].Message.Content == nil {
		return "", errors.New("standard chat Choices is empty or content is nil")
	}

	// 调研日志组件，打印日志，临时先用fmt
	fmt.Println(utils.PrintJsonString(resp))

	return *resp.Choices[0].Message.Content.StringValue, nil
}
