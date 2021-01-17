package sq

const (
	// Class 00 - Successful Completion
	CodeSuccessfulCompletion = "00000"

	// Class 01 - Warning
	CodeWarning                          = "01000"
	CodeDynamicResultSetsReturned        = "0100C"
	CodeImplicitZeroBitPadding           = "01008"
	CodeNullValueEliminatedInSetFunction = "01003"
	CodePrivilegeNotGranted              = "01007"
	CodePrivilegeNotRevoked              = "01006"
	CodeStringDataRightTruncationWarning = "01004"
	CodeDeprecatedFeature                = "01P01"

	// Class 02 - No Data (this is also a warning class per the SQL standard)
	CodeNoData                                = "02000"
	CodeNoAdditionalDynamicResultSetsReturned = "02001"

	// Class 03 - SQL Statement Not Yet Complete
	CodeSQLStatementNotYetComplete = "03000"

	// Class 08 - Connection Exception
	CodeConnectionException                           = "08000"
	CodeConnectionDoesNotExist                        = "08003"
	CodeConnectionFailure                             = "08006"
	CodeSQLClientUnableToEstablishSqlconnection       = "08001"
	CodeSqlserverRejectedEstablishmentOfSqlconnection = "08004"
	CodeTransactionResolutionUnknown                  = "08007"
	CodeProtocolViolation                             = "08P01"

	// Class 09 - Triggered Action Exception
	CodeTriggeredActionException = "09000"

	// Class 0A - Feature Not Supported
	CodeFeatureNotSupported = "0A000"

	// Class 0B - Invalid Transaction Initiation
	CodeInvalidTransactionInitiation = "0B000"

	// Class 0F - Locator Exception
	CodeLocatorException            = "0F000"
	CodeInvalidLocatorSpecification = "0F001"

	// Class 0L - Invalid Grantor
	CodeInvalidGrantor        = "0L000"
	CodeInvalidGrantOperation = "0LP01"

	// Class 0P - Invalid Role Specification
	CodeInvalidRoleSpecification = "0P000"

	// Class 0Z - Diagnostics Exception
	CodeDiagnosticsException                           = "0Z000"
	CodeStackedDiagnosticsAccessedWithoutActiveHandler = "0Z002"

	// Class 20 - Case Not Found
	CodeCaseNotFound = "20000"

	// Class 21 - Cardinality Violation
	CodeCardinalityViolation = "21000"

	// Class 22 - Data Exception
	CodeDataException                         = "22000"
	CodeArraySubscriptError                   = "2202E"
	CodeCharacterNotInRepertoire              = "22021"
	CodeDatetimeFieldOverflow                 = "22008"
	CodeDivisionByZero                        = "22012"
	CodeErrorInAssignment                     = "22005"
	CodeEscapeCharacterConflict               = "2200B"
	CodeIndicatorOverflow                     = "22022"
	CodeIntervalFieldOverflow                 = "22015"
	CodeInvalidArgumentForLogarithm           = "2201E"
	CodeInvalidArgumentForNtileFunction       = "22014"
	CodeInvalidArgumentForNthValueFunction    = "22016"
	CodeInvalidArgumentForPowerFunction       = "2201F"
	CodeInvalidArgumentForWidthBucketFunction = "2201G"
	CodeInvalidCharacterValueForCast          = "22018"
	CodeInvalidDatetimeFormat                 = "22007"
	CodeInvalidEscapeCharacter                = "22019"
	CodeInvalidEscapeOctet                    = "2200D"
	CodeInvalidEscapeSequence                 = "22025"
	CodeNonstandardUseOfEscapeCharacter       = "22P06"
	CodeInvalidIndicatorParameterValue        = "22010"
	CodeInvalidParameterValue                 = "22023"
	CodeInvalidPrecedingOrFollowingSize       = "22013"
	CodeInvalidRegularExpression              = "2201B"
	CodeInvalidRowCountInLimitClause          = "2201W"
	CodeInvalidRowCountInResultOffsetClause   = "2201X"
	CodeInvalidTablesampleArgument            = "2202H"
	CodeInvalidTablesampleRepeat              = "2202G"
	CodeInvalidTimeZoneDisplacementValue      = "22009"
	CodeInvalidUseOfEscapeCharacter           = "2200C"
	CodeMostSpecificTypeMismatch              = "2200G"
	CodeNullValueNotAllowed                   = "22004"
	CodeNullValueNoIndicatorParameter         = "22002"
	CodeNumericValueOutOfRange                = "22003"
	CodeSequenceGeneratorLimitExceeded        = "2200H"
	CodeStringDataLengthMismatch              = "22026"
	CodeStringDataRightTruncation             = "22001"
	CodeSubstringError                        = "22011"
	CodeTrimError                             = "22027"
	CodeUnterminatedCString                   = "22024"
	CodeZeroLengthCharacterString             = "2200F"
	CodeFloatingPointException                = "22P01"
	CodeInvalidTextRepresentation             = "22P02"
	CodeInvalidBinaryRepresentation           = "22P03"
	CodeBadCopyFileFormat                     = "22P04"
	CodeUntranslatableCharacter               = "22P05"
	CodeNotAnXMLDocument                      = "2200L"
	CodeInvalidXMLDocument                    = "2200M"
	CodeInvalidXMLContent                     = "2200N"
	CodeInvalidXMLComment                     = "2200S"
	CodeInvalidXMLProcessingInstruction       = "2200T"
	CodeDuplicateJSONObjectKeyValue           = "22030"
	CodeInvalidJSONText                       = "22032"
	CodeInvalidSQLJSONSubscript               = "22033"
	CodeMoreThanOneSQLJSONItem                = "22034"
	CodeNoSQLJSONItem                         = "22035"
	CodeNonNumericSQLJSONItem                 = "22036"
	CodeNonUniqueKeysInAJSONObject            = "22037"
	CodeSingletonSQLJSONItemRequired          = "22038"
	CodeSQLJSONArrayNotFound                  = "22039"
	CodeSQLJSONMemberNotFound                 = "2203A"
	CodeSQLJSONNumberNotFound                 = "2203B"
	CodeSQLJSONObjectNotFound                 = "2203C"
	CodeTooManyJSONArrayElements              = "2203D"
	CodeTooManyJSONObjectMembers              = "2203E"
	CodeSQLJSONScalarRequired                 = "2203F"

	// Class 23 - Integrity Constraint Violation
	CodeIntegrityConstraintViolation = "23000"
	CodeRestrictViolation            = "23001"
	CodeNotNullViolation             = "23502"
	CodeForeignKeyViolation          = "23503"
	CodeUniqueViolation              = "23505"
	CodeCheckViolation               = "23514"
	CodeExclusionViolation           = "23P01"

	// Class 24 - Invalid Cursor State
	CodeInvalidCursorState = "24000"

	// Class 25 - Invalid Transaction State
	CodeInvalidTransactionState                         = "25000"
	CodeActiveSQLTransaction                            = "25001"
	CodeBranchTransactionAlreadyActive                  = "25002"
	CodeHeldCursorRequiresSameIsolationLevel            = "25008"
	CodeInappropriateAccessModeForBranchTransaction     = "25003"
	CodeInappropriateIsolationLevelForBranchTransaction = "25004"
	CodeNoActiveSQLTransactionForBranchTransaction      = "25005"
	CodeReadOnlySQLTransaction                          = "25006"
	CodeSchemaAndDataStatementMixingNotSupported        = "25007"
	CodeNoActiveSQLTransaction                          = "25P01"
	CodeInFailedSQLTransaction                          = "25P02"
	CodeIdleInTransactionSessionTimeout                 = "25P03"

	// Class 26 - Invalid SQL Statement Name
	CodeInvalidSQLStatementName = "26000"

	// Class 27 - Triggered Data Change Violation
	CodeTriggeredDataChangeViolation = "27000"

	// Class 28 - Invalid Authorization Specification
	CodeInvalidAuthorizationSpecification = "28000"
	CodeInvalidPassword                   = "28P01"

	// Class 2B - Dependent Privilege Descriptors Still Exist
	CodeDependentPrivilegeDescriptorsStillExist = "2B000"
	CodeDependentObjectsStillExist              = "2BP01"

	// Class 2D - Invalid Transaction Termination
	CodeInvalidTransactionTermination = "2D000"

	// Class 2F - SQL Routine Exception
	CodeSQLRoutineException               = "2F000"
	CodeFunctionExecutedNoReturnStatement = "2F005"
	CodeModifyingSQLDataNotPermitted      = "2F002"
	CodeProhibitedSQLStatementAttempted   = "2F003"
	CodeReadingSQLDataNotPermitted        = "2F004"

	// Class 34 - Invalid Cursor Name
	CodeInvalidCursorName = "34000"

	// Class 38 - External Routine Exception
	CodeExternalRoutineException                = "38000"
	CodeExternalContainingSQLNotPermitted       = "38001"
	CodeExternalModifyingSQLDataNotPermitted    = "38002"
	CodeExternalProhibitedSQLStatementAttempted = "38003"
	CodeExternalReadingSQLDataNotPermitted      = "38004"

	// Class 39 - External Routine Invocation Exception
	CodeExternalRoutineInvocationException   = "39000"
	CodeExternalInvalidSqlstateReturned      = "39001"
	CodeExternalNullValueNotAllowed          = "39004"
	CodeExternalTriggerProtocolViolated      = "39P01"
	CodeExternalSRFProtocolViolated          = "39P02"
	CodeExternalEventTriggerProtocolViolated = "39P03"

	// Class 3B - Savepoint Exception
	CodeSavepointException            = "3B000"
	CodeInvalidSavepointSpecification = "3B001"

	// Class 3D - Invalid Catalog Name
	CodeInvalidCatalogName = "3D000"

	// Class 3F - Invalid Schema Name
	CodeInvalidSchemaName = "3F000"

	// Class 40 - Transaction Rollback
	CodeTransactionRollback                     = "40000"
	CodeTransactionIntegrityConstraintViolation = "40002"
	CodeSerializationFailure                    = "40001"
	CodeStatementCompletionUnknown              = "40003"
	CodeDeadlockDetected                        = "40P01"

	// Class 42 - Syntax Error or Access Rule Violation
	CodeSyntaxErrorOrAccessRuleViolation   = "42000"
	CodeSyntaxError                        = "42601"
	CodeInsufficientPrivilege              = "42501"
	CodeCannotCoerce                       = "42846"
	CodeGroupingError                      = "42803"
	CodeWindowingError                     = "42P20"
	CodeInvalidRecursion                   = "42P19"
	CodeInvalidForeignKey                  = "42830"
	CodeInvalidName                        = "42602"
	CodeNameTooLong                        = "42622"
	CodeReservedName                       = "42939"
	CodeDatatypeMismatch                   = "42804"
	CodeIndeterminateDatatype              = "42P18"
	CodeCollationMismatch                  = "42P21"
	CodeIndeterminateCollation             = "42P22"
	CodeWrongObjectType                    = "42809"
	CodeGeneratedAlways                    = "428C9"
	CodeUndefinedColumn                    = "42703"
	CodeUndefinedFunction                  = "42883"
	CodeUndefinedTable                     = "42P01"
	CodeUndefinedParameter                 = "42P02"
	CodeUndefinedObject                    = "42704"
	CodeDuplicateColumn                    = "42701"
	CodeDuplicateCursor                    = "42P03"
	CodeDuplicateDatabase                  = "42P04"
	CodeDuplicateFunction                  = "42723"
	CodeDuplicatePreparedStatement         = "42P05"
	CodeDuplicateSchema                    = "42P06"
	CodeDuplicateTable                     = "42P07"
	CodeDuplicateAlias                     = "42712"
	CodeDuplicateObject                    = "42710"
	CodeAmbiguousColumn                    = "42702"
	CodeAmbiguousFunction                  = "42725"
	CodeAmbiguousParameter                 = "42P08"
	CodeAmbiguousAlias                     = "42P09"
	CodeInvalidColumnReference             = "42P10"
	CodeInvalidColumnDefinition            = "42611"
	CodeInvalidCursorDefinition            = "42P11"
	CodeInvalidDatabaseDefinition          = "42P12"
	CodeInvalidFunctionDefinition          = "42P13"
	CodeInvalidPreparedStatementDefinition = "42P14"
	CodeInvalidSchemaDefinition            = "42P15"
	CodeInvalidTableDefinition             = "42P16"
	CodeInvalidObjectDefinition            = "42P17"

	// Class 44 - WITH CHECK OPTION Violation
	CodeWithCheckOptionViolation = "44000"

	// Class 53 - Insufficient Resources
	CodeInsufficientResources      = "53000"
	CodeDiskFull                   = "53100"
	CodeOutOfMemory                = "53200"
	CodeTooManyConnections         = "53300"
	CodeConfigurationLimitExceeded = "53400"

	// Class 54 - Program Limit Exceeded
	CodeProgramLimitExceeded = "54000"
	CodeStatementTooComplex  = "54001"
	CodeTooManyColumns       = "54011"
	CodeTooManyArguments     = "54023"

	// Class 55 - Object Not In Prerequisite State
	CodeObjectNotInPrerequisiteState = "55000"
	CodeObjectInUse                  = "55006"
	CodeCantChangeRuntimeParam       = "55P02"
	CodeLockNotAvailable             = "55P03"
	CodeUnsafeNewEnumValueUsage      = "55P04"

	// Class 57 - Operator Intervention
	CodeOperatorIntervention = "57000"
	CodeQueryCanceled        = "57014"
	CodeAdminShutdown        = "57P01"
	CodeCrashShutdown        = "57P02"
	CodeCannotConnectNow     = "57P03"
	CodeDatabaseDropped      = "57P04"

	// Class 58 - System Error (errors external to PostgreSQL itself)
	CodeSystemError   = "58000"
	CodeIOError       = "58030"
	CodeUndefinedFile = "58P01"
	CodeDuplicateFile = "58P02"

	// Class 72 - Snapshot Failure
	CodeSnapshotTooOld = "72000"

	// Class F0 - Configuration File Error
	CodeConfigFileError = "F0000"
	CodeLockFileExists  = "F0001"

	// Class HV - Foreign Data Wrapper Error (SQL/MED)
	CodeFDWError                             = "HV000"
	CodeFDWColumnNameNotFound                = "HV005"
	CodeFDWDynamicParameterValueNeeded       = "HV002"
	CodeFDWFunctionSequenceError             = "HV010"
	CodeFDWInconsistentDescriptorInformation = "HV021"
	CodeFDWInvalidAttributeValue             = "HV024"
	CodeFDWInvalidColumnName                 = "HV007"
	CodeFDWInvalidColumnNumber               = "HV008"
	CodeFDWInvalidDataType                   = "HV004"
	CodeFDWInvalidDataTypeDescriptors        = "HV006"
	CodeFDWInvalidDescriptorFieldIdentifier  = "HV091"
	CodeFDWInvalidHandle                     = "HV00B"
	CodeFDWInvalidOptionIndex                = "HV00C"
	CodeFDWInvalidOptionName                 = "HV00D"
	CodeFDWInvalidStringLengthOrBufferLength = "HV090"
	CodeFDWInvalidStringFormat               = "HV00A"
	CodeFDWInvalidUseOfNullPointer           = "HV009"
	CodeFDWTooManyHandles                    = "HV014"
	CodeFDWOutOfMemory                       = "HV001"
	CodeFDWNoSchemas                         = "HV00P"
	CodeFDWOptionNameNotFound                = "HV00J"
	CodeFDWReplyHandle                       = "HV00K"
	CodeFDWSchemaNotFound                    = "HV00Q"
	CodeFDWTableNotFound                     = "HV00R"
	CodeFDWUnableToCreateExecution           = "HV00L"
	CodeFDWUnableToCreateReply               = "HV00M"
	CodeFDWUnableToEstablishConnection       = "HV00N"
	CodePLPGSQLError                         = "P0000"
	CodeRaiseException                       = "P0001"
	CodeNoDataFound                          = "P0002"
	CodeTooManyRows                          = "P0003"
	CodeAssertFailure                        = "P0004"
	CodeInternalError                        = "XX000"
	CodeDataCorrupted                        = "XX001"
	CodeIndexCorrupted                       = "XX002"
)
