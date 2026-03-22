package convoy

// GraphEdge is one dependency arc: From must complete before To can run.
type GraphEdge struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// GraphLayer buckets subtasks sharing the same dag_layer (longest dependency depth).
type GraphLayer struct {
	Layer      int      `json:"layer"`
	SubtaskIDs []string `json:"subtask_ids"`
}

// GraphDetail is the payload for GET /api/convoys/{id}/graph (MC-style DAG view).
type GraphDetail struct {
	ConvoyID         string         `json:"convoy_id"`
	TopologicalOrder []string       `json:"topological_order"`
	Edges            []GraphEdge    `json:"edges"`
	Layers           []GraphLayer   `json:"layers"`
	GraphSummary     map[string]any `json:"graph"`
}
