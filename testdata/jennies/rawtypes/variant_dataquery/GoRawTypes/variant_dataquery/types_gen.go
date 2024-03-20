package variant_dataquery

type Query struct {
	Expr string `json:"expr"`
	Instant *bool `json:"instant,omitempty"`
}
func (resource Query) ImplementsDataqueryVariant() {}


