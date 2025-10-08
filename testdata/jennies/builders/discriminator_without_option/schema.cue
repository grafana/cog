package discriminator_without_option

AnEnum: "A" | "B"

NoShowFieldOption: {
	field: AnEnum & "A"
	text: string
}

ShowFieldOption: {
	field: AnEnum
	text: string
}
