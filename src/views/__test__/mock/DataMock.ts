import { InsuranceTypes } from '../../enums/Insurance';
import { IFIData, IFIInsurance } from '../../types/IFIData';
import { IUser } from '../../types/IUser';

export const userJohnDoe: IUser = {
  id: '12345',
  name: 'John Doe',
  email: 'john@example.com',
  jwt: () => new Promise(() => '12345'),
};

export const insuranceMedicalJohnDoe: IFIInsurance = {
  uuid: '94be9707-d419-404d-a3e2-5be54019fe85',
  type: InsuranceTypes.MEDICAL_PLAN,
  title: 'Medical Plan',
  description: 'Medical Plan Description',
  account_id: 'TESTACCT01',
  score: 0.5,
  clauses: ['Clause 1', 'Clause 2'],
  current_clauses: ['Clause 2', 'Clause 3'],
  current_premium: 1,
  current_cover: 10,
  offer_premium: 1.5,
  offer_cover: 11,
  is_active: true,
  is_claimed: false,
  yoy_deduction_rate: 0.5,
  is_insurance_ng_acct: false
}

export const fiDataJohnDoe: IFIData = {
  status: false,
  data: {
    name: 'John Doe',
    date_of_birth: new Date(2011,10,30),
    pancard: 'TESTPAN01',
    ckyc_compliance: true,
    age_score: 0,
    wealth_score: 0,
    debt_score: 0,
    investment_score: 0,
    phone: '999999999999999',
    shared_data_sources: 0,
    insurance: [insuranceMedicalJohnDoe]
  }
}
