package sandbox

#SomethingType: "foo" | "bar" @cog(kind="enum", memberNames="fooVal|barVal")

#SubContainer: {
	lala: string
}

SomeStructFoo: {
	type: #SomethingType & "foo"
	value: string
	other: #SubContainer
}

SomeStructBar: {
	type: #SomethingType & "bar"
	value: string
	otherValue: int
}

#BucketAggregationType: "terms" | "filters" | "geohash_grid" | "date_histogram" | "histogram" | "nested"

#BaseBucketAggregation: {
	id:        string
	type:      #BucketAggregationType
	settings?: _
}
#BucketAggregationWithField: {
	#BaseBucketAggregation
	field?: string
}

#Histogram: {
	#BucketAggregationWithField
	type: #BucketAggregationType & {
		"histogram"
	}
	settings?: _
}
