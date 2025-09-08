package disjunctions_of_scalars;

import com.fasterxml.jackson.core.JsonGenerator;
import com.fasterxml.jackson.databind.JsonSerializer;
import com.fasterxml.jackson.databind.SerializerProvider;

import java.io.IOException;

public class DisjunctionOfScalarsSerializer extends JsonSerializer<DisjunctionOfScalars> {

    @Override
    public void serialize(DisjunctionOfScalars value, JsonGenerator gen, SerializerProvider serializerProvider) throws IOException {
        if (value.string != null) {
            gen.writeObject(value.string);
        }
        else if (value.bool != null) {
            gen.writeObject(value.bool);
        }
        else if (value.arrayOfString != null) {
            gen.writeObject(value.arrayOfString);
        }
        else if (value.float32 != null) {
            gen.writeObject(value.float32);
        }
        else if (value.int64 != null) {
            gen.writeObject(value.int64);
        }
    }
}
