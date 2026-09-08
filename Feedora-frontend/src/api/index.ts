import { httpDelete, httpGet, httpPost, httpPut } from './request';

// 统一真实接口调用层：方法名 → REST 端点映射。
// 所有数据均来自后端真实接口，不再使用任何本地 mock 假数据。
// transform：后端返回结构与前端页面期望不一致时的适配（如评论树数组 → PageResult）。
function realEndpoint(methodName: string, args: any[]): { method: 'get' | 'post' | 'put' | 'delete'; url: string; params?: any; data?: any; transform?: (res: any) => any } {
  const [a, b, c] = args;
  const maps: Record<string, () => { method: 'get' | 'post' | 'put' | 'delete'; url: string; params?: any; data?: any; transform?: (res: any) => any }> = {
    login: () => ({ method: 'post', url: '/auth/login', data: { account: a, password: b } }),
    register: () => ({ method: 'post', url: '/auth/register', data: a }),
    getUsers: () => ({ method: 'get', url: '/users', params: a }),
    getUser: () => ({ method: 'get', url: `/users/${a}` }),
    updateUser: () => ({ method: 'put', url: '/users/me/profile', data: a }),
    getTags: () => ({ method: 'get', url: '/tags' }),
    getTag: () => ({ method: 'get', url: `/tags/${a}` }),
    getTopics: () => ({ method: 'get', url: '/topics', params: a }),
    getTopic: () => ({ method: 'get', url: `/topics/${a}` }),
    getCircles: () => ({ method: 'get', url: '/circles', params: a }),
    getCircle: () => ({ method: 'get', url: `/circles/${a}` }),
    createCircle: () => ({ method: 'post', url: '/circles', data: a }),
    joinCircle: () => ({ method: 'post', url: `/circles/${a}/join`, data: { reason: b } }),
    leaveCircle: () => ({ method: 'post', url: `/circles/${a}/leave` }),
    getPosts: () => ({ method: 'get', url: '/posts', params: a }),
    getPost: () => ({ method: 'get', url: `/posts/${a}` }),
    createPost: () => ({ method: 'post', url: '/posts', data: a }),
    updatePost: () => ({ method: 'put', url: `/posts/${a}`, data: b }),
    hidePost: () => ({ method: 'put', url: `/posts/${a}/hide` }),
    unhidePost: () => ({ method: 'put', url: `/posts/${a}/unhide` }),
    deletePost: () => ({ method: 'delete', url: `/posts/${a}` }),
    likePost: () => ({ method: 'post', url: `/posts/${a}/like` }),
    favoritePost: () => ({ method: 'post', url: `/posts/${a}/favorite` }),
    sharePost: () => ({ method: 'post', url: `/posts/${a}/share`, data: { channel: 'copy_link' } }),
    repostPost: () => ({ method: 'post', url: `/posts/${a}/repost`, data: { repostComment: b } }),
    getComments: () => ({
      method: 'get', url: `/posts/${a}/comments`,
      // 后端返回根评论树数组，页面按 PageResult 消费，这里包装对齐。
      transform: (r: any) => (Array.isArray(r) ? { list: r, total: r.length, page: 1, pageSize: r.length } : r),
    }),
    addComment: () => ({ method: 'post', url: '/comments', data: { postId: a, content: b } }),
    getMyComments: () => ({ method: 'get', url: '/users/me/comments', params: a }),
    deleteComment: () => ({ method: 'delete', url: `/comments/${a}` }),
    getCircleMembers: () => ({ method: 'get', url: `/circles/${a}/members`, params: b }),
    setMemberRole: () => ({ method: 'put', url: `/circles/${a}/members/${b}/role`, data: { role: c } }),
    muteMember: () => ({ method: 'put', url: `/circles/${a}/members/${b}/mute`, data: { duration: c, reason: args[3] } }),
    unmuteMember: () => ({ method: 'put', url: `/circles/${a}/members/${b}/unmute` }),
    removeMember: () => ({ method: 'delete', url: `/circles/${a}/members/${b}` }),
    getHotRanks: () => ({ method: 'get', url: '/hot/ranks', params: { rankType: a, timeRange: b } }),
    getNotificationList: () => ({ method: 'get', url: '/notifications' }),
    markNotificationRead: () => ({ method: 'put', url: `/notifications/${a}/read` }),
    getUnreadCount: () => ({ method: 'get', url: '/notifications/unread-count' }),
    search: () => ({ method: 'get', url: '/search', params: { keyword: a, type: b, ...(c || {}) } }),
    searchSuggest: () => ({ method: 'get', url: '/search/suggest', params: { keyword: a } }),
    getHotKeywords: () => ({ method: 'get', url: '/search/hot-keywords' }),
    checkIn: () => ({ method: 'post', url: '/growth/check-in' }),
    getTasks: () => ({ method: 'get', url: '/growth/tasks', params: { type: a } }),
    claimTask: () => ({ method: 'post', url: `/growth/tasks/${a}/claim` }),
    getRankings: () => ({ method: 'get', url: '/growth/rankings', params: { type: a, range: b } }),
    getMyPosts: () => ({ method: 'get', url: '/users/me/posts', params: a }),
    getMyLikedPosts: () => ({ method: 'get', url: '/users/me/liked-posts', params: a }),
    getMyFavoritePosts: () => ({ method: 'get', url: '/users/me/favorite-posts', params: a }),
    // 关注体系
    followUser: () => ({ method: 'post', url: `/users/${a}/follow` }),
    unfollowUser: () => ({ method: 'delete', url: `/users/${a}/follow` }),
    getUserFollowState: () => ({ method: 'get', url: `/users/${a}/follow/state` }),
    getFollowingUsers: () => ({ method: 'get', url: `/users/${a?.userId ?? 'me'}/following`, params: { page: a?.page, pageSize: a?.pageSize } }),
    getFollowers: () => ({ method: 'get', url: `/users/${a}/followers`, params: b }),
    getFollowingFeed: () => ({ method: 'get', url: '/users/me/following-feed', params: { page: a?.page, pageSize: a?.pageSize, feedTab: a?.feedTab ?? 'all' } }),
    getFollowingCircles: () => ({ method: 'get', url: '/users/me/following-circles', params: { page: a?.page, pageSize: a?.pageSize } }),
    getFollowingTopics: () => ({ method: 'get', url: '/users/me/following-topics', params: { page: a?.page, pageSize: a?.pageSize } }),
    getFollowingTags: () => ({ method: 'get', url: '/users/me/following-tags', params: { page: a?.page, pageSize: a?.pageSize } }),
    // 后台管理
    adminGetUsers: () => ({ method: 'get', url: '/admin/users', params: a }),
    adminGetPosts: () => ({ method: 'get', url: '/admin/posts', params: a }),
    adminGetComments: () => ({ method: 'get', url: '/admin/comments', params: a }),
    adminGetTags: () => ({ method: 'get', url: '/admin/tags' }),
    adminCreateTag: () => ({ method: 'post', url: '/admin/tags', data: a }),
    adminUpdateTag: () => ({ method: 'put', url: `/admin/tags/${a}`, data: b }),
    adminGetTopics: () => ({ method: 'get', url: '/admin/topics', params: a }),
    adminCreateTopic: () => ({ method: 'post', url: '/admin/topics', data: a }),
    adminUpdateTopic: () => ({ method: 'put', url: `/admin/topics/${a}`, data: b }),
    adminGetCircles: () => ({ method: 'get', url: '/admin/circles', params: a }),
    adminStats: () => ({ method: 'get', url: '/admin/dashboard/stats' }),
    getAdminLogs: () => ({ method: 'get', url: '/admin/operation-logs', params: a }),
  };
  return maps[methodName]?.() || { method: 'get', url: `/__unimplemented__/${methodName}`, params: {} };
}

