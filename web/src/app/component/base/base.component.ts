import { Component } from '@angular/core';
import {PaginationComponent} from "@app/component/pagination/pagination.component";
import { AppConstant } from '@app/constants/app.contants';
@Component({
  selector: 'app-base',
  standalone: true,
  imports: [],
  template: `

  `
})
export class BaseComponent {
	protected constant = AppConstant;
}
