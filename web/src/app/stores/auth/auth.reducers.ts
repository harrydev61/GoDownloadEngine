import { createReducer, on } from "@ngrx/store";
import {loginSuccess} from "./actions/login.actions";
import { initialAuthState } from "./auth.state";
import { signupSuccess } from "./actions/signup.actions";
export const initialState = initialAuthState;

export const AuthReducers = createReducer(initialState, 
    on(loginSuccess, (state, action) => {return {...state, user: action.user}}),
    on(signupSuccess, (state, action) => {return {...state, user: action.user}})
)