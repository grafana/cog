package constraints;

import java.util.HashMap;
import java.util.LinkedList;
import java.util.Map;
import java.util.List;

public class RefStruct {
    public Map<String, String> labels;
    public List<String> tags;
    public RefStruct() {
        this.labels = new HashMap<>();
        this.tags = new LinkedList<>();
    }
    public RefStruct(Map<String, String> labels,List<String> tags) {
        this.labels = labels;
        this.tags = tags;
    }
}
