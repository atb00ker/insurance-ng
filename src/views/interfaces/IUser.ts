export interface IUser {
  id?: string;
  name?: string;
  email?: string;
  jwt: () => Promise<string>;
}

export interface IAuth {
  loginWithRedirect: (options?: any) => Promise<void>;
  logout: (options?: any) => void;
  isAuthenticated?: boolean;
  isReady?: boolean;
  user: IUser;
}
