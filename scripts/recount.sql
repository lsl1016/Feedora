-- recount.sql 一次性校准存量计数漂移（可选执行，不自动运行）。
-- 背景：旧代码发帖时 users.post_count 双重递增、删帖不对称回减、评论删除不回减
-- 用户评论数等缺陷已修复，但历史脏数据不会自愈，需本脚本校准。
-- 适用 MySQL 8.0+；执行前建议备份相关表。

-- 1. 用户发帖数：未删除的已发布帖子（含转发帖）。
UPDATE users u
SET post_count = (
    SELECT COUNT(*) FROM posts p
    WHERE p.author_id = u.id
      AND p.deleted_at IS NULL
      AND p.status = 'published'
);

-- 2. 用户评论数：未删除的正常评论。
UPDATE users u
SET comment_count = (
    SELECT COUNT(*) FROM comments c
    WHERE c.user_id = u.id
      AND c.deleted_at IS NULL
      AND c.status = 'normal'
);

-- 3. 帖子评论数。
UPDATE posts p
SET comment_count = (
    SELECT COUNT(*) FROM comments c
    WHERE c.post_id = p.id
      AND c.deleted_at IS NULL
      AND c.status = 'normal'
);

-- 4. 话题帖子数与参与人数（去重作者）。
UPDATE topics t
SET post_count = (
        SELECT COUNT(*) FROM post_topics pt
        JOIN posts p ON p.id = pt.post_id
            AND p.deleted_at IS NULL
            AND p.status = 'published'
        WHERE pt.topic_id = t.id
    ),
    participant_count = (
        SELECT COUNT(DISTINCT p.author_id) FROM post_topics pt
        JOIN posts p ON p.id = pt.post_id
            AND p.deleted_at IS NULL
            AND p.status = 'published'
        WHERE pt.topic_id = t.id
    );

-- 5. 圈子帖子数。
UPDATE circles ci
SET post_count = (
    SELECT COUNT(*) FROM posts p
    WHERE p.circle_id = ci.id
      AND p.deleted_at IS NULL
      AND p.status = 'published'
);

-- 6. 标签使用数。
UPDATE tags tg
SET use_count = (
    SELECT COUNT(*) FROM post_tags pt
    JOIN posts p ON p.id = pt.post_id
        AND p.deleted_at IS NULL
        AND p.status = 'published'
    WHERE pt.tag_id = tg.id
);

-- 7. 用户等级：每 100 积分升 1 级，最低 1 级（与 growth_repository.AddPointLog 保持一致）。
UPDATE users
SET level = GREATEST(1, FLOOR(point_count / 100) + 1);
