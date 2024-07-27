import { Component } from '@angular/core';
import {LayoutPartialBreadcrumbComponent} from "@app/component/breadcrumb/breadcrumb.component";
import {PaginationComponent} from "@app/component/pagination/pagination.component";

@Component({
  selector: 'downloader-app-downloader-task',
  standalone: true,
    imports: [
        LayoutPartialBreadcrumbComponent,
        PaginationComponent
    ],
  templateUrl: './download-task.component.html',
  styleUrl: './download-task.component.scss'
})
export class DownloaderTaskComponent {
  breadcrumbs = [
    {label: 'Download task', url: ['/', 'download','task']}
  ]
}
