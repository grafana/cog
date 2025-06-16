package disjunctions;


public class BoolOrRef {
    protected Boolean bool;
    protected SomeStruct someStruct;
    protected BoolOrRef() {}
    public static BoolOrRef createBool(Boolean bool) {
        BoolOrRef boolOrRef = new BoolOrRef();
        boolOrRef.bool = bool;
        return boolOrRef;
    }
    public static BoolOrRef createSomeStruct(SomeStruct someStruct) {
        BoolOrRef boolOrRef = new BoolOrRef();
        boolOrRef.someStruct = someStruct;
        return boolOrRef;
    }
}
