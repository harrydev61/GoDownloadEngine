import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { LayoutPartialSidebarComponent } from '../_partials/sidebar/sidebar.component';
import { LayoutPartialHeaderComponent } from '../_partials/header/header.component';
import { LayoutPartialFooterComponent } from '../_partials/footer/footer.component';
import {LayoutPartialBreadcrumbComponent} from "@app/layouts/_partials/breadcrumb/breadcrumb.component";

@Component({
  selector: 'app-download',
  standalone: true,
  imports: [
    RouterOutlet,
    LayoutPartialSidebarComponent,
    LayoutPartialHeaderComponent,
    LayoutPartialFooterComponent,
    LayoutPartialBreadcrumbComponent],
  templateUrl: './download.component.html',
  styleUrl: './download.component.scss'
})
export class LayoutDefaultComponent {

}
