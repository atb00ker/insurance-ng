package account_aggregator

import (
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/google/uuid"
)

const (
	SetuApiCreateConsentPath = "Consent"
	SetuApiConsentStatusPath = "/Consent/handle/%s"
	SetuApiConsentSignedPath = "/Consent/%s"
	SetuApiFiRequest         = "/FI/request"
	SetuApiFiDataFetch       = "/FI/fetch/%s"
	RahasyaApiGetKeys        = "/ecc/v1/generateKey"
	RahasyaApiDecrypt        = "/ecc/v1/decrypt"
)

const (
	// Consent Mode
	ConsentModeView   = "VIEW"
	ConsentModeStore  = "STORE"
	ConsentModeQUERY  = "QUERY"
	ConsentModeSTREAM = "STREAM"
	// Fetch Type
	FetchTypeOnetime  = "ONETIME"
	FetchTypePeriodic = "PERIODIC"
	// Consent Types
	ConsentTypesProfile     = "PROFILE"
	ConsentTypesTransaction = "TRANSACTIONS"
	ConsentTypesSummary     = "SUMMARY"
	// FI Types
	FiTypesDeposit             = "DEPOSIT"
	FiTypesMutualFunds         = "MUTUAL_FUNDS"
	FiTypesInsurancePolicies   = "INSURANCE_POLICIES"
	FiTypesTermDeposit         = "TERM_DEPOSIT"
	FiTypesRecurringDeposit    = "RECURRING_DEPOSIT"
	FiTypesSIP                 = "SIP"
	FiTypesGovernmentSecrities = "GOVT_SECURITIES"
	FiTypesEquities            = "EQUITIES"
	FiTypesBonds               = "BONDS"
	FiTypesDebentures          = "DEBENTURES"
	FiTypesETF                 = "ETF"
)

const (
	PurposeWealthManagement int = 101
	PurposeReportings
	PurposeStatement
	PurposeMonitoring
	PurposeOneTime
)

// Create Consent Types //
type createConsentRequestInput struct {
	Phone string `json:"phone"`
}

type createConsentResponseOutput struct {
	ConsentHandle string `json:"consent_handle"`
}

type setuCreateConsentResponse struct {
	Ver           string    `json:"ver"`
	Timestamp     time.Time `json:"timestamp"`
	Txnid         uuid.UUID `json:"txnid"`
	Customer      idType    `json:"Customer"`
	ConsentHandle uuid.UUID `json:"ConsentHandle"`
	ErrorMsg      string    `json:"errorMsg"`
	ErrorCode     string    `json:"errorCode"`
}

type setuCreateConsentRequest struct {
	Ver           string         `json:"ver"`
	Timestamp     string         `json:"timestamp"`
	Txnid         uuid.UUID      `json:"txnid"`
	ConsentDetail consentDetails `json:"ConsentDetail"`
	jwt.StandardClaims
}

type consentDetails struct {
	ConsentStart  string        `json:"consentStart"`
	ConsentExpiry string        `json:"consentExpiry"`
	ConsentMode   string        `json:"consentMode"`
	FetchType     string        `json:"fetchType"`
	ConsentTypes  []string      `json:"consentTypes"`
	FiTypes       []string      `json:"fiTypes"`
	DataConsumer  idType        `json:"DataConsumer"`
	Customer      idType        `json:"Customer"`
	Purpose       purpose       `json:"Purpose"`
	FIDataRange   fIDataRange   `json:"FIDataRange"`
	DataLife      dataTimeRange `json:"DataLife"`
	Frequency     dataTimeRange `json:"Frequency"`
	DataFilter    []dataFilter  `json:"DataFilter"`
}

type idType struct {
	Id string `json:"id"`
}

type purpose struct {
	Code     string          `json:"code"`
	RefUri   string          `json:"refUri"`
	Text     string          `json:"text"`
	Category purposeCategory `json:"Category"`
}

type purposeCategory struct {
	Type string `json:"type"`
}

type fIDataRange struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type dataTimeRange struct {
	Unit  string `json:"unit"`
	Value int    `json:"value"`
}

