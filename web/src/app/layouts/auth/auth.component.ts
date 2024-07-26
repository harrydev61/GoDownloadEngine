import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { LayoutPartialFooterComponent } from '../_partials/footer/footer.component';

@Component({
  selector: 'layout-app-auth',
  standalone: true,
  imports: [RouterOutlet,LayoutPartialFooterComponent],
  templateUrl: './auth.component.html',
  styleUrl: './auth.component.scss'
})
export class LayoutAuthComponent {

}
