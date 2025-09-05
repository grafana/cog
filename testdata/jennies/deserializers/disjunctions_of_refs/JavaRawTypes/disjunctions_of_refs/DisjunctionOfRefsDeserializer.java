package disjunctions_of_refs;

import com.fasterxml.jackson.core.JsonParser;
import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.DeserializationContext;
import com.fasterxml.jackson.databind.JsonDeserializer;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;

import java.io.IOException;

public class DisjunctionOfRefsDeserializer extends JsonDeserializer<DisjunctionOfRefs> {

    @Override
    public DisjunctionOfRefs deserialize(JsonParser jp, DeserializationContext cxt) throws IOException {
        ObjectMapper mapper = (ObjectMapper) jp.getCodec();
        JsonNode root = mapper.readTree(jp);
        
        DisjunctionOfRefs disjunctionOfRefs = new DisjunctionOfRefs();
        if (!root.has("type")) {
            throw new IOException("Cannot find discriminator for DisjunctionOfRefs");
        }
        String discriminator = root.get("type").asText();  
        
        switch (discriminator) {
        case "A":
            disjunctionOfRefs.myRefA = mapper.convertValue(root, MyRefA.class);
            break;
        case "B":
            disjunctionOfRefs.myRefB = mapper.convertValue(root, MyRefB.class);
            break;
        }
        
        return disjunctionOfRefs;
    }
}
