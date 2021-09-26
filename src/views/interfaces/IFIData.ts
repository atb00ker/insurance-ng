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
  wealth_score: number;
  debt_score: number;
  investment_score: number;
  insurance: IFIInsurance[];
}

export interface IFIInsurance {
  uuid: string;
  type: string;
  title: string;
  description: string;
  account_id: string;
  score: number;
  current_premium: number;
  clauses: Array<string>;
  current_clauses: Array<string>;
  current_cover: number;
  offer_premium: number;
  offer_cover: number;
  is_active: boolean;
  yoy_deduction_rate: number;
  is_insuranceng_account: string;
}
