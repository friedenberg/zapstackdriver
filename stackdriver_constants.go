package zapstackdriver

const (
	SDKeyType        = "@type"
	SDValueErrorType = "type.googleapis.com/google.devtools.clouderrorreporting.v1beta1.ReportedErrorEvent"

	SDKeyServiceContext        = "serviceContext"
	SDKeyServiceContextService = "service"
	SDKeyServiceContextVersion = "version"

	SDKeySourceLocation         = "sourceLocation"
	SDKeySourceLocationFile     = "file"
	SDKeySourceLocationLine     = "line"
	SDKeySourceLocationFunction = "function"

	SDKeyErrorContextSourceReferences           = "sourceReferences"
	SDKeyErrorContextSourceReferencesRepository = "repository"
	SDKeyErrorContextSourceReferencesRevisionID = "revisionId"

	SDKeyErrorContextUser = "user"

	SDKeyErrorContextHTTPRequest                   = "httpRequest"
	SDKeyErrorContextHTTPRequestMethod             = "method"
	SDKeyErrorContextHTTPRequestURL                = "url"
	SDKeyErrorContextHTTPRequestRemoteIP           = "remoteIp"
	SDKeyErrorContextHTTPRequestUserAgent          = "userAgent"
	SDKeyErrorContextHTTPRequestReferrer           = "referrer"
	SDKeyErrorContextHTTPRequestResponseStatusCode = "responseStatusCode"

	SDKeyErrorContextReportLocation         = "reportLocation"
	SDKeyErrorContextReportLocationFilePath = "filePath"
	SDKeyErrorContextReportLineNumber       = "lineNumber"
	SDKeyErrorContextReportFunctionName     = "functionName"
)
