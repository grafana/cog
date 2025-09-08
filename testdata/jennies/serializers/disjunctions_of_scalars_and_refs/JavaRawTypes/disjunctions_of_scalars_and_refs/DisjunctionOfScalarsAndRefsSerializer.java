package disjunctions_of_scalars_and_refs;

import com.fasterxml.jackson.core.JsonGenerator;
import com.fasterxml.jackson.databind.JsonSerializer;
import com.fasterxml.jackson.databind.SerializerProvider;

import java.io.IOException;

public class DisjunctionOfScalarsAndRefsSerializer extends JsonSerializer<DisjunctionOfScalarsAndRefs> {

    @Override
    public void serialize(DisjunctionOfScalarsAndRefs value, JsonGenerator gen, SerializerProvider serializerProvider) throws IOException {
        if (value.string != null) {
            gen.writeObject(value.string);
        }
        else if (value.bool != null) {
            gen.writeObject(value.bool);
        }
        else if (value.arrayOfString != null) {
            gen.writeObject(value.arrayOfString);
        }
        else if (value.myRefA != null) {
            gen.writeObject(value.myRefA);
        }
        else if (value.myRefB != null) {
            gen.writeObject(value.myRefB);
        }
        else if (value.any != null) {
            gen.writeObject(value.any);
        }
    }
}
