export enum RouterPath {
  CreateConsent = '/',
  Register = '/register/',
  Dashboard = '/dashboard/',
  InsuranceDetails = '/insurance/:insurance_uuid',
}

export enum ServerPath {
  Register = '/api/register/',
  CreateConsent = '/api/account_aggregator/consent/',
  ConsentStatus = '/api/account_aggregator/consent/status/',
  GetUserData = '/api/insurance/',
  InsurancePurchase = '/api/insurance/purchase/',
}
