import { Component } from '@angular/core';
import {LayoutPartialBreadcrumbComponent} from "@app/component/breadcrumb/breadcrumb.component";

@Component({
  selector: 'app-profile',
  standalone: true,
  imports: [
    LayoutPartialBreadcrumbComponent
  ],
  templateUrl: './profile.component.html',
  styleUrl: './profile.component.scss'
})
export class SettingUserProfileComponent {
  breadcrumbs = [
    {label: 'Setting', url: []},
    {label: 'User', url: []},
    {label: 'Profile', url: ['/','setting','user','profile']}
  ]
  protected userAvatarDefaultPath = 'assets/images/user-default.png'

}
