import {Injectable} from '@angular/core';
import {LocalStorageService} from "@services/local-storage.service";
import {HttpService} from "@services/http.service";
import {environment} from "../../environments/environment";
import {AuthResponseData} from "@models/auth-response-data.model";
import {User} from "@models/user.model";
import {BehaviorSubject, Observable} from "rxjs";

@Injectable({
  providedIn: 'root'
})
export class AuthenticationService {
  public key = 'current_user';
  private originalUserKey = 'user_original';
  private apiServerPath = environment.backendServer.paths;
  private currentUserSubject: BehaviorSubject<any>;
  public currentUser: Observable<any>;
  timeoutInterval: any;
  private originalUserSubject: BehaviorSubject<any>;
  public originalUser: Observable<User>;

  constructor(private localStorageService: LocalStorageService,
              private httpService: HttpService) {
    // @ts-ignore
    this.currentUserSubject = new BehaviorSubject<any>(this.localStorageService.getCurrentUser());
    this.currentUser = this.currentUserSubject.asObservable();
    this.originalUserSubject = new BehaviorSubject<any>(this.getUserOriginalFromLocalStorage() || null);
    this.originalUser = this.originalUserSubject.asObservable();
  }

  login(email: string, password: string) {
    return this.httpService.post(this.apiServerPath.auth.login, {
      email: email,
      password: password
    })
  }

  logout() {
    this.localStorageService.delete(this.key);
    this.currentUserSubject.next(null);
    if (this.timeoutInterval) {
      clearTimeout(this.timeoutInterval);
      this.timeoutInterval = null;
    }
  }

  formatUser(data: AuthResponseData) {
    const createdDate = new Date(parseInt(data.createdTime) * 1000);
    const expirationDate = new Date(
      createdDate.getTime() + parseInt(data.expiresIn) * 1000
    );
    const user = new User(
      data.id,
      data.email,
      data.username,
      data.fullName,
      data.avatar,
      data.roleId,
      data.roleName,
      data.accessToken,
      data.refreshToken,
      createdDate,
      expirationDate,
      data.verifyOtp || false,
      data.apps,
      data.status,
      data.groupUserPermission
    );
    return user;
  }

  getErrorMessage(message: string) {
    switch (message) {
      case 'EMAIL_NOT_FOUND':
        return 'Email Not Found';
      case 'INVALID_PASSWORD':
        return 'Invalid Password';
      case 'EMAIL_EXISTS':
        return 'Email already exists';
      default:
        return 'Unknown error occurred. Please try again';
    }
  }

  setUserInLocalStorage(user: User) {
    this.localStorageService.setUser(user);
    this.runTimeoutInterval(user);
  }

  setOriginalUserInLocalStorage(originalUser: any) {
    this.localStorageService.setItem(originalUser, this.originalUserKey)
  }

  setUserSwitchInLocalStorage(user: User) {
    this.localStorageService.setUser(user);
  }

  runTimeoutInterval(user: User) {
    const todaysDate = new Date().getTime();
    const expirationDate = user.expireDate.getTime();
    const timeInterval = expirationDate - todaysDate;

    this.timeoutInterval = setTimeout(() => {
      //logout functionality or get the refresh token
    }, timeInterval);
  }

  getUserFromLocalStorage() {
    const userData = this.localStorageService.getCurrentUser()
    if (userData) {
      const createdDate = new Date(userData.createdDate);
      const expirationDate = new Date(userData.expirationDate);
      const user = new User(
        userData.userId,
        userData.email,
        userData.username,
        userData.fullName,
        userData.avatar,
        userData.roleId,
        userData.roleName,
        userData.accessToken,
        userData.refreshToken,
        createdDate,
        expirationDate,
        userData.verifyOtp,
        userData.apps,
        userData.status,
        userData.groupUserPermission,
      );
      // this.runTimeoutInterval(user);
      return user;
    }
    return null;
  }

  getUserOriginalFromLocalStorage() {
    const userData: any = this.localStorageService.get(this.originalUserKey);
    return JSON.parse(userData || null);
  }

  getAppSelectedFromLocalStorage() {
    const appData: any = this.localStorageService.get('domain_selected');
    return JSON.parse(appData || null);
  }

  removeUserOriginal() {
    this.localStorageService.remove(this.originalUserKey);
  }

}

