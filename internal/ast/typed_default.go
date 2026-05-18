package ast

// BuildTypedDefault converts a raw default value into a TypedDefault, dispatching
// on the declared type's Kind. For refs (where the referent cannot be resolved
// from the type alone), it falls back to shape inference via InferTypedDefault.
func BuildTypedDefault(declared Type, raw any) TypedDefault {
	switch declared.Kind {
	case KindArray:
		valueType := declared.AsArray().ValueType
		return TypedDefault{
			Kind:  KindArray,
			Array: &ArrayDefault{Items: rawToArrayItems(valueType, raw)},
		}
	case KindMap:
		valueType := declared.AsMap().ValueType
		return TypedDefault{
			Kind: KindMap,
			Map:  &MapDefault{Entries: rawToMapEntries(valueType, raw)},
		}
	case KindStruct:
		return TypedDefault{
			Kind:   KindStruct,
			Struct: &StructDefault{Fields: rawToStructFields(declared.AsStruct(), raw)},
		}
	case KindEnum:
		return TypedDefault{Kind: KindEnum, Enum: &EnumDefault{Value: raw}}
	case KindRef:
		return TypedDefault{Kind: KindRef, Ref: &RefDefault{Inner: InferTypedDefault(raw)}}
	default:
		return TypedDefault{Kind: KindScalar, Scalar: &ScalarDefault{Value: raw}}
	}
}

// InferTypedDefault best-effort infers the shape of a raw value when the
// declared type cannot be resolved (e.g. behind a Ref at parse time).
func InferTypedDefault(raw any) TypedDefault {
	switch v := raw.(type) {
	case map[string]any:
		fields := make(map[string]TypedDefault, len(v))
		for k, val := range v {
			fields[k] = InferTypedDefault(val)
		}
		return TypedDefault{Kind: KindStruct, Struct: &StructDefault{Fields: fields}}
	case []any:
		items := make([]TypedDefault, 0, len(v))
		for _, val := range v {
			items = append(items, InferTypedDefault(val))
		}
		return TypedDefault{Kind: KindArray, Array: &ArrayDefault{Items: items}}
	default:
		return TypedDefault{Kind: KindScalar, Scalar: &ScalarDefault{Value: raw}}
	}
}

func rawToArrayItems(valueType Type, raw any) []TypedDefault {
	rawItems, ok := raw.([]any)
	if !ok {
		return nil
	}
	items := make([]TypedDefault, 0, len(rawItems))
	for _, rawItem := range rawItems {
		items = append(items, BuildTypedDefault(valueType, rawItem))
	}
	return items
}

func rawToMapEntries(valueType Type, raw any) map[string]TypedDefault {
	rawEntries, ok := raw.(map[string]any)
	if !ok {
		return nil
	}
	entries := make(map[string]TypedDefault, len(rawEntries))
	for key, val := range rawEntries {
		entries[key] = BuildTypedDefault(valueType, val)
	}
	return entries
}

func rawToStructFields(structType StructType, raw any) map[string]TypedDefault {
	rawFields, ok := raw.(map[string]any)
	if !ok {
		return nil
	}
	known := make(map[string]Type, len(structType.Fields))
	for _, field := range structType.Fields {
		known[field.Name] = field.Type
	}
	fields := make(map[string]TypedDefault, len(rawFields))
	for name, rawVal := range rawFields {
		if declared, ok := known[name]; ok {
			fields[name] = BuildTypedDefault(declared, rawVal)
		} else {
			fields[name] = InferTypedDefault(rawVal)
		}
	}
	return fields
}

// TypedDefaultToAny converts a TypedDefault back into a raw `any` value. Used by
// jennies that still operate on raw values internally.
func TypedDefaultToAny(td TypedDefault) any {
	switch td.Kind {
	case KindScalar:
		if td.Scalar != nil {
			return td.Scalar.Value
		}
	case KindEnum:
		if td.Enum != nil {
			return td.Enum.Value
		}
	case KindArray:
		if td.Array == nil {
			return nil
		}
		items := make([]any, 0, len(td.Array.Items))
		for _, item := range td.Array.Items {
			items = append(items, TypedDefaultToAny(item))
		}
		return items
	case KindMap:
		if td.Map == nil {
			return nil
		}
		entries := make(map[string]any, len(td.Map.Entries))
		for key, val := range td.Map.Entries {
			entries[key] = TypedDefaultToAny(val)
		}
		return entries
	case KindStruct:
		if td.Struct == nil {
			return nil
		}
		fields := make(map[string]any, len(td.Struct.Fields))
		for key, val := range td.Struct.Fields {
			fields[key] = TypedDefaultToAny(val)
		}
		return fields
	case KindRef:
		if td.Ref != nil {
			return TypedDefaultToAny(td.Ref.Inner)
		}
	}
	return nil
}
