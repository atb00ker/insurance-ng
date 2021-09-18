export enum ConsentStatus {
  UserConsentsPending = 'PENDING',
  UserConsentsReady = 'READY',
  UserConsentsActive = 'ACTIVE',
  UserConsentsFetched = 'FETCHED',
  UserConsentsRejected = 'REJECTED',
  UserConsentsRevoked = 'REVOKED',
  UserConsentsNorecord = 'NOTFOUND',
  UserConsentsUnknown = 'UNKNOWN',
}

export enum ConsentType {
  Transactions = 'Transactions',
  Summary = 'Summary',
  Profile = 'Profile',
}
