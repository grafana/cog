package defaults;

import java.util.List;

public class DefaultsStructComplexField {
    public String uid;
    public DefaultsStructComplexFieldNested nested;
    public List<String> array;
    public DefaultsStructComplexField() {}
    
    public DefaultsStructComplexField(String uid,DefaultsStructComplexFieldNested nested,List<String> array) {
        this.uid = uid;
        this.nested = nested;
        this.array = array;
    }
}
