import { User } from "@app/models/user.model";
import { createAction, props } from "@ngrx/store";

const SIGNUP_START = '[AUTH_PAGE]___Sign up start';
const SIGNUP_SUCCESS = '[AUTH_PAGE]___Sign up sucess';
const SIGNUP_FAILED = '[AUTH_PAGE]__Sign up failed';

export const signupStart = createAction(SIGNUP_START, props<{user: User}>());
export const signupSuccess = createAction(SIGNUP_SUCCESS, props<{user: User}>());
export const signupFailed = createAction(SIGNUP_FAILED, props<{message: string}>());

