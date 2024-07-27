import { Component } from '@angular/core';
import {LayoutPartialBreadcrumbComponent} from "@app/component/breadcrumb/breadcrumb.component";
import {PaginationComponent} from "@app/component/pagination/pagination.component";

@Component({
  selector: 'app-roles',
  standalone: true,
  imports: [
    LayoutPartialBreadcrumbComponent,
    PaginationComponent
  ],
  templateUrl: './roles.component.html',
  styleUrl: './roles.component.scss'
})
export class AdminRolesComponent {
  breadcrumbs = [
    {label: 'admin', url: ['/', 'admin']},
    {label: 'roles', url: ['/', 'admin', 'roles']}
  ]
}
