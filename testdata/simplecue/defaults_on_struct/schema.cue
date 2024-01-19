Enum: "a" | "b" | "c"

HeatmapColorOptions: {
    scheme: string
    fill: string
    exponent: float32
    min?: float32
    max?: float32
    enum: Enum
}
container: {
    color: HeatmapColorOptions | *{
        scheme: "Oranges"
        fill: "dark-orange"
        enum: Enum & (*"b" | _)
    }
}
