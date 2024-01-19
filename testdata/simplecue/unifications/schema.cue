#InlineScript: string | {
    inline?: string
}

#MetricAggregationWithInlineScript: {
    settings?: {
        script?: #InlineScript
    }
}

#Average: {
    #MetricAggregationWithInlineScript
    type: "avg"
    settings?: {
        script?: #InlineScript
        missing?: string
    }
}
