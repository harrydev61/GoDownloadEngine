import { Routes } from '@angular/router';
import { PagesAuthLoginComponent } from '@app/pages/auth/components/login/login.component';


export const AuthRoutes: Routes = [
  { path: '**', redirectTo: 'login' },
  { path: 'login', component: PagesAuthLoginComponent },
];
