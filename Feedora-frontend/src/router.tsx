import { Navigate, createBrowserRouter } from 'react-router-dom';
import { FrontLayout } from './components/layout/FrontLayout';
import { AdminLayout } from './components/layout/AdminLayout';
import { LoginPage } from './pages/auth/LoginPage';
import { RegisterPage } from './pages/auth/RegisterPage';
import { HomePage } from './pages/home/HomePage';
import { HotPage } from './pages/home/HotPage';
import { PostCreatePage } from './pages/post/PostCreatePage';
import { PostDetailPage } from './pages/post/PostDetailPage';
import { PostEditPage } from './pages/post/PostEditPage';
import { TagDetailPage } from './pages/tag/TagDetailPage';
import { CircleListPage } from './pages/circle/CircleListPage';
import { CircleCreatePage } from './pages/circle/CircleCreatePage';
import { CircleHomePage } from './pages/circle/CircleHomePage';
import { CircleMemberManagePage } from './pages/circle/CircleMemberManagePage';
import { TopicListPage } from './pages/topic/TopicListPage';
import { TopicDetailPage } from './pages/topic/TopicDetailPage';
import { FollowingCenterPage } from './pages/following/FollowingCenterPage';
import { SearchResultPage } from './pages/search/SearchResultPage';
import { NotificationPage } from './pages/notification/NotificationPage';
import { MyProfilePage } from './pages/user/MyProfilePage';
import { UserProfilePage } from './pages/user/UserProfilePage';
import { CreatorCenterPage } from './pages/user/CreatorCenterPage';
import { TaskCenterPage } from './pages/growth/TaskCenterPage';
import { RankingPage } from './pages/growth/RankingPage';
import { ActivityCenterPage } from './pages/activity/ActivityCenterPage';
import { ActivityDetailPage } from './pages/activity/ActivityDetailPage';
import { AnnouncementListPage, AnnouncementDetailPage } from './pages/announcement/AnnouncementPages';
import { WorkspaceLayout } from './pages/workspace/WorkspaceLayout';
import { WorkspaceDashboardPage } from './pages/workspace/WorkspaceDashboardPage';
import { WorkspaceNoteListPage } from './pages/workspace/WorkspaceNoteListPage';
import { WorkspaceNoteEditPage } from './pages/workspace/WorkspaceNoteEditPage';
import { WorkspaceKnowledgeBaseListPage } from './pages/workspace/WorkspaceKnowledgeBaseListPage';
import { WorkspaceKnowledgeBaseDetailPage } from './pages/workspace/WorkspaceKnowledgeBaseDetailPage';
import { WorkspaceAiChatPage } from './pages/workspace/WorkspaceAiChatPage';
import { WorkspaceModelConfigPage } from './pages/workspace/WorkspaceModelConfigPage';
import { AdminActivityPage, AdminAiModelPage, AdminAiTokenConfigPage, AdminAnnouncementPage, AdminBadgePage, AdminBannerPage, AdminCirclePage, AdminCommentPage, AdminDashboardPage, AdminOperationDashboardPage, AdminOperationLogPage, AdminPointLevelPage, AdminPostPage, AdminRecommendPage, AdminReportPage, AdminReviewPage, AdminSearchPage, AdminSensitiveWordPage, AdminTagPage, AdminTaskPage, AdminTopicPage, AdminUserPage } from './pages/admin/AdminPages';

export const router = createBrowserRouter([
  { path:'/login', element:<LoginPage/> },
  { path:'/register', element:<RegisterPage/> },
  { path:'/', element:<FrontLayout/>, children:[
    { index:true, element:<HomePage/> },
    { path:'hot', element:<HotPage/> },
    { path:'following', element:<Navigate to="/following-center" replace/> },
    { path:'following-center', element:<FollowingCenterPage/> },
    { path:'posts/create', element:<PostCreatePage/> },
    { path:'posts/:postId', element:<PostDetailPage/> },
    { path:'posts/:postId/edit', element:<PostEditPage/> },
    { path:'tags/:tagId', element:<TagDetailPage/> },
    { path:'circles', element:<CircleListPage/> },
    { path:'circles/create', element:<CircleCreatePage/> },
    { path:'circles/:circleId', element:<CircleHomePage/> },
    { path:'circles/:circleId/manage/members', element:<CircleMemberManagePage/> },
    { path:'topics', element:<TopicListPage/> },
    { path:'topics/:topicId', element:<TopicDetailPage/> },
    { path:'search', element:<SearchResultPage/> },
    { path:'notifications', element:<NotificationPage/> },
    { path:'users/me', element:<MyProfilePage/> },
    { path:'users/:userId', element:<UserProfilePage/> },
    { path:'tasks', element:<TaskCenterPage/> },
    { path:'rankings', element:<RankingPage/> },
    { path:'creator', element:<CreatorCenterPage/> },
    { path:'activities', element:<ActivityCenterPage/> },
    { path:'activities/:activityId', element:<ActivityDetailPage/> },
    { path:'announcements', element:<AnnouncementListPage/> },
    { path:'announcements/:announcementId', element:<AnnouncementDetailPage/> },
    { path:'workspace', element:<WorkspaceLayout/>, children:[
      { index:true, element:<WorkspaceDashboardPage/> },
      { path:'notes', element:<WorkspaceNoteListPage/> },
      { path:'notes/create', element:<WorkspaceNoteEditPage/> },
      { path:'notes/:noteId', element:<WorkspaceNoteEditPage/> },
      { path:'knowledge-bases', element:<WorkspaceKnowledgeBaseListPage/> },
      { path:'knowledge-bases/:knowledgeBaseId', element:<WorkspaceKnowledgeBaseDetailPage/> },
      { path:'ai-chat', element:<WorkspaceAiChatPage/> },
      { path:'model-config', element:<WorkspaceModelConfigPage/> },
    ]},
  ]},
  { path:'/admin', element:<AdminLayout/>, children:[
    { index:true, element:<AdminDashboardPage/> },
    { path:'users', element:<AdminUserPage/> },
    { path:'posts', element:<AdminPostPage/> },
    { path:'comments', element:<AdminCommentPage/> },
    { path:'reviews', element:<AdminReviewPage/> },
    { path:'reports', element:<AdminReportPage/> },
    { path:'tags', element:<AdminTagPage/> },
    { path:'circles', element:<AdminCirclePage/> },
    { path:'topics', element:<AdminTopicPage/> },
    { path:'search', element:<AdminSearchPage/> },
    { path:'recommend', element:<AdminRecommendPage/> },
    { path:'points', element:<AdminPointLevelPage/> },
    { path:'badges', element:<AdminBadgePage/> },
    { path:'tasks', element:<AdminTaskPage/> },
    { path:'activities', element:<AdminActivityPage/> },
    { path:'banners', element:<AdminBannerPage/> },
    { path:'announcements', element:<AdminAnnouncementPage/> },
    { path:'sensitive-words', element:<AdminSensitiveWordPage/> },
    { path:'logs', element:<AdminOperationLogPage/> },
    { path:'operation-dashboard', element:<AdminOperationDashboardPage/> },
    { path:'ai/models', element:<AdminAiModelPage/> },
    { path:'ai/token-config', element:<AdminAiTokenConfigPage/> },
  ]},
  { path:'*', element:<Navigate to="/" replace/> },
]);
