import {Component} from '@angular/core';
import {RouterOutlet} from '@angular/router';
import {LayoutPartialSidebarComponent} from '../_partials/sidebar/sidebar.component';
import {LayoutPartialHeaderComponent} from '../_partials/header/header.component';
import {LayoutPartialFooterComponent} from '../_partials/footer/footer.component';
import {ModeService} from "@services/mode.service";
import {NgClass} from "@angular/common";

@Component({
  selector: 'app-downloader',
  standalone: true,
  imports: [
    RouterOutlet,
    LayoutPartialSidebarComponent,
    LayoutPartialHeaderComponent,
    LayoutPartialFooterComponent,
    NgClass
  ],
  templateUrl: './downloader.component.html',
  styleUrl: './downloader.component.scss'
})
export class LayoutDownloaderComponent {
  constructor() {

  }
}
