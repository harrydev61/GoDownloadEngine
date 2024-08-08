import { createReducer, on } from "@ngrx/store";
import {loginStart, loginExcuting, loginEnd} from "../actions/login.actions";
import { environment } from "src/environments/environment";

export const initialState = environment.state.login.loginConfirm;

export const loginReducer = createReducer(initialState, 
    on(loginStart, (state) => environment.state.login.loginStart),
    on(loginExcuting, (state) => environment.state.login.loginExcute),
    on(loginEnd, (state) => environment.state.login.loginEnd),
)