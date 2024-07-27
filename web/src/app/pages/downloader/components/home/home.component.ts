import {Component} from '@angular/core';
import {BaseComponent} from "@app/component/base/base.component";
import {PaginationComponent} from "@app/component/pagination/pagination.component";
import {LayoutPartialBreadcrumbComponent} from "@app/component/breadcrumb/breadcrumb.component";
import {MatSlideToggleModule} from '@angular/material/slide-toggle';

@Component({
  selector: 'downloader-app-home',
  standalone: true,
  imports: [
    BaseComponent,
    PaginationComponent,
    LayoutPartialBreadcrumbComponent,
    MatSlideToggleModule,],
  templateUrl: './home.component.html',
  styleUrl: './home.component.scss'
})
export class DownloaderHomeComponent extends BaseComponent {

}
