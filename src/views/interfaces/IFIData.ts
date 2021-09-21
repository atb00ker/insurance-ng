export interface IFIData {
  status: boolean;
  data: IFIUserData;
  error?: string;
}

export interface IFIUserData {
  name: string;
  date_of_birth: Date;
  pancard: string;
  ckyc_compliance: boolean;
  age_score: number;
  insurance: IFIInsurance[];
}

export interface IFIInsurance {
  title: string;
  description: string;
  account_id: string;
  score: number;
  current_premium: number;
  current_cover: number;
  offer_premium: number;
  offer_cover: number;
  type: string;
}
