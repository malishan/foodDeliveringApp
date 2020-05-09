package apicontext

import "context"

type ctxType string

const (
	// APICtx - defining a separate type to avoid colliding with basic type
	APICtx ctxType = "apiCtx"
)

const (
	CustomerRole = "customer"
	AdminRole    = "admin"
)

// APIContext contains context of client
type APIContext struct {
	RestaurantID  string // Restaurant ID for identification of restaurant
	BranchID      string // Branch ID for identification of a particular restaurant chain
	RoleID        string // To know client type
	Token         string // Token is the api token
	RequestID     string // RequestID - used to track logs across a request-response cycle
	CorrelationID string // CorrelationID - used to track logs across a user's session
	UserID        string
}

// AppContext is the combination of native context and APIContext
type AppContext struct {
	context.Context
	APIContext
}

// GetAPICtx returns the api context from the context provided
func GetAPICtx(ctx context.Context) (APIContext, bool) {
	if ctx == nil {
		return APIContext{}, false
	}
	apiCtx, exists := ctx.Value(APICtx).(APIContext)
	return apiCtx, exists
}

// WithAPICtx returns a new context with the api context provided
func WithAPICtx(ctx context.Context, apictx APIContext) context.Context {
	return context.WithValue(ctx, APICtx, apictx)
}

// UpgradeCtx embeds native context and APIContext
func UpgradeCtx(ctx context.Context) AppContext {
	var tContext AppContext
	apiCtx, _ := GetAPICtx(ctx)

	tContext.Context = ctx
	tContext.APIContext = apiCtx
	return tContext
}
