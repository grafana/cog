package disjunction_anonymous;


public class MyStruct {
    public StringOrBoolOrFloat64OrInt64 scalars;
    public MyStructSameKind sameKind;
    public Object refs;
    public StructAOrStringOrInt64 mixed;
    public MyStruct() {
        this.scalars = new disjunction_anonymous.StringOrBoolOrFloat64OrInt64();
        this.sameKind = MyStructSameKind.A;
        this.refs = new Object();
        this.mixed = new disjunction_anonymous.StructAOrStringOrInt64();
    }
    public MyStruct(StringOrBoolOrFloat64OrInt64 scalars,MyStructSameKind sameKind,Object refs,StructAOrStringOrInt64 mixed) {
        this.scalars = scalars;
        this.sameKind = sameKind;
        this.refs = refs;
        this.mixed = mixed;
    }
}
