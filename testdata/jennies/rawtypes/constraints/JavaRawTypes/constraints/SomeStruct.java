package constraints;

import java.util.HashMap;
import java.util.LinkedList;
import java.util.Map;
import java.util.List;

public class SomeStruct {
    public Long id;
    public Long maybeId;
    public Long greaterThanZero;
    public Long negative;
    public String title;
    public Map<String, String> labels;
    public List<String> tags;
    public SomeStruct() {
        this.id = 0L;
        this.greaterThanZero = 0L;
        this.negative = 0L;
        this.title = "";
        this.labels = new HashMap<>();
        this.tags = new LinkedList<>();
    }
    public SomeStruct(Long id,Long maybeId,Long greaterThanZero,Long negative,String title,Map<String, String> labels,List<String> tags) {
        this.id = id;
        this.maybeId = maybeId;
        this.greaterThanZero = greaterThanZero;
        this.negative = negative;
        this.title = title;
        this.labels = labels;
        this.tags = tags;
    }
}
