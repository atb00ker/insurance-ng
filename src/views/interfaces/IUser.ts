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

export interface IUserMetrics {
  name: string;
  datapoint: {
    age: {
      title: string;
      explaination: string;
      value: string;
      score: number;
    };
    health: {
      title: string;
      explaination: string;
      value: number;
      score: number;
    };
    travel: {
      title: string;
      explaination: string;
      value: number;
      score: number;
    };
    dept: {
      title: string;
      explaination: string;
      value: number;
      score: number;
    };
    wealth: {
      title: string;
      explaination: string;
      value: number;
      score: number;
    };
    investment: {
      title: string;
      explaination: string;
      score: number;
      sip_percent: number;
      mutualfund_percent: number;
      nps_percent: number;
      ppf_percent: number;
      epf_percent: number;
    };
    motor: {
      title: string;
      explaination: string;
      value: number;
      score: number;
    };
  };
}
