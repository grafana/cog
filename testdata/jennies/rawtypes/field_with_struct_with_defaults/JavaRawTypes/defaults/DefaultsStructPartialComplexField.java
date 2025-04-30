package defaults;


public class DefaultsStructPartialComplexField {
    public String uid;
    public Long intVal;
    public DefaultsStructPartialComplexField() {
        this.uid = "";
        this.intVal = 0L;
    }
    public DefaultsStructPartialComplexField(String uid,Long intVal) {
        this.uid = uid;
        this.intVal = intVal;
    }
}
