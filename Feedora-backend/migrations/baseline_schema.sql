/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `circle_members` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `circle_id` bigint DEFAULT NULL,
  `user_id` bigint DEFAULT NULL,
  `role` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT 'member',
  `status` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT 'normal',
  `mute_reason` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `muted_until` datetime(3) DEFAULT NULL,
  `joined_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_circle_user` (`circle_id`,`user_id`),
  KEY `idx_circle_role` (`circle_id`,`role`),
  KEY `idx_circle_status` (`circle_id`,`status`),
  KEY `idx_circle_members_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `circles` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `owner_id` bigint DEFAULT NULL,
  `name` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `avatar` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `description` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `category` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `join_type` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT 'direct',
  `post_permission` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT 'all',
  `rules` text COLLATE utf8mb4_unicode_ci,
  `status` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT 'normal',
  `member_count` bigint DEFAULT '0',
  `post_count` bigint DEFAULT '0',
  `featured_count` bigint DEFAULT '0',
  `is_recommended` tinyint(1) DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_circle_name` (`name`),
  KEY `idx_circles_owner_id` (`owner_id`),
  KEY `idx_circles_status` (`status`),
  KEY `idx_circles_is_recommended` (`is_recommended`),
  KEY `idx_circles_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `comment_likes` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `comment_id` bigint DEFAULT NULL,
  `user_id` bigint DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_comment_like_user` (`comment_id`,`user_id`),
  KEY `idx_comment_likes_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `comments` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `post_id` bigint DEFAULT NULL,
  `user_id` bigint DEFAULT NULL,
  `parent_id` bigint DEFAULT '0',
  `root_id` bigint DEFAULT '0',
  `reply_to_user_id` bigint DEFAULT NULL,
  `content` varchar(2000) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `status` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT 'normal',
  `like_count` bigint DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_post_parent` (`post_id`,`parent_id`),
  KEY `idx_post_created` (`post_id`),
  KEY `idx_user_created` (`user_id`),
  KEY `idx_comments_root_id` (`root_id`),
  KEY `idx_comments_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `event_outbox` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `event_id` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `topic` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `event_type` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `aggregate_id` bigint DEFAULT NULL,
  `user_id` bigint DEFAULT '0',
  `trace_id` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `payload` json DEFAULT NULL,
  `status` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT 'pending',
  `retry_count` bigint DEFAULT '0',
  `next_retry_at` datetime(3) DEFAULT NULL,
  `dispatched_at` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_event_id` (`event_id`),
  KEY `idx_status_next_retry` (`status`,`next_retry_at`),
  KEY `idx_event_outbox_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `files` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint DEFAULT NULL,
  `biz_type` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `filename` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `object_key` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `url` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `size` bigint DEFAULT '0',
  `mime_type` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_files_user_id` (`user_id`),
  KEY `idx_files_biz_type` (`biz_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `notifications` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint DEFAULT NULL,
  `actor_id` bigint DEFAULT '0',
  `type` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `title` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `content` varchar(1000) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `target_type` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `target_id` bigint DEFAULT NULL,
  `target_url` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `read_status` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT 'unread',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_read_created` (`user_id`,`read_status`,`created_at`),
  KEY `idx_user_created` (`user_id`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `operation_logs` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `admin_id` bigint DEFAULT NULL,
  `admin_name` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `action` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `target_type` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `target_id` bigint DEFAULT NULL,
  `detail` varchar(1000) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_operation_logs_admin_id` (`admin_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `post_favorites` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `post_id` bigint DEFAULT NULL,
  `user_id` bigint DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_post_fav_user` (`post_id`,`user_id`),
  KEY `idx_post_favorites_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `post_images` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `post_id` bigint DEFAULT NULL,
  `image_url` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `sort_order` bigint DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_post_images_post_id` (`post_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `post_likes` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `post_id` bigint DEFAULT NULL,
  `user_id` bigint DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_post_like_user` (`post_id`,`user_id`),
  KEY `idx_post_likes_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `post_tags` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `post_id` bigint DEFAULT NULL,
  `tag_id` bigint DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_post_tag` (`post_id`,`tag_id`),
  KEY `idx_post_tags_tag_id` (`tag_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `post_topics` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `post_id` bigint DEFAULT NULL,
  `topic_id` bigint DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_post_topic` (`post_id`,`topic_id`),
  KEY `idx_post_topics_topic_id` (`topic_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `posts` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `author_id` bigint DEFAULT NULL,
  `title` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `content_md` longtext COLLATE utf8mb4_unicode_ci,
  `summary` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `post_type` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT 'original',
  `source_post_id` bigint DEFAULT NULL,
  `repost_comment` varchar(1000) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `circle_id` bigint DEFAULT NULL,
  `visibility` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT 'public',
  `status` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT 'published',
  `cover_url` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `view_count` bigint DEFAULT '0',
  `like_count` bigint DEFAULT '0',
  `comment_count` bigint DEFAULT '0',
  `favorite_count` bigint DEFAULT '0',
  `share_count` bigint DEFAULT '0',
  `repost_count` bigint DEFAULT '0',
  `hot_score` bigint DEFAULT '0',
  `scheduled_at` datetime(3) DEFAULT NULL,
  `published_at` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_author_status` (`author_id`,`status`),
  KEY `idx_circle_status` (`circle_id`,`status`),
  KEY `idx_posts_hot_score` (`hot_score`),
  KEY `idx_posts_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `tags` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `description` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `status` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT 'enabled',
  `use_count` bigint DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_tag_name` (`name`),
  KEY `idx_tags_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `topics` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `description` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `cover_url` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `is_official` tinyint(1) DEFAULT '0',
  `is_recommended` tinyint(1) DEFAULT '0',
  `participant_count` bigint DEFAULT '0',
  `post_count` bigint DEFAULT '0',
  `status` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT 'enabled',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_topic_name` (`name`),
  KEY `idx_topics_is_recommended` (`is_recommended`),
  KEY `idx_topics_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_follows` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `follower_id` bigint DEFAULT NULL,
  `followee_id` bigint DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_follow` (`follower_id`,`followee_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_point_logs` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint DEFAULT NULL,
  `action` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `point` bigint DEFAULT NULL,
  `biz_type` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `biz_id` bigint DEFAULT NULL,
  `remark` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_action_biz` (`user_id`,`action`,`biz_type`,`biz_id`),
  KEY `idx_user_created` (`user_id`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `account` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `password_hash` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `nickname` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `avatar` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `bio` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `role` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT 'user',
  `status` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT 'normal',
  `follower_count` bigint DEFAULT '0',
  `following_count` bigint DEFAULT '0',
  `post_count` bigint DEFAULT '0',
  `comment_count` bigint DEFAULT '0',
  `like_count` bigint DEFAULT '0',
  `point_count` bigint DEFAULT '0',
  `level` bigint DEFAULT '1',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_account` (`account`),
  KEY `idx_users_status` (`status`),
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `worker_event_records` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `event_id` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `worker_name` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `status` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT 'success',
  `error_message` varchar(1000) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_event_worker` (`event_id`,`worker_name`),
  KEY `idx_worker_created` (`worker_name`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
