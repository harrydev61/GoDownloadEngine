import { Actions, createEffect, ofType } from "@ngrx/effects";
import { loginFailed, loginStart, loginSuccess } from "./actions/login.actions";
import { exhaustMap, map } from "rxjs";
import { AuthenticationService } from "@app/services/authenticate.service";
import { Store } from "@ngrx/store";
import { IAuthState } from "./auth.state";
import { Router } from "@angular/router";
import { Injectable } from "@angular/core";

@Injectable()

export class AuthEffects {
    constructor(
        private actions$: Actions,
        private authenticateService: AuthenticationService,
        private store: Store<IAuthState>,
        private router: Router
    ) {}

    login$ = createEffect(() => {
        return this.actions$.pipe(
            ofType(loginStart),
            exhaustMap((action) => {
                console.log(action)
                return this.authenticateService.login(action.email, action.password).pipe(
                    map((response: any) => {
                        console.log(response)
                        const data = response.data || false;
                        if (data) {
                            return loginSuccess({user: data})
                        } else {
                            return loginFailed({message: response.message || null})
                        }
                    })
                )
            })
        )
    })
}