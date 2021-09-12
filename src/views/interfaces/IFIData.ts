import { IUser } from './IUser';

export interface IFIData {
  FipId: string,
  RahasyaData: RahasyaList[]
}

export interface RahasyaList {
 data: string,
 errorInfo: string,
}
