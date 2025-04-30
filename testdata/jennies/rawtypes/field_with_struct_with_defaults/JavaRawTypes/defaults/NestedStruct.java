package defaults;


public class NestedStruct {
    public String stringVal;
    public Long intVal;
    public NestedStruct() {
        this.stringVal = "";
        this.intVal = 0L;
    }
    public NestedStruct(String stringVal,Long intVal) {
        this.stringVal = stringVal;
        this.intVal = intVal;
    }
}
