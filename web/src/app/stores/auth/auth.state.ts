import { User } from "@app/models/user.model";

export interface IAuthState {
    user: User | null;
    isAuthticated: boolean;
    error: string
}

export const initialAuthState: IAuthState = {
    user: null,
    isAuthticated: false,
    error: ''
};