import { Component } from '@angular/core';
import {LayoutPartialBreadcrumbComponent} from "@app/component/breadcrumb/breadcrumb.component";
import {PaginationComponent} from "@app/component/pagination/pagination.component";

@Component({
  selector: 'app-users',
  standalone: true,
  imports: [
    LayoutPartialBreadcrumbComponent,
    PaginationComponent
  ],
  templateUrl: './users.component.html',
  styleUrl: './users.component.scss'
})
export class AdminUsersComponent {
  protected userImgDfPath ='../../../assets/img/user-default.png';
  breadcrumbs = [
    {label: 'admin', url: ['/', 'admin']},
    {label: 'users', url: ['/', 'admin', 'users']}
  ]

}
