import {Component, OnInit} from '@angular/core';
import {FormControl, FormGroup, ReactiveFormsModule, Validators} from "@angular/forms";
import {AuthenticationService} from "@services/authenticate.service";
import {BaseComponent} from "@app/component/base/base.component";
import {catchError, map} from "rxjs";
import {User} from "@models/user.model";
import { IAuthState } from '@app/stores/auth/auth.state';
import { Store } from '@ngrx/store';
import { loginStart } from '@app/stores/auth/actions/login.actions';

@Component({
  selector: 'pages-auth-login-app-login',
  standalone: true,
  imports: [
    ReactiveFormsModule
  ],
  templateUrl: './login.component.html',
  styleUrl: './login.component.scss'
})
export class PagesAuthLoginComponent extends BaseComponent implements OnInit {
  protected googleIconPath = '../../../../assets/images/social/google.png'
  protected facebookIconPath = '../../../../assets/images/social/facebook.png'
  protected loginForm: any;
  constructor(
    private authenticationService: AuthenticationService,
    private store: Store<IAuthState>,
  ) {
    super();
  }
  ngOnInit(): void {
    this.loginForm = new FormGroup({
      email: new FormControl('', [Validators.required, Validators.email]),
      password: new FormControl('', [Validators.required]),
    });

  }
  checkFormValid() {
    return this.loginForm.invalid;
  }

  onSubmit() {
    const email = this.loginForm.value.email;
    const password = this.loginForm.value.password;
    this.store.dispatch(loginStart({email, password}))
  }
}
