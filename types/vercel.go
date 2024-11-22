package types

type CreateVercelProjectResponse struct {
	AccountID                            string                `json:"accountId"`
	Analytics                            Analytics             `json:"analytics"`
	AutoAssignCustomDomains              bool                  `json:"autoAssignCustomDomains"`
	AutoAssignCustomDomainsUpdatedBy     string                `json:"autoAssignCustomDomainsUpdatedBy"`
	AutoExposeSystemEnvs                 bool                  `json:"autoExposeSystemEnvs"`
	BuildCommand                         *string               `json:"buildCommand"`
	CommandForIgnoringBuildStep          *string               `json:"commandForIgnoringBuildStep"`
	ConcurrencyBucketName                string                `json:"concurrencyBucketName"`
	ConnectBuildsEnabled                 bool                  `json:"connectBuildsEnabled"`
	ConnectConfigurationID               *string               `json:"connectConfigurationId"`
	CreatedAt                            int64                 `json:"createdAt"`
	Crons                                Crons                 `json:"crons"`
	CustomEnvironments                   []interface{}         `json:"customEnvironments"`
	CustomerSupportCodeVisibility        bool                  `json:"customerSupportCodeVisibility"`
	DataCache                            DataCache             `json:"dataCache"`
	DeploymentExpiration                 *DeploymentExpiration `json:"deploymentExpiration"`
	DevCommand                           *string               `json:"devCommand"`
	DirectoryListing                     bool                  `json:"directoryListing"`
	EnableAffectedProjectsDeployments    bool                  `json:"enableAffectedProjectsDeployments"`
	EnablePreviewFeedback                *bool                 `json:"enablePreviewFeedback"`
	EnableProductionFeedback             *bool                 `json:"enableProductionFeedback"`
	EnvVars                              []EnvVar              `json:"env"`
	Framework                            *string               `json:"framework"`
	GitComments                          GitComments           `json:"gitComments"`
	GitForkProtection                    bool                  `json:"gitForkProtection"`
	GitLFS                               bool                  `json:"gitLFS"`
	HasActiveBranches                    bool                  `json:"hasActiveBranches"`
	HasFloatingAliases                   bool                  `json:"hasFloatingAliases"`
	ID                                   string                `json:"id"`
	InstallCommand                       *string               `json:"installCommand"`
	IpBuckets                            []IpBucket            `json:"ipBuckets"`
	LastAliasRequest                     *AliasRequest         `json:"lastAliasRequest"`
	LastRollbackTarget                   *RollbackTarget       `json:"lastRollbackTarget"`
	LatestDeployments                    []Deployment          `json:"latestDeployments"`
	Live                                 bool                  `json:"live"`
	Microfrontends                       Microfrontends        `json:"microfrontends"`
	NodeVersion                          string                `json:"nodeVersion"`
	OIDCTokenConfig                      OIDCTokenConfig       `json:"oidcTokenConfig"`
	OutputDirectory                      *string               `json:"outputDirectory"`
	Paused                               bool                  `json:"paused"`
	Permissions                          Permissions           `json:"permissions"`
	ProductionDeploymentsFastLane        bool                  `json:"productionDeploymentsFastLane"`
	ProtectionBypass                     interface{}           `json:"protectionBypass"`
	PublicSource                         *bool                 `json:"publicSource"`
	ResourceConfig                       ResourceConfig        `json:"resourceConfig"`
	RootDirectory                        *string               `json:"rootDirectory"`
	Security                             Security              `json:"security"`
	ServerlessFunctionRegion             *string               `json:"serverlessFunctionRegion"`
	ServerlessFunctionZeroConfigFailover bool                  `json:"serverlessFunctionZeroConfigFailover"`
	SkewProtectionBoundaryAt             int64                 `json:"skewProtectionBoundaryAt"`
	SkewProtectionMaxAge                 int64                 `json:"skewProtectionMaxAge"`
	SkipGitConnectDuringLink             bool                  `json:"skipGitConnectDuringLink"`
	SourceFilesOutsideRootDirectory      bool                  `json:"sourceFilesOutsideRootDirectory"`
	SpeedInsights                        SpeedInsights         `json:"speedInsights"`
	SSOProtection                        *SSOProtection        `json:"ssoProtection"`
	Targets                              []Target              `json:"targets"`
	TransferCompletedAt                  int64                 `json:"transferCompletedAt"`
	TransferStartedAt                    int64                 `json:"transferStartedAt"`
	TransferToAccountID                  string                `json:"transferToAccountId"`
	TransferredFromAccountID             string                `json:"transferredFromAccountId"`
	TrustedIps                           []TrustedIP           `json:"trustedIps"`
	UpdatedAt                            int64                 `json:"updatedAt"`
	WebAnalytics                         WebAnalytics          `json:"webAnalytics"`
}

