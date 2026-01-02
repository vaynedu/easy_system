package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/vaynedu/exam_system/model"
	"github.com/vaynedu/exam_system/utils"
)

// 解析Markdown表格
func parseMarkdownTable(table string) ([]model.ExamQuestion, error) {
	lines := strings.Split(table, "\n")
	var questions []model.ExamQuestion

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
			continue // 跳过字段数量不匹配的行
		}

		// 解析题目类型
		qType := int8(0)
		if t, err := strconv.ParseInt(trimFields[0], 10, 8); err == nil {
			qType = int8(t)
		}

		// 创建问题对象
		q := model.ExamQuestion{
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
		}

		questions = append(questions, q)
	}

	return questions, nil
}

func main() {
	markdownText := `|题目类型|题干|选项A|选项B|选项C|选项D|正确答案|答案解析|题目备注|一级分类（tag）|二级分类（second_tag）|
| ---- | ---- | ---- | ---- | ---- | ---- | ---- | ---- | ---- | ---- | ---- |
|0|以下关于Redis缓存雪崩、缓存击穿、缓存穿透的描述，正确的是（  ）|缓存雪崩是指大量缓存同时失效，导致大量请求直接访问数据库<br>缓存击穿是指一个热点key突然失效，大量请求访问该key<br>缓存穿透是指请求一个不存在的key，导致请求直接访问数据库|缓存雪崩是指请求一个不存在的key，导致请求直接访问数据库<br>缓存击穿是指大量缓存同时失效，导致大量请求直接访问数据库<br>缓存穿透是指一个热点key突然失效，大量请求访问该key|缓存雪崩是指一个热点key突然失效，大量请求访问该key<br>缓存击穿是指请求一个不存在的key，导致请求直接访问数据库<br>缓存穿透是指大量缓存同时失效，导致大量请求直接访问数据库|以上都不对|A|缓存雪崩是指在某一时刻，大量缓存同时过期失效，使得大量请求直接落到数据库上，给数据库带来巨大压力；缓存击穿是指一个非常热门的key，在某个时间点突然失效，此时大量的请求会直接访问数据库；缓存穿透是指客户端请求一个数据库和缓存中都不存在的key，这样每次请求都会直接打到数据库上。所以选项A描述正确。|无|Redis缓存问题|缓存雪崩、击 穿、穿透概念|
|1|Redis缓存雪崩是指大量缓存______，导致大量请求直接访问数据库。| | | | |同时失效|当大量缓存同时失效时，原本应该从缓存中获取数据的请求就会直接访问数据库，这就是缓存雪崩的定义。|无|Redis缓存问题|缓存雪崩概念|
|2|请简述Redis缓存击穿是什么，如何预防以及解决方案。| | | | |缓存击穿是指一个热点key在某个时刻突然失效，此时大量的请求会直接访问数据库，给数据库带来巨大压力。<br><br>预防措施：<br>1. 设置热点key永不过期：对于一些非常热门的key， 不设置过期时间，这样可以避免热点key突然失效的问题。<br>2. 定时更新：在缓存失效前，通过定时任务提前更新缓存。<br><br>解决方案：<br>1. 互斥锁：当发现缓存失效时，先获取一个互斥锁，只有获取到锁的线程才能去数据库查询数据并更新缓存， 其他线程等待。<br>2. 异步更新：当缓存失效时，直接返回旧数据，同时启动一个异步线程去更新缓存。|无|Redis缓存问题|缓存击穿||
|2|请说明Redis缓存穿透是什么，如何预防以及解决方案。| | | | |缓存穿透是指客户端请求一个数据库和缓存中都不存在的key，这样每次请求都会直接打到数据库上，增加数据库的负担。<br><br>预防措施：<br>1. 布隆过滤器：在请求进入缓存之前，先 通过布隆过滤器判断该key是否可能存在，如果不存在则直接返回，避免请求到达数据库。<br>2. 缓存空值：当查询一个不存在的key时，将该key对应的空值也缓存起来，并设置一个较短的过期时间。<br><br>解决方案：<br>1. 对请求进行合法性校验：过滤 掉一些明显不合法的请求，如请求的key格式不符合要求等。<br>2. 监控和报警：实时监控请求情况，当发现大量请求同一个不存在的key时，及时进行报警和处理。|无|Redis缓存问题|缓存穿透||
|2|请阐述Redis缓存雪崩、缓存击穿、缓存穿透的区别。| | | | |1. 概念不同：<br> - 缓存雪崩是指大量缓存同时失效，导致大量请求直接访问数据库。<br> - 缓存击穿是指一个热点key突然失效，大量请求访问该key。<br> - 缓存穿透是指请求一个不存 在的key，导致请求直接访问数据库。<br><br>2. 影响范围不同：<br> - 缓存雪崩影响的是多个缓存key，波及范围较广。<br> - 缓存击穿主要影响的是单个热点key。<br> - 缓存穿透主要是由于请求不存在的key导致的，可能是个别恶意请求或者数据不一致导致。<br><br>3. 解决思路不同：<br> - 缓存雪崩主要通过分散缓存过期时间、使用集群等方式解决。<br> - 缓存击穿主要通过设置热点key永不过期、互斥锁等方式解决。<br> - 缓存穿透主要通过布隆过滤器、缓存空值等方式解决。|无|Redis缓存问题| 缓存雪崩、击穿、穿透区别||`

	data, err := parseMarkdownTable(markdownText)
	if err != nil {
		log.Fatal("处理Markdown表格失败:", err)
	}

	fmt.Println(utils.PrintJsonString(data))
	fmt.Println("所有题目处理完成")
}
