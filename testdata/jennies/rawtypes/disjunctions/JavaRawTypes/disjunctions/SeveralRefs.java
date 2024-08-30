package disjunctions;


public class SeveralRefs {
    public SomeStruct someStruct;
    public SomeOtherStruct someOtherStruct;
    public YetAnotherStruct yetAnotherStruct;
    public SeveralRefs() {}

    public SeveralRefs(SomeStruct someStruct,SomeOtherStruct someOtherStruct,YetAnotherStruct yetAnotherStruct) {
        this.someStruct = someStruct;
        this.someOtherStruct = someOtherStruct;
        this.yetAnotherStruct = yetAnotherStruct;
    }
    
}
