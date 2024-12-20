package disjunctions;


public class BoolOrRef {
    public Boolean bool;
    public SomeStruct someStruct;
    public BoolOrRef() {
    }
    
    public BoolOrRef(Boolean bool,SomeStruct someStruct) {
        this.bool = bool;
        this.someStruct = someStruct;
    }
}