type ApiMethod = (...args: any[]) => Promise<any>;

// api 为真实接口代理：任意方法名按 realEndpoint 映射发起真实 HTTP 请求。
// 未映射到后端接口的方法（关注、AI 工作台、活动、公告等占位能力）本地静默失败：
// 不发请求、不弹错误提示，页面保持空态，避免每次加载出现「资源不存在」弹窗。
export const api = new Proxy({} as Record<string, ApiMethod>, {
  get(_target, prop) {
    const methodName = String(prop);
    return async (...args: any[]) => {
      const endpoint = realEndpoint(methodName, args);
      if (endpoint.url.startsWith('/__unimplemented__')) {
        console.warn(`[api] 方法未映射到后端接口: ${methodName}`);
        return Promise.reject(new Error(`api.${methodName} 未实现`));
      }
      if (endpoint.method === 'get') {
        const data = await httpGet(endpoint.url, endpoint.params);
        return endpoint.transform ? endpoint.transform(data) : data;
      }
      if (endpoint.method === 'post') return httpPost(endpoint.url, endpoint.data);
      if (endpoint.method === 'put') return httpPut(endpoint.url, endpoint.data);
      return httpDelete(endpoint.url, endpoint.params);
    };
  },
});

// 兼容旧的模块导出（现均指向真实接口）。
export const realApi = api;
export const mockApi = api;
