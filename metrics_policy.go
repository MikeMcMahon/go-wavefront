package wavefront

// MetricsPolicy represents the global metrics policy for a given Wavefront domain
type MetricsPolicy struct {
	PolicyRules        []PolicyRule `json:"policyRules,omitempty"`
	Customer           string       `json:"customer,omitempty"`
	UpdaterId          string       `json:"updaterId,omitempty"`
	UpdatedEpochMillis int          `json:"updatedEpochMillis,omitempty"`
}

type UpdateMetricsPolicyRequest struct {
	PolicyRules []PolicyRule `json:"policyRules,omitempty"`
}

type PolicyRule struct {
	Accounts    []string    `json:"accounts,omitempty"`
	UserGroups  []UserGroup `json:"userGroups,omitempty"`
	Roles       []string    `json:"roles,omitempty"`
	Name        string      `json:"name,omitempty"`
	Tags        []string    `json:"tags,omitempty"`
	Description string      `json:"description,omitempty"`
	Prefixes    []string    `json:"prefixes,omitempty"`
	TagsAnded   bool        `json:"tagsAnded,omitempty"`
	AccessType  string      `json:"accessType,omitempty"`
}

// MetricsPolicyAPI is used to perform MetricsPolicy-related operations against the Wavefront API
type MetricsPolicyAPI struct {
	// client is the Wavefront client used to perform Dashboard-related operations
	client Wavefronter
}

const baseMetricsPolicyPath = "/api/v2/metricspolicy"

// MetricsPolicyAPI is used to return a client for MetricsPolicy-related operations
func (c *Client) MetricsPolicyAPI() *MetricsPolicyAPI {
	return &MetricsPolicyAPI{client: c}
}

func (m *MetricsPolicyAPI) Get() (*MetricsPolicy, error) {
	metricsPolicy := MetricsPolicy{}
	err := doRest(
		"GET",
		baseMetricsPolicyPath,
		m.client,
		doDirectResponse(&metricsPolicy),
	)
	return &metricsPolicy, err
}

func (m *MetricsPolicyAPI) Update(policyRules *UpdateMetricsPolicyRequest) (*MetricsPolicy, error) {
	metricsPolicy := MetricsPolicy{}
	err := doRest(
		"POST",
		baseMetricsPolicyPath,
		m.client,
		doPayload(policyRules),
		doDirectResponse(&metricsPolicy),
	)
	return &metricsPolicy, err
}