type UpdateVercelProjectResponse struct {
	AccountID                            string                `json:"accountId" form:"accountId" schema:"accountId"`
	Analytics                            Analytics             `json:"analytics" form:"analytics" schema:"analytics"`
	AutoAssignCustomDomains              bool                  `json:"autoAssignCustomDomains" form:"autoAssignCustomDomains" schema:"autoAssignCustomDomains"`
	AutoAssignCustomDomainsUpdatedBy     string                `json:"autoAssignCustomDomainsUpdatedBy" form:"autoAssignCustomDomainsUpdatedBy" schema:"autoAssignCustomDomainsUpdatedBy"`
	AutoExposeSystemEnvs                 bool                  `json:"autoExposeSystemEnvs" form:"autoExposeSystemEnvs" schema:"autoExposeSystemEnvs"`
	BuildCommand                         *string               `json:"buildCommand,omitempty" form:"buildCommand" schema:"buildCommand"`
	CommandForIgnoringBuildStep          *string               `json:"commandForIgnoringBuildStep,omitempty" form:"commandForIgnoringBuildStep" schema:"commandForIgnoringBuildStep"`
	ConcurrencyBucketName                string                `json:"concurrencyBucketName" form:"concurrencyBucketName" schema:"concurrencyBucketName"`
	ConnectBuildsEnabled                 bool                  `json:"connectBuildsEnabled" form:"connectBuildsEnabled" schema:"connectBuildsEnabled"`
	ConnectConfigurationID               *string               `json:"connectConfigurationId,omitempty" form:"connectConfigurationId" schema:"connectConfigurationId"`
	CreatedAt                            int64                 `json:"createdAt" form:"createdAt" schema:"createdAt"`
	Crons                                Crons                 `json:"crons" form:"crons" schema:"crons"`
	CustomEnvironments                   []interface{}         `json:"customEnvironments,omitempty" form:"customEnvironments" schema:"customEnvironments"`
	CustomerSupportCodeVisibility        bool                  `json:"customerSupportCodeVisibility" form:"customerSupportCodeVisibility" schema:"customerSupportCodeVisibility"`
	DataCache                            DataCache             `json:"dataCache" form:"dataCache" schema:"dataCache"`
	DeploymentExpiration                 *DeploymentExpiration `json:"deploymentExpiration,omitempty" form:"deploymentExpiration" schema:"deploymentExpiration"`
	DevCommand                           *string               `json:"devCommand,omitempty" form:"devCommand" schema:"devCommand"`
	DirectoryListing                     bool                  `json:"directoryListing" form:"directoryListing" schema:"directoryListing"`
	EnableAffectedProjectsDeployments    bool                  `json:"enableAffectedProjectsDeployments" form:"enableAffectedProjectsDeployments" schema:"enableAffectedProjectsDeployments"`
	EnablePreviewFeedback                *bool                 `json:"enablePreviewFeedback,omitempty" form:"enablePreviewFeedback" schema:"enablePreviewFeedback"`
	EnableProductionFeedback             *bool                 `json:"enableProductionFeedback,omitempty" form:"enableProductionFeedback" schema:"enableProductionFeedback"`
	EnvVars                              []EnvVar              `json:"env" form:"env" schema:"env"`
	Framework                            *string               `json:"framework,omitempty" form:"framework" schema:"framework"`
	GitComments                          GitComments           `json:"gitComments" form:"gitComments" schema:"gitComments"`
	GitForkProtection                    bool                  `json:"gitForkProtection" form:"gitForkProtection" schema:"gitForkProtection"`
	GitLFS                               bool                  `json:"gitLFS" form:"gitLFS" schema:"gitLFS"`
	HasActiveBranches                    bool                  `json:"hasActiveBranches" form:"hasActiveBranches" schema:"hasActiveBranches"`
	HasFloatingAliases                   bool                  `json:"hasFloatingAliases" form:"hasFloatingAliases" schema:"hasFloatingAliases"`
	ID                                   string                `json:"id" form:"id" schema:"id"`
	InstallCommand                       *string               `json:"installCommand,omitempty" form:"installCommand" schema:"installCommand"`
	IpBuckets                            []IpBucket            `json:"ipBuckets" form:"ipBuckets" schema:"ipBuckets"`
	LastAliasRequest                     *AliasRequest         `json:"lastAliasRequest,omitempty" form:"lastAliasRequest" schema:"lastAliasRequest"`
	LastRollbackTarget                   *RollbackTarget       `json:"lastRollbackTarget,omitempty" form:"lastRollbackTarget" schema:"lastRollbackTarget"`
	LatestDeployments                    []Deployment          `json:"latestDeployments" form:"latestDeployments" schema:"latestDeployments"`
	Live                                 bool                  `json:"live" form:"live" schema:"live"`
	Microfrontends                       Microfrontends        `json:"microfrontends" form:"microfrontends" schema:"microfrontends"`
	NodeVersion                          string                `json:"nodeVersion" form:"nodeVersion" schema:"nodeVersion"`
	OIDCTokenConfig                      OIDCTokenConfig       `json:"oidcTokenConfig" form:"oidcTokenConfig" schema:"oidcTokenConfig"`
	OutputDirectory                      *string               `json:"outputDirectory,omitempty" form:"outputDirectory" schema:"outputDirectory"`
	Paused                               bool                  `json:"paused" form:"paused" schema:"paused"`
	Permissions                          Permissions           `json:"permissions" form:"permissions" schema:"permissions"`
	ProductionDeploymentsFastLane        bool                  `json:"productionDeploymentsFastLane" form:"productionDeploymentsFastLane" schema:"productionDeploymentsFastLane"`
	ProtectionBypass                     interface{}           `json:"protectionBypass" form:"protectionBypass" schema:"protectionBypass"`
	PublicSource                         *bool                 `json:"publicSource,omitempty" form:"publicSource" schema:"publicSource"`
	ResourceConfig                       ResourceConfig        `json:"resourceConfig" form:"resourceConfig" schema:"resourceConfig"`
	RootDirectory                        *string               `json:"rootDirectory,omitempty" form:"rootDirectory" schema:"rootDirectory"`
	Security                             Security              `json:"security" form:"security" schema:"security"`
	ServerlessFunctionRegion             *string               `json:"serverlessFunctionRegion,omitempty" form:"serverlessFunctionRegion" schema:"serverlessFunctionRegion"`
	ServerlessFunctionZeroConfigFailover bool                  `json:"serverlessFunctionZeroConfigFailover" form:"serverlessFunctionZeroConfigFailover" schema:"serverlessFunctionZeroConfigFailover"`
	SkewProtectionBoundaryAt             int64                 `json:"skewProtectionBoundaryAt" form:"skewProtectionBoundaryAt" schema:"skewProtectionBoundaryAt"`
	SkewProtectionMaxAge                 int64                 `json:"skewProtectionMaxAge" form:"skewProtectionMaxAge" schema:"skewProtectionMaxAge"`
	SkipGitConnectDuringLink             bool                  `json:"skipGitConnectDuringLink" form:"skipGitConnectDuringLink" schema:"skipGitConnectDuringLink"`
	SourceFilesOutsideRootDirectory      bool                  `json:"sourceFilesOutsideRootDirectory" form:"sourceFilesOutsideRootDirectory" schema:"sourceFilesOutsideRootDirectory"`
	SpeedInsights                        SpeedInsights         `json:"speedInsights" form:"speedInsights" schema:"speedInsights"`
	SSOProtection                        *SSOProtection        `json:"ssoProtection,omitempty" form:"ssoProtection" schema:"ssoProtection"`
	Targets                              []Target              `json:"targets" form:"targets" schema:"targets"`
	Tier                                 string                `json:"tier" form:"tier" schema:"tier"`
	TransferCompletedAt                  int64                 `json:"transferCompletedAt" form:"transferCompletedAt" schema:"transferCompletedAt"`
	TransferStartedAt                    int64                 `json:"transferStartedAt" form:"transferStartedAt" schema:"transferStartedAt"`
	TransferToAccountID                  string                `json:"transferToAccountId" form:"transferToAccountId" schema:"transferToAccountId"`
	TransferredFromAccountID             string                `json:"transferredFromAccountId" form:"transferredFromAccountId" schema:"transferredFromAccountId"`
	TrustedIps                           []TrustedIP           `json:"trustedIps" form:"trustedIps" schema:"trustedIps"`
	UpdatedAt                            int64                 `json:"updatedAt" form:"updatedAt" schema:"updatedAt"`
	WebAnalytics                         WebAnalytics          `json:"webAnalytics" form:"webAnalytics" schema:"webAnalytics"`
}

