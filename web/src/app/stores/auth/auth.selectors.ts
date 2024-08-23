import { createFeatureSelector, createSelector } from "@ngrx/store";
import { IAuthState } from "./auth.state";

export const selectAuthState = createFeatureSelector<IAuthState>('auth');

export const selectLoginState = createSelector(
    selectAuthState,
    (state: IAuthState) => state.user
)