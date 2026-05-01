package constraints;

import java.util.Objects;
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
    public String regex;
    public String negativeRegex;
    public List<String> minMaxList;
    public List<String> uniqueList;
    public List<Long> fullConstraintList;
    public SomeStruct() {
        this.id = 0L;
        this.greaterThanZero = 0L;
        this.negative = 0L;
        this.title = "";
        this.labels = new HashMap<>();
        this.tags = new LinkedList<>();
        this.regex = "";
        this.negativeRegex = "";
        this.minMaxList = new LinkedList<>();
        this.uniqueList = new LinkedList<>();
        this.fullConstraintList = new LinkedList<>();
    }
    public SomeStruct(Long id,Long maybeId,Long greaterThanZero,Long negative,String title,Map<String, String> labels,List<String> tags,String regex,String negativeRegex,List<String> minMaxList,List<String> uniqueList,List<Long> fullConstraintList) {
        this.id = id;
        this.maybeId = maybeId;
        this.greaterThanZero = greaterThanZero;
        this.negative = negative;
        this.title = title;
        this.labels = labels;
        this.tags = tags;
        this.regex = regex;
        this.negativeRegex = negativeRegex;
        this.minMaxList = minMaxList;
        this.uniqueList = uniqueList;
        this.fullConstraintList = fullConstraintList;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof SomeStruct)) return false;
        SomeStruct o = (SomeStruct) other;
        if (!Objects.equals(this.id, o.id)) return false;
        if (!Objects.equals(this.maybeId, o.maybeId)) return false;
        if (!Objects.equals(this.greaterThanZero, o.greaterThanZero)) return false;
        if (!Objects.equals(this.negative, o.negative)) return false;
        if (!Objects.equals(this.title, o.title)) return false;
        if (!Objects.equals(this.labels, o.labels)) return false;
        if (!Objects.equals(this.tags, o.tags)) return false;
        if (!Objects.equals(this.regex, o.regex)) return false;
        if (!Objects.equals(this.negativeRegex, o.negativeRegex)) return false;
        if (!Objects.equals(this.minMaxList, o.minMaxList)) return false;
        if (!Objects.equals(this.uniqueList, o.uniqueList)) return false;
        if (!Objects.equals(this.fullConstraintList, o.fullConstraintList)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.id, this.maybeId, this.greaterThanZero, this.negative, this.title, this.labels, this.tags, this.regex, this.negativeRegex, this.minMaxList, this.uniqueList, this.fullConstraintList);
    }
}
