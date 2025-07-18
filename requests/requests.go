// Package requests is tool for sending HTTP requests.
// Features:
// - Retry mechanism. Retry count, time between retries.
// - Client which provide basic settings to created requests. Nesting cookies, headers, etc.
// - Cancel action if context is canceled.
// - Export response to provided structure (JSON).
// - FormData writer.
// - Bytes writer.
package requests