type dataFilter struct {
	Type     string `json:"type"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

// Consent Artefact Status Types //
type setuConsentStatusRequest struct {
	Path string `json:"path"`
	jwt.StandardClaims
}

type setuConsentStatusResponse struct {
	Ver           string        `json:"ver"`
	Txnid         uuid.UUID     `json:"txnid"`
	Timestamp     time.Time     `json:"timestamp"`
	ConsentHandle string        `json:"ConsentHandle"`
	ConsentStatus consentStatus `json:"ConsentStatus"`
	ErrorMsg      string        `json:"errorMsg"`
	ErrorCode     string        `json:"errorCode"`
}

type consentStatus struct {
	Id     uuid.UUID `json:"id"`
	Status string    `json:"status"`
}

// Signed Consent Types //
type setuSignedConsentResponse struct {
	Ver             string     `json:"ver"`
	Txnid           uuid.UUID  `json:"txnid"`
	ConsentUse      consentUse `json:"ConsentUse"`
	ConsentId       uuid.UUID  `json:"consentId"`
	CreateTimestamp time.Time  `json:"createTimestamp"`
	SignedConsent   string     `json:"signedConsent"`
	Timestamp       time.Time  `json:"timestamp"`
	Status          string     `json:"status"`
	ErrorMsg        string     `json:"errorMsg"`
	ErrorCode       string     `json:"errorCode"`
}

type consentUse struct {
	Count           int       `json:"count"`
	LastUseDateTime uuid.UUID `json:"lastUseDateTime"`
	LogUri          string    `json:"logUri"`
}

// Rahasya Get Keys Types //
type rahasyaKeyResponse struct {
	PrivateKey  string             `json:"privateKey"`
	ErrorInfo   string             `json:"errorInfo"`
	KeyMaterial rahasyaKeyMaterial `json:"KeyMaterial"`
}

type rahasyaKeyMaterial struct {
	CryptoAlg   string          `json:"cryptoAlg"`
	Curve       string          `json:"curve"`
	Params      string          `json:"params"`
	DHPublicKey rahasyaDhPublic `json:"DHPublicKey"`
	Nonce       string          `json:"Nonce"`
}

type rahasyaDhPublic struct {
	Expiry     time.Time `json:"expiry"`
	Parameters string    `json:"Parameters"`
	KeyValue   string    `json:"KeyValue"`
}

// Fi Session //
type setuFiSessionRequest struct {
	Ver         string             `json:"ver"`
	Timestamp   string             `json:"timestamp"`
	Txnid       uuid.UUID          `json:"txnid"`
	FIDataRange fIDataRange        `json:"FIDataRange"`
	Consent     fiConsent          `json:"Consent"`
	KeyMaterial rahasyaKeyMaterial `json:"KeyMaterial"`
	jwt.StandardClaims
}

type setuFiSessionResponse struct {
	Ver       string    `json:"ver"`
	Timestamp time.Time `json:"timestamp"`
	Txnid     uuid.UUID `json:"txnid"`
	SessionId uuid.UUID `json:"sessionId"`
	ConsentId uuid.UUID `json:"consentId"`
}

type fiConsent struct {
	Id               uuid.UUID `json:"id"`
	DigitalSignature string    `json:"digitalSignature"`
}

// Fi Data //
type setuFiDataResponse struct {
	Ver       string             `json:"ver"`
	Timestamp string             `json:"timestamp"`
	Txnid     uuid.UUID          `json:"txnid"`
	FI        []fiEncryptionData `json:"FI"`
	jwt.StandardClaims
}

type fiEncryptionData struct {
	KeyMaterial rahasyaKeyMaterial `json:"KeyMaterial"`
	Data        []fiData           `json:"data"`
	FipId       string             `json:"fipId"`
}

type fiData struct {
	EncryptedFI     string    `json:"encryptedFI"`
	LinkRefNumber   uuid.UUID `json:"linkRefNumber"`
	MaskedAccNumber string    `json:"maskedAccNumber"`
}

// Rahasya Decrypt Types //
type rahasyaDecryptRequest struct {
	Base64Data        string             `json:"base64Data"`
	Base64RemoteNonce string             `json:"base64RemoteNonce"`
	Base64YourNonce   string             `json:"base64YourNonce"`
	OurPrivateKey     string             `json:"ourPrivateKey"`
	RemoteKeyMaterial rahasyaKeyMaterial `json:"remoteKeyMaterial"`
}

type rahasyaDataResponse struct {
	Base64Data string `json:"base64Data,omitempty"`
	Data       string `json:"data"`
	ErrorInfo  string `json:"errorInfo"`
}

type rahasyaDataResponseCollection struct {
	RahasyaData []rahasyaDataResponse
	FipId       string
}
