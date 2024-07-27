import { HttpInterceptorFn } from '@angular/common/http';

export const ErrorInterceptor: HttpInterceptorFn = (req, next) => {
  console.log(req.url);
  return next(req);
};
