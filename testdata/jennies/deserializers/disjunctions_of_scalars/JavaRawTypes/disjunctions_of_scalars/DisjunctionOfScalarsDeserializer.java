package disjunctions_of_scalars;

import com.fasterxml.jackson.core.JsonParser;
import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.DeserializationContext;
import com.fasterxml.jackson.databind.JsonDeserializer;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;

import java.io.IOException;
import java.util.List;

public class DisjunctionOfScalarsDeserializer extends JsonDeserializer<DisjunctionOfScalars> {

    @Override
    public DisjunctionOfScalars deserialize(JsonParser jp, DeserializationContext cxt) throws IOException {
        ObjectMapper mapper = (ObjectMapper) jp.getCodec();
        JsonNode root = mapper.readTree(jp);
        
        DisjunctionOfScalars disjunctionOfScalars = new DisjunctionOfScalars();
        if (root.isTextual()) {
            disjunctionOfScalars.string = mapper.convertValue(root, String.class);
        }
        else if (root.isBoolean()) {
            disjunctionOfScalars.bool = mapper.convertValue(root, Boolean.class);
        }
        else if (root.isArray()) {
            disjunctionOfScalars.arrayOfString = mapper.convertValue(root, new TypeReference<List<String>>() {});
        }
        else if (root.isFloatingPointNumber()) {
            disjunctionOfScalars.float32 = mapper.convertValue(root, Float.class);
        }
        else if (root.isIntegralNumber()) {
            disjunctionOfScalars.int64 = mapper.convertValue(root, Long.class);
        }
        
        return disjunctionOfScalars;
    }
}
