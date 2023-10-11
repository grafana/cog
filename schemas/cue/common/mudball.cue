package common

// TODO docs
AxisPlacement:      "auto" | "top" | "right" | "bottom" | "left" | "hidden" @cog(kind="enum")

// TODO docs
AxisColorMode:      "text" | "series"                                       @cog(kind="enum")

// TODO docs
VisibilityMode:     "auto" | "never" | "always"                             @cog(kind="enum")

// TODO docs
GraphDrawStyle:     "line" | "bars" | "points"                              @cog(kind="enum")

// TODO docs
GraphTransform:     "constant" | "negative-Y"                               @cog(kind="enum",memberNames="Constant|NegativeY")

// TODO docs
LineInterpolation:  "linear" | "smooth" | "stepBefore" | "stepAfter"        @cog(kind="enum")

// TODO docs
ScaleDistribution:  "linear" | "log" | "ordinal" | "symlog"                 @cog(kind="enum")

// TODO docs
GraphGradientMode:  "none" | "opacity" | "hue" | "scheme"                   @cog(kind="enum")

// TODO docs
StackingMode:       "none" | "normal" | "percent"                           @cog(kind="enum")

// TODO docs
BarAlignment:       -1 | 0 | 1                                              @cog(kind="enum",memberNames="Before|Center|After")

// TODO docs
ScaleOrientation:   0 | 1                                                   @cog(kind="enum",memberNames="Horizontal|Vertical")

// TODO docs
ScaleDirection:     1 | 1 | -1 | -1                                         @cog(kind="enum",memberNames="Up|Right|Down|Left")

// TODO docs
LineStyle: {
	fill?: "solid" | "dash" | "dot" | "square"
	dash?: [...number]
} @cuetsy(kind="interface")

// TODO docs
LineConfig: {
	lineColor?:         string
	lineWidth?:         number
	lineInterpolation?: LineInterpolation
	lineStyle?:         LineStyle

	// Indicate if null values should be treated as gaps or connected.
	// When the value is a number, it represents the maximum delta in the
	// X axis that should be considered connected.  For timeseries, this is milliseconds
	spanNulls?:         bool | number
} @cuetsy(kind="interface")

// TODO docs
BarConfig: {
	barAlignment?:   BarAlignment
	barWidthFactor?: number
	barMaxWidth?:    number
} @cuetsy(kind="interface")

// TODO docs
FillConfig: {
	fillColor?:   string
	fillOpacity?: number
	fillBelowTo?: string
} @cuetsy(kind="interface")

// TODO docs
PointsConfig: {
	showPoints?:  VisibilityMode
	pointSize?:   number
	pointColor?:  string
	pointSymbol?: string
} @cuetsy(kind="interface")

// TODO docs
ScaleDistributionConfig: {
	type: ScaleDistribution
	log?: number
	linearThreshold?: number
} @cuetsy(kind="interface")

// TODO docs
AxisConfig: {
	axisPlacement?:     AxisPlacement
	axisColorMode?:     AxisColorMode
	axisLabel?:         string
	axisWidth?:         number
	axisSoftMin?:       number
	axisSoftMax?:       number
	axisGridShow?:      bool
	scaleDistribution?: ScaleDistributionConfig
	axisCenteredZero?:   bool
} @cuetsy(kind="interface")

// TODO docs
HideSeriesConfig: {
	tooltip: bool
	legend:  bool
	viz:     bool
} @cuetsy(kind="interface")

// TODO docs
StackingConfig: {
	mode?:  StackingMode
	group?: string
} @cuetsy(kind="interface")

// TODO docs
StackableFieldConfig: {
	stacking?: StackingConfig
} @cuetsy(kind="interface")

// TODO docs
HideableFieldConfig: {
	hideFrom?: HideSeriesConfig
} @cuetsy(kind="interface")

// TODO docs
GraphTresholdsStyleMode: "off" | "line" | "dashed" | "area" | "line+area" | "dashed+area" | "series" @cog(kind="enum",memberNames="Off|Line|Dashed|Area|LineAndArea|DashedAndArea|Series")

// TODO docs
GraphThresholdsStyleConfig: {
	mode: GraphTresholdsStyleMode
} @cuetsy(kind="interface")

