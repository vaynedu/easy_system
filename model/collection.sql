-- 收藏题目表
CREATE TABLE IF NOT EXISTS `exam_question_collection` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '收藏ID',
  `question_id` int(11) NOT NULL COMMENT '题目ID',
  `tag` varchar(50) NOT NULL COMMENT '一级分类',
  `second_tag` varchar(100) NOT NULL COMMENT '二级分类',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '收藏时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_question_id` (`question_id`),
  KEY `idx_tag` (`tag`),
  KEY `idx_second_tag` (`second_tag`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='收藏题目表';