type GetVercelEnvironmentVariablesResponse struct {
	EnvsArray []EnvVariable `json:"envsarrayRequired"`
}

type Analytics struct {
	CanceledAt          *int64 `json:"canceledAt"`
	DisabledAt          int64  `json:"disabledAt"`
	EnabledAt           int64  `json:"enabledAt"`
	ID                  string `json:"id"`
	PaidAt              *int64 `json:"paidAt"`
	SampleRatePercent   *int64 `json:"sampleRatePercent"`
	SpendLimitInDollars *int64 `json:"spendLimitInDollars"`
}

type Crons struct {
	Definitions []CronDefinition `json:"definitions"`
}

type CronDefinition struct {
	Host         string `json:"host"`
	Path         string `json:"path"`
	Schedule     string `json:"schedule"`
	DeploymentID string `json:"deploymentId"`
	DisabledAt   *int64 `json:"disabledAt"`
	EnabledAt    int64  `json:"enabledAt"`
	UpdatedAt    int64  `json:"updatedAt"`
}

type DataCache struct {
	StorageSizeBytes *int64 `json:"storageSizeBytes"`
	Unlimited        bool   `json:"unlimited"`
	UserDisabled     bool   `json:"userDisabled"`
}

type DeploymentExpiration struct {
	DeploymentsToKeep        int64 `json:"deploymentsToKeep"`
	ExpirationDays           int64 `json:"expirationDays"`
	ExpirationDaysCanceled   int64 `json:"expirationDaysCanceled"`
	ExpirationDaysErrored    int64 `json:"expirationDaysErrored"`
	ExpirationDaysProduction int64 `json:"expirationDaysProduction"`
}

