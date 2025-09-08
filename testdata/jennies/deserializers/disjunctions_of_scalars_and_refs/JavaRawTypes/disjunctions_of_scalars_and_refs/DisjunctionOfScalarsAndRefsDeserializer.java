package disjunctions_of_scalars_and_refs;

import com.fasterxml.jackson.core.JsonParser;
import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.DeserializationContext;
import com.fasterxml.jackson.databind.JsonDeserializer;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;

import java.io.IOException;
import java.util.List;

public class DisjunctionOfScalarsAndRefsDeserializer extends JsonDeserializer<DisjunctionOfScalarsAndRefs> {

    @Override
    public DisjunctionOfScalarsAndRefs deserialize(JsonParser jp, DeserializationContext cxt) throws IOException {
        ObjectMapper mapper = (ObjectMapper) jp.getCodec();
        JsonNode root = mapper.readTree(jp);
        
        DisjunctionOfScalarsAndRefs disjunctionOfScalarsAndRefs = new DisjunctionOfScalarsAndRefs();
        if (root.isTextual()) {
            disjunctionOfScalarsAndRefs.string = mapper.convertValue(root, String.class);
        }
        else if (root.isBoolean()) {
            disjunctionOfScalarsAndRefs.bool = mapper.convertValue(root, Boolean.class);
        }
        else if (root.isArray()) {
            disjunctionOfScalarsAndRefs.arrayOfString = mapper.convertValue(root, new TypeReference<List<String>>() {});
        }
        else if (root.isObject() && couldBe(mapper, root, MyRefA.class)) {
            disjunctionOfScalarsAndRefs.myRefA = mapper.convertValue(root, MyRefA.class);
        }
        else if (root.isObject() && couldBe(mapper, root, MyRefB.class)) {
            disjunctionOfScalarsAndRefs.myRefB = mapper.convertValue(root, MyRefB.class);
        }
        else if (root.isObject()) {
            disjunctionOfScalarsAndRefs.any = mapper.convertValue(root, Object.class);
        }
        
        return disjunctionOfScalarsAndRefs;
    }
    
    private <T> boolean couldBe(ObjectMapper mapper, JsonNode root, Class<T> clazz) {
        try {
            mapper.convertValue(root, clazz);
        } catch (Exception e) {
           return false;
        }
        
        return true;
    }
}
