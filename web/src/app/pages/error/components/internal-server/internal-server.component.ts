import { Component } from '@angular/core';

@Component({
  selector: 'app-internal-server',
  standalone: true,
  imports: [],
  templateUrl: './internal-server.component.html',
  styleUrl: './internal-server.component.scss'
})
export class ErrorInternalServerComponent {
  protected internalServerImagePath = "../../../assets/images/illustrations/500.svg";

}
