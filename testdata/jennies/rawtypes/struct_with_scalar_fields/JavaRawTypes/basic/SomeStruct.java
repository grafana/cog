package basic;


// This
// is
// a
// comment
public class SomeStruct {
    // Anything can go in there.
    // Really, anything.
    public Object fieldAny;
    public Boolean fieldBool;
    public Byte fieldBytes;
    public String fieldString;
    public String fieldStringWithConstantValue;
    public Float fieldFloat32;
    public Double fieldFloat64;
    public Integer fieldUint8;
    public Short fieldUint16;
    public Integer fieldUint32;
    public Long fieldUint64;
    public Integer fieldInt8;
    public Short fieldInt16;
    public Integer fieldInt32;
    public Long fieldInt64;
    public SomeStruct() {
        this.fieldAny = new Object();
        this.fieldBool = false;
        this.fieldBytes = (byte) 0;
        this.fieldString = "";
        this.fieldStringWithConstantValue = "";
        this.fieldFloat32 = 0.0f;
        this.fieldFloat64 = 0.0;
        this.fieldUint8 = 0;
        this.fieldUint16 = 0;
        this.fieldUint32 = 0;
        this.fieldUint64 = 0L;
        this.fieldInt8 = 0;
        this.fieldInt16 = 0;
        this.fieldInt32 = 0;
        this.fieldInt64 = 0L;
    }
    public SomeStruct(Object fieldAny,Boolean fieldBool,Byte fieldBytes,String fieldString,String fieldStringWithConstantValue,Float fieldFloat32,Double fieldFloat64,Integer fieldUint8,Short fieldUint16,Integer fieldUint32,Long fieldUint64,Integer fieldInt8,Short fieldInt16,Integer fieldInt32,Long fieldInt64) {
        this.fieldAny = fieldAny;
        this.fieldBool = fieldBool;
        this.fieldBytes = fieldBytes;
        this.fieldString = fieldString;
        this.fieldStringWithConstantValue = fieldStringWithConstantValue;
        this.fieldFloat32 = fieldFloat32;
        this.fieldFloat64 = fieldFloat64;
        this.fieldUint8 = fieldUint8;
        this.fieldUint16 = fieldUint16;
        this.fieldUint32 = fieldUint32;
        this.fieldUint64 = fieldUint64;
        this.fieldInt8 = fieldInt8;
        this.fieldInt16 = fieldInt16;
        this.fieldInt32 = fieldInt32;
        this.fieldInt64 = fieldInt64;
    }
}
