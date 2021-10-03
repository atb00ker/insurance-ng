export enum RouterPath {
  Home = '/',
  Features = '/features/',
  Register = '/register/',
  Dashboard = '/dashboard/',
  About = '/about/',
  InsuranceDetails = '/insurance/:insurance_uuid',
}

export enum ServerPath {
  Register = '/api/register/',
  CreateConsent = '/api/account_aggregator/consent/',
  ConsentNotificationMock = '/api/account_aggregator/Mock/Consent/Notification/',
  GetUserData = '/api/insurance/',
  DataWebsocket = '/api/ws/insurance/',
  InsurancePurchase = '/api/insurance/purchase/',
}
