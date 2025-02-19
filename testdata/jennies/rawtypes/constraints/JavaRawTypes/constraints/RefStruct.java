package constraints;

import java.util.Map;
import java.util.List;

public class RefStruct {
    public Map<String, String> labels;
    public List<String> tags;
    public RefStruct() {
    }
    public RefStruct(Map<String, String> labels,List<String> tags) {
        this.labels = labels;
        this.tags = tags;
    }
}
