export type IUser = {
  id?: string;
  name?: string;
  email?: string;
  jwt: () => Promise<string>;
};

export type IAuth = {
  loginWithRedirect: (options?: any) => Promise<void>;
  logout: (options?: any) => void;
  isAuthenticated?: boolean;
  isReady?: boolean;
  user: IUser;
};

export type IUserProfileScores = {
  title: string;
  explaination: string;
  score: number;
};
