import type { CreatePostRequest, NotificationItem, PageResult, Post, PostQuery } from '@feedora/contracts';

export interface CommunitySearchInput {
  keyword: string;
  type?: string;
  page?: number;
  pageSize?: number;
}

export interface CommunityPort {
  getPosts(query: PostQuery): Promise<PageResult<Post>>;
  getPost(postId: number): Promise<Post>;
  createPost(input: CreatePostRequest): Promise<Post>;
  getNotifications(): Promise<NotificationItem[]>;
  search(input: CommunitySearchInput): Promise<unknown>;
}
