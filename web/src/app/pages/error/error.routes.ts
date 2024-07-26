import { Routes } from '@angular/router';
import { ErrorNotFoundComponent } from './components/not-found/not-found.component';

export const ErrorRoutes: Routes = [
  { path: '', component: ErrorNotFoundComponent },
  { path: '404', component: ErrorNotFoundComponent },
];
