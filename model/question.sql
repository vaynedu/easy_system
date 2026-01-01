DROP TABLE IF EXISTS exam_questions;
CREATE TABLE IF NOT EXISTS exam_questions (
                                              id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '题目唯一自增ID',
                                              question_type TINYINT NOT NULL COMMENT '题型：0=选择题，1=填空题，2=问答题',
                                              question_title VARCHAR(500) NOT NULL COMMENT '题干内容',
    option_a VARCHAR(200) DEFAULT '' COMMENT '选项A（仅选择题有效）',
    option_b VARCHAR(200) DEFAULT '' COMMENT '选项B（仅选择题有效）',
    option_c VARCHAR(200) DEFAULT '' COMMENT '选项C（仅选择题有效）',
    option_d VARCHAR(200) DEFAULT '' COMMENT '选项D（仅选择题有效）',
    correct_answer VARCHAR(2000) NOT NULL COMMENT '正确答案（选择题：A/B/C/D；填空/问答：具体答案）',
    answer_analysis VARCHAR(2000) DEFAULT '' COMMENT '答案解析（可选，问答题推荐填写）',
    question_remark VARCHAR(500) DEFAULT '' COMMENT '题目备注（如来源、难度等）',
    tag VARCHAR(50) DEFAULT '' COMMENT '题目一级分类（大类别）',
    second_tag VARCHAR(100) DEFAULT '' COMMENT '题目二级分类（细分类别）',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    update_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE  CURRENT_TIMESTAMP COMMENT '更新时间',
    upload_type TINYINT(1) NOT NULL COMMENT '题目录入方式，默认0=手动 1=excel表格 2=豆包AI 3=阿里AI 4=云雾AI',
    PRIMARY KEY (id)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '多题型通用题库表（选择/填空/问答）';


-- 给现有题库表新增一级分类（tag）、二级分类（second_tag）字段
ALTER TABLE exam_questions
    ADD COLUMN tag VARCHAR(50) DEFAULT '' COMMENT '题目一级分类（大类别）' AFTER question_remark,
ADD COLUMN second_tag VARCHAR(100) DEFAULT '' COMMENT '题目二级分类（细分类别）' AFTER tag;


ALTER TABLE exam_questions MODIFY COLUMN   correct_answer VARCHAR(2000) NOT NULL COMMENT '正确答案（选择题：A/B/C/D；填空/问答：具体答案）';


ALTER TABLE  exam_questions
    ADD COLUMN   update_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE  CURRENT_TIMESTAMP COMMENT '更新时间';

ALTER TABLE  exam_questions
    ADD COLUMN   upload_type TINYINT(1) NOT NULL COMMENT '题目录入方式，默认0=手动 1=excel表格 2=豆包AI 3=阿里AI 4=云雾AI';