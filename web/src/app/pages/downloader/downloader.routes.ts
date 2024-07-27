import {Routes} from '@angular/router';
import {DownloaderHomeComponent} from '@app/pages/downloader/components/home/home.component';
import {DownloaderTaskComponent} from "@app/pages/downloader/components/download-task/download-task.component";

export const DownloaderRoutes: Routes = [
  {path: 'home', component: DownloaderHomeComponent},
  {path: 'download/task', component: DownloaderTaskComponent},
];