type EnvVar struct {
	Comment         *string      `json:"comment"`
	ConfigurationID *string      `json:"configurationId"`
	ContentHint     *ContentHint `json:"contentHint"`
	CreatedAt       int64        `json:"createdAt"`
	CreatedBy       *string      `json:"createdBy"`
	UpdatedAt       int64        `json:"updatedAt"`
	UpdatedBy       *string      `json:"updatedBy"`
	Value           string       `json:"value"`
	VSMValue        string       `json:"vsmValue"`
}

type ContentHint struct {
	StoreID string `json:"storeId"`
	Type    string `json:"type"`
}

type GitComments struct {
	OnCommit      bool `json:"onCommit"`
	OnPullRequest bool `json:"onPullRequest"`
}

type IpBucket struct {
	Bucket       string `json:"bucket"`
	SupportUntil int64  `json:"supportUntil"`
}

type AliasRequest struct {
	FromDeploymentID string `json:"fromDeploymentId"`
	JobStatus        string `json:"jobStatus"`
	RequestedAt      int64  `json:"requestedAt"`
	ToDeploymentID   string `json:"toDeploymentId"`
	Type             string `json:"type"`
}

type RollbackTarget struct {
	DeploymentID string `json:"deploymentId"`
}

type Deployment struct {
	Alias            []string         `json:"alias"`
	AliasAssigned    interface{}      `json:"aliasAssigned"`
	AliasError       *AliasError      `json:"aliasError"`
	AliasFinal       *string          `json:"aliasFinal"`
	AutomaticAliases []AutomaticAlias `json:"automaticAliases"`
	BranchMatcher    BranchMatcher    `json:"branchMatcher"`
	BuildingAt       int64            `json:"buildingAt"`
	Checks           Checks           `json:"checks"`
}

