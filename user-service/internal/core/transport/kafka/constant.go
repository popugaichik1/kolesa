package core_kafka

// Define custom key types to avoid key collisions
type ContextKey string

const (
    OperationID ContextKey = "operationID" // For tracking, debugging
    OpUserID    ContextKey = "opUserID"  
)

const (
    noTimeout = -1
)

var FlushTimeOut = 5000