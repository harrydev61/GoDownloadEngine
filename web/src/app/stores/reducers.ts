
import { isDevMode } from '@angular/core';

import {
    ActionReducer,
    ActionReducerMap,
    createFeatureSelector,
    createSelector,
    MetaReducer
} from '@ngrx/store';
import { loginReducer } from './auth/reducers/login.reducers';

export interface State {

}

export const reducers: ActionReducerMap<State> = {
    login: loginReducer
}

export const metaReducers: MetaReducer<State>[] = isDevMode() ? [] : []