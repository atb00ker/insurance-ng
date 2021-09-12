export enum ConsentStatus {
	UserConsentsPending  = "PENDING",
	UserConsentsReady    = "READY",
	UserConsentsActive   = "ACTIVE",
	UserConsentsFetched   = "FETCHED",
	UserConsentsRejected = "REJECTED",
	UserConsentsRevoked  = "REVOKED",
	UserConsentsNorecord = "NOTFOUND",
	UserConsentsUnknown  = "UNKNOWN"
}
