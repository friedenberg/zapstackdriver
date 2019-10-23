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
	SDKeyErrorContextSourceReferencesRevisionId = "revisionId"

	SDKeyErrorContextUser = "user"

	SDKeyErrorContextHttpRequest                   = "httpRequest"
	SDKeyErrorContextHttpRequestMethod             = "method"
	SDKeyErrorContextHttpRequestUrl                = "url"
	SDKeyErrorContextHttpRequestRemoteIp           = "remoteIp"
	SDKeyErrorContextHttpRequestUserAgent          = "userAgent"
	SDKeyErrorContextHttpRequestReferrer           = "referrer"
	SDKeyErrorContextHttpRequestResponseStatusCode = "responseStatusCode"

	SDKeyErrorContextReportLocation         = "reportLocation"
	SDKeyErrorContextReportLocationFilePath = "filePath"
	SDKeyErrorContextReportLineNumber       = "lineNumber"
	SDKeyErrorContextReportFunctionName     = "functionName"
)
