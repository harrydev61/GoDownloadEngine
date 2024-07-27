import { Routes } from '@angular/router';
import {AdminDashboardComponent} from "./components/dashboard/dashboard.component";
import {AdminUsersComponent} from "@app/pages/admin/components/users/users.component";
import {AdminRolesComponent} from "@app/pages/admin/components/roles/roles.component";

export const AdminRoutes: Routes = [
  {path:'', redirectTo:'dashboard', pathMatch: 'full'},
  {path:'dashboard', component: AdminDashboardComponent},
  {path:'users', component: AdminUsersComponent},
  {path:'roles', component: AdminRolesComponent}
];
