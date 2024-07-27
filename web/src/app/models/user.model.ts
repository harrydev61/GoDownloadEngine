import {environment} from "../../environments/environment";

export class User {
  constructor(
    public userId: number,
    public email: string,
    public username: string,
    public fullName: string,
    public avatar: string,
    public roleId: number,
    public roleName: string,
    public accessToken: string,
    public refreshToken: string,
    public createdDate: Date,
    public expirationDate: Date,
    public verifyOtp: boolean = false,
    public apps: string,
    public status: number,
    public groupUserPermission: any[] = [],
  ) {}

  get expireDate() {
    return this.expirationDate;
  }

  get userAccessToken() {
    return this.accessToken;
  }

  get userRefreshToken() {
    return this.refreshToken;
  }

  get userUserName() {
    return this.username;
  }

  public isSuperAdmin() {
    return this.roleId === environment.roles.superAdmin;
  }

  public isAdmin() {
    return this.roleId === environment.roles.superAdmin || this.roleId === environment.roles.admin;
  }
}
