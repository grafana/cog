package disjunction_anonymous;


public class StructAOrStringOrInt64 {
    protected StructA structA;
    protected String string;
    protected Long int64;
    protected StructAOrStringOrInt64() {}
    public static StructAOrStringOrInt64 createStructA(StructA structA) {
        StructAOrStringOrInt64 structAOrStringOrInt64 = new StructAOrStringOrInt64();
        structAOrStringOrInt64.structA = structA;
        return structAOrStringOrInt64;
    }
    public static StructAOrStringOrInt64 createString(String string) {
        StructAOrStringOrInt64 structAOrStringOrInt64 = new StructAOrStringOrInt64();
        structAOrStringOrInt64.string = string;
        return structAOrStringOrInt64;
    }
    public static StructAOrStringOrInt64 createInt64(Long int64) {
        StructAOrStringOrInt64 structAOrStringOrInt64 = new StructAOrStringOrInt64();
        structAOrStringOrInt64.int64 = int64;
        return structAOrStringOrInt64;
    }
}