type AliasError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type AutomaticAlias struct {
	Pattern string `json:"pattern"`
	Type    string `json:"type"`
}

type BranchMatcher struct {
	Pattern string `json:"pattern"`
	Type    string `json:"type"`
}

type Checks struct {
	Conclusion string `json:"conclusion"`
	State      string `json:"state"`
}

type Microfrontends struct {
	DefaultRoute string   `json:"defaultRoute"`
	Enabled      bool     `json:"enabled"`
	GroupIDs     []string `json:"groupIds"`
	IsDefaultApp bool     `json:"isDefaultApp"`
	UpdatedAt    int64    `json:"updatedAt"`
}

type VercelProjectRequestBody struct {
	Name                                 string                `json:"name,omitempty" form:"name" schema:"name"`
	BuildCommand                         string                `json:"buildCommand,omitempty" form:"buildCommand" schema:"buildCommand"`
	CommandForIgnoringBuildStep          string                `json:"commandForIgnoringBuildStep,omitempty" form:"commandForIgnoringBuildStep" schema:"commandForIgnoringBuildStep"`
	DevCommand                           string                `json:"devCommand,omitempty" form:"devCommand" schema:"devCommand"`
	EnableAffectedProjectsDeployments    bool                  `json:"enableAffectedProjectsDeployments,omitempty" form:"enableAffectedProjectsDeployments" schema:"enableAffectedProjectsDeployments"`
	EnvironmentVariables                 []EnvironmentVariable `json:"environmentVariables,omitempty" form:"environmentVariables" schema:"environmentVariables"`
	Framework                            string                `json:"framework,omitempty" form:"framework" schema:"framework"`
	GitRepository                        GitRepository         `json:"gitRepository,omitempty" form:"gitRepository" schema:"gitRepository"`
	InstallCommand                       string                `json:"installCommand,omitempty" form:"installCommand" schema:"installCommand"`
	OIDCTokenConfig                      OIDCTokenConfig       `json:"oidcTokenConfig,omitempty" form:"oidcTokenConfig" schema:"oidcTokenConfig"`
	OutputDirectory                      string                `json:"outputDirectory,omitempty" form:"outputDirectory" schema:"outputDirectory"`
	PublicSource                         bool                  `json:"publicSource,omitempty" form:"publicSource" schema:"publicSource"`
	RootDirectory                        string                `json:"rootDirectory,omitempty" form:"rootDirectory" schema:"rootDirectory"`
	ServerlessFunctionRegion             string                `json:"serverlessFunctionRegion,omitempty" form:"serverlessFunctionRegion" schema:"serverlessFunctionRegion"`
	ServerlessFunctionZeroConfigFailover interface{}           `json:"serverlessFunctionZeroConfigFailover,omitempty" form:"serverlessFunctionZeroConfigFailover" schema:"serverlessFunctionZeroConfigFailover"`
	SkipGitConnectDuringLink             bool                  `json:"skipGitConnectDuringLink,omitempty" form:"skipGitConnectDuringLink" schema:"skipGitConnectDuringLink"`
}

