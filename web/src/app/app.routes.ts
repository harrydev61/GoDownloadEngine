import { Routes } from '@angular/router';
import { LayoutErrorComponent } from './layouts/error/error.component';
import { LayoutDefaultComponent } from '@app/layouts/download/download.component';
import { LayoutAuthComponent } from './layouts/auth/auth.component';

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
      import('./pages/auth/auth.routes').then((m) => m.AuthRoutes),
    data: { preload: false },
  },
  {
    path: '',
    component: LayoutDefaultComponent,
    title: 'Download engine',
    loadChildren: () =>
      import('@app/pages/download/download.routes').then((m) => m.DownloadRoutes),
    data: { preload: false },
  },
  {
    path: 'error',
    component: LayoutErrorComponent,
    title: 'error',
    loadChildren: () =>
      import('./pages/error/error.routes').then((m) => m.ErrorRoutes),
    data: { preload: false },
  },
  /**
   * Page 404
   */
  { path: '**', redirectTo: 'error/404' },
];
