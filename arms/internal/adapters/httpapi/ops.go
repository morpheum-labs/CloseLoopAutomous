package httpapi

import "net/http"

// opsSummary exposes lightweight operator metrics (schema + build + product lifecycle counts).
func (h *Handlers) opsSummary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	active, deleted, err := h.Autopilot.Products.CountLifecycle(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"schema_version_expected": h.ExpectedSchemaVersion,
		"build_version":           h.BuildVersion,
		"build_commit":            h.BuildCommit,
		"products": map[string]any{
			"active":  active,
			"deleted": deleted,
		},
	})
}
