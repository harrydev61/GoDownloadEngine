import { User } from "@app/models/user.model";
import { createAction, props } from "@ngrx/store";
const LOGIN_START = '[AUTH_PAGE]___Login start';
const LOGIN_SUCCESS = '[AUTH_PAGE]___Login successful';
const LOGIN_FAILED = '[AUTH_PAGE]___Login failed';

export const loginStart =  createAction(LOGIN_START, props<{email: string; password: string}>());
export const loginSuccess =  createAction(LOGIN_SUCCESS, props<{user: User}>());
export const loginFailed =  createAction(LOGIN_FAILED, props<{message: string}>());