// TODO docs
LegendPlacement: "bottom" | "right" @cuetsy(kind="type")

// TODO docs
// Note: "hidden" needs to remain as an option for plugins compatibility
LegendDisplayMode: "list" | "table" | "hidden" @cog(kind="enum")

// TODO docs
SingleStatBaseOptions: {
	OptionsWithTextFormatting
	reduceOptions: ReduceDataOptions
	orientation:   VizOrientation
} @cuetsy(kind="interface")

// TODO docs
ReduceDataOptions: {
	// If true show each row value
	values?: bool
	// if showing all values limit
	limit?: number
	// When !values, pick one value for the whole field
	calcs: [...string]
	// Which fields to show.  By default this is only numeric fields
	fields?: string
} @cuetsy(kind="interface")

// TODO docs
VizOrientation: "auto" | "vertical" | "horizontal" @cog(kind="enum")

// TODO docs
OptionsWithTooltip: {
	tooltip: VizTooltipOptions
} @cuetsy(kind="interface")

// TODO docs
OptionsWithLegend: {
	legend: VizLegendOptions
} @cuetsy(kind="interface")

// TODO docs
OptionsWithTimezones: {
	timezone?: [...string]
} @cuetsy(kind="interface")

// TODO docs
OptionsWithTextFormatting: {
	text?: VizTextDisplayOptions
} @cuetsy(kind="interface")

// TODO docs
BigValueColorMode: "value" | "background" | "background_solid" | "none" @cog(kind="enum", memberNames="Value|Background|BackgroundSolid|None")

// TODO docs
BigValueGraphMode: "none" | "line" | "area" @cog(kind="enum")

// TODO docs
BigValueJustifyMode: "auto" | "center" @cog(kind="enum")

// TODO docs
BigValueTextMode: "auto" | "value" | "value_and_name" | "name" | "none" @cog(kind="enum",memberNames="Auto|Value|ValueAndName|Name|None")

// TODO -- should not be table specific!
// TODO docs
FieldTextAlignment: "auto" | "left" | "right" | "center" @cuetsy(kind="type")

// Controls the value alignment in the TimelineChart component
TimelineValueAlignment: "center" | "left" | "right" @cuetsy(kind="type")

// TODO docs
VizTextDisplayOptions: {
	// Explicit title text size
	titleSize?: number
	// Explicit value text size
	valueSize?: number
} @cuetsy(kind="interface")

// TODO docs
TooltipDisplayMode: "single" | "multi" | "none" @cog(kind="enum")

// TODO docs
SortOrder: "asc" | "desc" | "none" @cog(kind="enum",memberNames="Ascending|Descending|None")

// TODO docs
GraphFieldConfig: {
	LineConfig
	FillConfig
	PointsConfig
	AxisConfig
	BarConfig
	StackableFieldConfig
	HideableFieldConfig
	drawStyle?:       GraphDrawStyle
	gradientMode?:    GraphGradientMode
	thresholdsStyle?: GraphThresholdsStyleConfig
	transform?:       GraphTransform
} @cuetsy(kind="interface")

// TODO docs
VizLegendOptions: {
	displayMode:  LegendDisplayMode
	placement:    LegendPlacement
	showLegend: 	bool
	asTable?:     bool
	isVisible?:   bool
	sortBy?:      string
	sortDesc?:    bool
	width?:       number
	calcs:        [...string]
} @cuetsy(kind="interface")

// Enum expressing the possible display modes
// for the bar gauge component of Grafana UI
BarGaugeDisplayMode: "basic" | "lcd" | "gradient" @cog(kind="enum")

// Allows for the table cell gauge display type to set the gauge mode.
BarGaugeValueMode: "color" | "text" | "hidden" @cog(kind="enum")

// TODO docs
VizTooltipOptions: {
	mode: TooltipDisplayMode
	sort: SortOrder
} @cuetsy(kind="interface")

Labels: {
	[string]: string
} @cuetsy(kind="interface")

// Compare two values
ComparisonOperation: "eq" | "neq" | "lt" | "lte" | "gt" | "gte" @cog(kind="enum",memberNames="EQ|NEQ|LT|LTE|GT|GTE")
