package disjunctions;


public class SeveralRefs {
    protected SomeStruct someStruct;
    protected SomeOtherStruct someOtherStruct;
    protected YetAnotherStruct yetAnotherStruct;
    protected SeveralRefs() {}
    public static SeveralRefs createSomeStruct(SomeStruct someStruct) {
        SeveralRefs severalRefs = new SeveralRefs();
        severalRefs.someStruct = someStruct;
        return severalRefs;
    }
    public static SeveralRefs createSomeOtherStruct(SomeOtherStruct someOtherStruct) {
        SeveralRefs severalRefs = new SeveralRefs();
        severalRefs.someOtherStruct = someOtherStruct;
        return severalRefs;
    }
    public static SeveralRefs createYetAnotherStruct(YetAnotherStruct yetAnotherStruct) {
        SeveralRefs severalRefs = new SeveralRefs();
        severalRefs.yetAnotherStruct = yetAnotherStruct;
        return severalRefs;
    }
}
