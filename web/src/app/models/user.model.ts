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
    public expireDate: Date,
    public createdDate: Date,
    public verifyOtp: boolean = false,
    public status: number,
  ) {}

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

  public getFullName(firstName: string, lastName: string) {
    return firstName + lastName
  }
}
