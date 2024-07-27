import {Routes} from '@angular/router';
import {LayoutErrorComponent} from './layouts/error/error.component';
import {LayoutDownloaderComponent} from '@app/layouts/downloader/downloader.component';
import {LayoutAuthComponent} from './layouts/auth/auth.component';

export const routes: Routes = [
  {
    path: '',
    redirectTo: '/home',
    pathMatch: 'full',
  },
  {
    path: 'auth',
    component: LayoutAuthComponent,
    title: 'auth',
    loadChildren: () =>
      import('@app/pages/auth/auth.routes').then((m) => m.AuthRoutes),
    data: {preload: false},
  },
  {
    path: '',
    component: LayoutDownloaderComponent,
    title: 'Download engine',
    loadChildren: () =>
      import('@app/pages/downloader/downloader.routes').then((m) => m.DownloaderRoutes),
    data: {preload: false},
  },
  {
    path: 'setting',
    component: LayoutDownloaderComponent,
    title: 'setting',
    loadChildren: () =>
      import('@app/pages/setting/setting.routes').then((m) => m.SettingRoutes),
    data: {preload: false},
  },
  {
    path: 'admin',
    component: LayoutDownloaderComponent,
    title: 'admin',
    loadChildren: () =>
      import('@app/pages/admin/admin.routes').then((m) => m.AdminRoutes),
    data: {preload: false},
  },
  {
    path: 'error',
    component: LayoutErrorComponent,
    title: 'error',
    loadChildren: () =>
      import('@app/pages/error/error.routes').then((m) => m.ErrorRoutes),
    data: {preload: false},
  },
  /**
   * Page 404
   */
  {path: '**', redirectTo: 'error/404'},
];
