package constant_references;


public class StructB {
    public Enum myEnum;
    public String myValue;
    public StructB() {
        this.myEnum = Enum.VALUE_B;
        this.myValue = "";
    }
    public StructB(String myValue) {
        this.myEnum = Enum.VALUE_B;
        this.myValue = myValue;
    }
}
