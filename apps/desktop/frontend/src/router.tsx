import { Navigate, createHashRouter } from 'react-router-dom';
import { DesktopShell } from './shell/DesktopShell';
import { CommunityPage } from './pages/CommunityPage';
import { WorkspaceDashboardPage } from './pages/WorkspaceDashboardPage';
import { NotesPage } from './pages/NotesPage';
import { KnowledgeBasesPage } from './pages/KnowledgeBasesPage';
import { RepositoryPage } from './pages/RepositoryPage';
import { SearchPage } from './pages/SearchPage';
import { NotificationsPage } from './pages/NotificationsPage';
import { SettingsPage } from './pages/SettingsPage';

export const router = createHashRouter([
  {
    path:'/',
    element:<DesktopShell />,
    children:[
      {index:true,element:<Navigate to="/community" replace />},
      {path:'community',element:<CommunityPage />},
      {path:'workspace',element:<WorkspaceDashboardPage />},
      {path:'workspace/notes',element:<NotesPage />},
      {path:'workspace/knowledge-bases',element:<KnowledgeBasesPage />},
      {path:'repository',element:<RepositoryPage />},
      {path:'search',element:<SearchPage />},
      {path:'notifications',element:<NotificationsPage />},
      {path:'settings',element:<SettingsPage />},
    ],
  },
]);
