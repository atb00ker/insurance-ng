export enum RouterPath {
  Home = '/',
  Features = '/features/',
  Register = '/register/',
  Dashboard = '/dashboard/',
  About = '/about/',
  InsuranceDetails = '/insurance/:insurance_uuid',
}

export enum ServerPath {
  Register = '/api/v1/register/',
  CreateConsent = '/api/v1/account_aggregator/consent/',
  ConsentNotificationMock = '/api/v1/account_aggregator/Mock/Consent/Notification/',
  GetUserData = '/api/v1/insurance/',
  DataWebsocket = '/api/v1/ws/insurance/',
  InsurancePurchase = '/api/v1/insurance/purchase/',
  InsuranceClaim = '/api/v1/insurance/claim/',
}
