import { IUser, IUserMetrics } from './IUser';

export interface IFIData {
  FipId: string;
  RahasyaData: RahasyaList[];
}

export interface RahasyaList {
  data: any;
  errorInfo: string;
}

export interface IInsuranceFIData {
  data: any;
  userInfo: IUserMetrics;
}
