import { Routes } from '@angular/router';
import { PagesAuthLoginComponent } from './login/login.component';


export const AuthRoutes: Routes = [
  { path: '**', redirectTo: 'login' },
  { path: 'login', component: PagesAuthLoginComponent },
];
