-- 001: 为 user_follows.followee_id 增加单列索引。
-- 背景：粉丝列表按 followee_id 过滤，而现有唯一索引 uk_user_follow(follower_id, followee_id)
-- 以 follower_id 打头，不覆盖该查询，会退化为全表扫描。
ALTER TABLE `user_follows` ADD INDEX `idx_user_follows_followee` (`followee_id`);
