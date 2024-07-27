import { Component } from '@angular/core';
import { LayoutPartialHeaderSearchComponent } from './search/search.component';
import {ModeService, ModeType} from "@services/mode.service";

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

  constructor(protected modeService: ModeService) {
  }

  protected toggleMode(mode:ModeType): void {
    this.modeService.updateMode(mode)
  }

  protected getDarkMode(): ModeType {
    return ModeType.DARK
  }

  protected getLightMode(): ModeType {
    return ModeType.LIGHT
  }

  protected readonly ModeType = ModeType;
}
