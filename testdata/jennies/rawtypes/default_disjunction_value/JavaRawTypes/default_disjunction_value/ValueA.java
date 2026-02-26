package default_disjunction_value;

import java.util.List;

public class ValueA {
    public String type;
    public List<String> anArray;
    public ValueB otherRef;
    public ValueA() {
    }
    public ValueA(String type,List<String> anArray,ValueB otherRef) {
        this.type = type;
        this.anArray = anArray;
        this.otherRef = otherRef;
    }
}
