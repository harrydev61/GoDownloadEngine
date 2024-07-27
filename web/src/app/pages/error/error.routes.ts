import { Routes } from '@angular/router';
import { ErrorNotFoundComponent } from './components/not-found/not-found.component';
import {ErrorInternalServerComponent} from "./components/internal-server/internal-server.component";
import {ErrorMaintenanceComponent} from "./components/maintenance/maintenance.component";

export const ErrorRoutes: Routes = [
  { path: '', component: ErrorNotFoundComponent },
  { path: '404', component: ErrorNotFoundComponent },
  { path: '500', component: ErrorInternalServerComponent },
  { path: 'maintenance', component: ErrorMaintenanceComponent },

];
