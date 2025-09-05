package disjunctions_of_refs;

import com.fasterxml.jackson.core.JsonGenerator;
import com.fasterxml.jackson.databind.JsonSerializer;
import com.fasterxml.jackson.databind.SerializerProvider;

import java.io.IOException;

public class DisjunctionOfRefsSerializer extends JsonSerializer<DisjunctionOfRefs> {

    @Override
    public void serialize(DisjunctionOfRefs value, JsonGenerator gen, SerializerProvider serializerProvider) throws IOException {
        if (value.myRefA != null) {
            gen.writeObject(value.myRefA);
        }
        else if (value.myRefB != null) {
            gen.writeObject(value.myRefB);
        }
    }
}
