package constant_references;


public class Struct {
    public String myValue;
    public Enum myEnum;
    public Struct() {
        this.myValue = "";
        this.myEnum = Enum.VALUE_A;
    }
    public Struct(String myValue,Enum myEnum) {
        this.myValue = myValue;
        this.myEnum = myEnum;
    }
}
