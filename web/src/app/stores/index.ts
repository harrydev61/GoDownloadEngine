
import { isDevMode } from '@angular/core';

import {
    ActionReducerMap,
    MetaReducer
} from '@ngrx/store';

import { IAuthState } from './auth/auth.state';
import { AuthReducers } from './auth/auth.reducers';
import { AuthEffects } from './auth/auth.effects';

export interface RootState {
    auth: IAuthState
}

export const rootReducers: ActionReducerMap<RootState> = {
    auth: AuthReducers
}

export const rootEffects = [AuthEffects]

export const metaReducers: MetaReducer<RootState>[] = isDevMode() ? [] : []