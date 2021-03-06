export type IFIData = {
  status: boolean;
  data: IFIUserData;
  error?: string;
};

export type IFIUserData = {
  name: string;
  date_of_birth: Date;
  pancard: string;
  ckyc_compliance: boolean;
  age_score: number;
  wealth_score: number;
  debt_score: number;
  investment_score: number;
  phone: string;
  shared_data_sources: number;
  insurance: IFIInsurance[];
};

export type IFIInsurance = {
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
  is_claimed: boolean;
  yoy_deduction_rate: number;
  is_insurance_ng_acct: boolean;
};
