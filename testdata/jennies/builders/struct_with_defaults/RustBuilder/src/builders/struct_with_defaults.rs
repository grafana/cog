use crate::cog;
use crate::types::struct_with_defaults::NestedStruct;
use crate::types::struct_with_defaults::Struct;

#[derive(Debug, Clone)]
pub struct NestedStructBuilder {
    internal: NestedStruct,
    errors: Vec<cog::BuildError>,
}

impl NestedStructBuilder {
    pub fn new() -> Self {
        Self {
            internal: NestedStruct::default(),
            errors: Vec::new(),
        }
    }
}

impl Default for NestedStructBuilder {
    fn default() -> Self {
        Self::new()
    }
}

impl NestedStructBuilder {
    pub fn string_val(mut self, string_val: String) -> Self {
        self.internal.string_val = string_val;

        self
    }
}

impl NestedStructBuilder {
    pub fn int_val(mut self, int_val: i64) -> Self {
        self.internal.int_val = int_val;

        self
    }
}

impl cog::Builder<NestedStruct> for NestedStructBuilder {
    fn build(&self) -> Result<NestedStruct, Vec<cog::BuildError>> {
        if !self.errors.is_empty() {
            return Err(self.errors.clone());
        }

        Ok(self.internal.clone())
    }
}

#[derive(Debug, Clone)]
pub struct StructBuilder {
    internal: Struct,
    errors: Vec<cog::BuildError>,
}

impl StructBuilder {
    pub fn new() -> Self {
        Self {
            internal: Struct::default(),
            errors: Vec::new(),
        }
    }
}

impl Default for StructBuilder {
    fn default() -> Self {
        Self::new()
    }
}

impl StructBuilder {
    pub fn all_fields(mut self, all_fields: impl cog::Builder<NestedStruct>) -> Self {
        let built0 = match all_fields.build() {
            Ok(val) => val,
            Err(mut err) => {
                self.errors.append(&mut err);
                return self;
            }
        };
        self.internal.all_fields = built0;

        self
    }
}

impl StructBuilder {
    pub fn partial_fields(mut self, partial_fields: impl cog::Builder<NestedStruct>) -> Self {
        let built0 = match partial_fields.build() {
            Ok(val) => val,
            Err(mut err) => {
                self.errors.append(&mut err);
                return self;
            }
        };
        self.internal.partial_fields = built0;

        self
    }
}

impl StructBuilder {
    pub fn empty_fields(mut self, empty_fields: impl cog::Builder<NestedStruct>) -> Self {
        let built0 = match empty_fields.build() {
            Ok(val) => val,
            Err(mut err) => {
                self.errors.append(&mut err);
                return self;
            }
        };
        self.internal.empty_fields = built0;

        self
    }
}

impl StructBuilder {
    pub fn complex_field(mut self, complex_field: serde_json::Value) -> Self {
        self.internal.complex_field = complex_field;

        self
    }
}

impl StructBuilder {
    pub fn partial_complex_field(mut self, partial_complex_field: serde_json::Value) -> Self {
        self.internal.partial_complex_field = partial_complex_field;

        self
    }
}

impl cog::Builder<Struct> for StructBuilder {
    fn build(&self) -> Result<Struct, Vec<cog::BuildError>> {
        if !self.errors.is_empty() {
            return Err(self.errors.clone());
        }

        Ok(self.internal.clone())
    }
}
