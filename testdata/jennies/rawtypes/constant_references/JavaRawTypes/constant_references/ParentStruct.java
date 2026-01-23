package constant_references;


public class ParentStruct {
    public Enum myEnum;
    public ParentStruct() {
        this.myEnum = Enum.VALUE_A;
    }
    public ParentStruct(Enum myEnum) {
        this.myEnum = myEnum;
    }
}
