import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { BaseComponent } from '@app/component/base/base.component';
import { AuthenticationService } from '@app/services/authenticate.service';

@Component({
  selector: 'app-page-auth-signup',
  standalone: true,
  imports: [ReactiveFormsModule],
  templateUrl: './signup.component.html',
  styleUrl: './signup.component.scss'
})
export class PageAuthSignUpComponent extends BaseComponent implements OnInit {
	protected signupForm!: FormGroup;
	constructor(
		private authenticationService: AuthenticationService
	) {
		super();
	}
	ngOnInit(): void {
		this.signupForm = new FormGroup({
			email: new FormControl('', [Validators.required, Validators.email]),
			password: new FormControl('', [Validators.required, Validators.minLength(5), Validators.maxLength(50)]),
			firstName: new FormControl('', [Validators.required, Validators.minLength(2), Validators.maxLength(50)]),
			lastName: new FormControl('', [Validators.required, Validators.minLength(2), Validators.maxLength(50)]),
			phone: new FormControl('', [Validators.required]),
			dob: new FormControl('', [Validators.required]),
			gender: new FormControl('', [Validators.required])
		})
	}

	onSubmit() {
		const data = this.signupForm.value;
		if (this.signupForm.valid) {
			this.authenticationService.signup(data)
		}
	}
}
