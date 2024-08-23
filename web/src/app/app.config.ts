import { ApplicationConfig } from '@angular/core';
import { provideRouter } from '@angular/router';

import { routes } from './app.routes';
import { provideAnimationsAsync } from '@angular/platform-browser/animations/async';
import {provideAnimations} from "@angular/platform-browser/animations";
import {provideHttpClient, withFetch, withInterceptors} from "@angular/common/http";
import {ErrorInterceptor} from "@app/interceptors/error.interceptor";
import { provideStore } from '@ngrx/store';
import {rootReducers, rootEffects, metaReducers} from "./stores/index"
import { provideEffects } from '@ngrx/effects';
export const appConfig: ApplicationConfig = {
  providers: [
    provideRouter(routes),
    provideAnimationsAsync(),
    provideAnimations(),
    provideHttpClient(
      withFetch(),
      withInterceptors([ErrorInterceptor])),
    provideStore(rootReducers, {metaReducers}),
    provideEffects(rootEffects)
    ]
};