type EnvironmentVariable struct {
	Key       string `json:"key,omitempty" form:"key" schema:"key"`
	Target    string `json:"target,omitempty" form:"target" schema:"target"`
	GitBranch string `json:"gitBranch,omitempty" form:"gitBranch" schema:"gitBranch"`
	Type      string `json:"type,omitempty" form:"type" schema:"type"`
	Value     string `json:"value,omitempty" form:"value" schema:"value"`
}

type GitRepository struct {
	Repo string `json:"repo,omitempty" form:"repo" schema:"repo"`
	Type string `json:"type,omitempty" form:"type" schema:"type"`
}

type OIDCTokenConfig struct {
	Enabled    bool   `json:"enabled,omitempty" form:"enabled" schema:"enabled"`
	IssuerMode string `json:"issuerMode,omitempty" form:"issuerMode" schema:"issuerMode"`
}

type Permissions struct {
	Monitoring []interface{} `json:"Monitoring"`
}

type ResourceConfig struct {
	TimeoutSeconds int64 `json:"timeoutSeconds"`
}

type Security struct {
	OAuth               string `json:"oauth"`
	OverrideFailedLogin bool   `json:"overrideFailedLogin"`
}

type SpeedInsights struct {
	Enabled bool `json:"enabled"`
}

type SSOProtection struct {
	SSOTimeout int64 `json:"ssoTimeout"`
	SSOEnabled bool  `json:"ssoEnabled"`
}

type Target struct {
	Name    string `json:"name"`
	Branch  string `json:"branch"`
	Pattern string `json:"pattern"`
	Type    string `json:"type"`
}

type TrustedIP struct {
	IP string `json:"ip"`
}

type WebAnalytics struct {
	ID            string `json:"id"`
	Enabled       bool   `json:"enabled"`
	Source        string `json:"source"`
	CustomMetrics string `json:"customMetrics"`
}

type Pagination struct {
	TotalCount int `json:"totalCount"`
}

type EnvVariable struct {
	Comment              string       `json:"comment"`
	ConfigurationID      *string      `json:"configurationId,omitempty"`
	ContentHint          *ContentHint `json:"contentHint,omitempty"`
	CreatedAt            int64        `json:"createdAt"`
	CreatedBy            *string      `json:"createdBy,omitempty"`
	CustomEnvironmentIDs []string     `json:"customEnvironmentIds"`
	Decrypted            bool         `json:"decrypted"`
	EdgeConfigID         *string      `json:"edgeConfigId,omitempty"`
	EdgeConfigTokenID    *string      `json:"edgeConfigTokenId,omitempty"`
	GitBranch            string       `json:"gitBranch"`
	ID                   string       `json:"id"`
	InternalContentHint  *ContentHint `json:"internalContentHint,omitempty"`
	Key                  string       `json:"key"`
	SunsetSecretID       *string      `json:"sunsetSecretId,omitempty"`
	System               bool         `json:"system"`
	Target               *Target      `json:"target"`
	Type                 string       `json:"type"`
	UpdatedAt            int64        `json:"updatedAt"`
	UpdatedBy            *string      `json:"updatedBy,omitempty"`
	Value                string       `json:"value"`
	VSMValue             string       `json:"vsmValue"`
}
