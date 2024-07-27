import { Component } from '@angular/core';
import {LayoutPartialBreadcrumbComponent} from "@app/component/breadcrumb/breadcrumb.component";

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [
    LayoutPartialBreadcrumbComponent
  ],
  templateUrl: './dashboard.component.html',
  styleUrl: './dashboard.component.scss'
})
export class AdminDashboardComponent {
  breadcrumbs = [
    {label: 'Admin', url: ['/', 'admin']},
    {label: 'Dashboard', url: ['/', 'dashboard']}
  ]
}
