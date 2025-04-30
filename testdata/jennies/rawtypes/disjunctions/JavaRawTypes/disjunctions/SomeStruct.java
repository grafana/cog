package disjunctions;


public class SomeStruct {
    public String type;
    public Object fieldAny;
    public SomeStruct() {
        this.type = "";
        this.fieldAny = new Object();
    }
    public SomeStruct(String type,Object fieldAny) {
        this.type = type;
        this.fieldAny = fieldAny;
    }
}
