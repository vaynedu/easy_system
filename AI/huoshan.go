package main

import (
	"context"
	"fmt"
	"os"

	"github.com/vaynedu/exam_system/utils"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
)

func Amain() {
	client := arkruntime.NewClientWithApiKey(os.Getenv("ARK_API_KEY"))
	ctx := context.Background()
	req := model.CreateChatCompletionRequest{
		Model: "doubao-1-5-pro-32k-250115",
		Messages: []*model.ChatCompletionMessage{
			&model.ChatCompletionMessage{
				Role: "user",
				Content: &model.ChatCompletionMessageContent{
					StringValue: volcengine.String("生成关于redis缓存雪崩、缓存击穿、缓存穿透的问答题，要包含是什么、如果预防、解决方案 区别。 按照此格式Excel表头：题目类型、题干、选项A、选项B、选项C、选项D、正确答案、答案解析、题目备注、一级分类、二级分类," +
						"其中题型取值：0（选择题）、1（填空题）、2（问答题）；非选择题选项可留空；分类需与系统一致\n" +
						""),
				},
				Name: nil,
			},
		},
	}
	fmt.Println(utils.PrintJsonString(req))
	fmt.Println("----- standard request -----")
	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		fmt.Printf("standard chat error: %v\n", err)
		return
	}
	fmt.Println(utils.PrintJsonString(resp))

	fmt.Println(*resp.Choices[0].Message.Content.StringValue)

}
