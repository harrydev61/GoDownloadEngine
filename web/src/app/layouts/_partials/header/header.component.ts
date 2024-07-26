import { Component } from '@angular/core';
import { LayoutPartialHeaderSearchComponent } from './search/search.component';

@Component({
  selector: 'layout-partial-app-header',
  standalone: true,
  imports: [LayoutPartialHeaderSearchComponent],
  templateUrl: './header.component.html',
  styleUrl: './header.component.scss'
})
export class LayoutPartialHeaderComponent {

  protected userImgDfPath ='../../../../assets/images/user-default.png'
  protected logoPath = '../../../../assets/logo/logo.png'

}
