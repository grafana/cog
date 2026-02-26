package default_disjunction_value;

import java.util.Map;

public class ValueB {
    public String type;
    public Map<String, Long> aMap;
    public Int64OrStringOrBool def;
    public ValueB() {
    }
    public ValueB(String type,Map<String, Long> aMap,Int64OrStringOrBool def) {
        this.type = type;
        this.aMap = aMap;
        this.def = def;
    }
}
