package refs;


public class RefToSomeStruct {
    public Object fieldAny;
    public RefToSomeStruct() {
        this.fieldAny = new Object();
    }
    public RefToSomeStruct(Object fieldAny) {
        this.fieldAny = fieldAny;
    }
}
