export interface AuthResponseData {
  avatar: string;
  id: number;
  email: string;
  username: string;
  fullName: string;
  roleId: number;
  roleName: string;
  accessToken: string;
  refreshToken: string;
  createdTime: string;
  expiresIn: string;
  verifyOtp: boolean;
  apps: string,
  status: number,
  groupUserPermission: any[],
}
