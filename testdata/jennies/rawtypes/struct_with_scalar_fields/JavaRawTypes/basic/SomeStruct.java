package basic;

import java.util.Objects;

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

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof SomeStruct)) return false;
        SomeStruct o = (SomeStruct) other;
        if (!Objects.equals(this.fieldAny, o.fieldAny)) return false;
        if (!Objects.equals(this.fieldBool, o.fieldBool)) return false;
        if (!Objects.equals(this.fieldBytes, o.fieldBytes)) return false;
        if (!Objects.equals(this.fieldString, o.fieldString)) return false;
        if (!Objects.equals(this.fieldStringWithConstantValue, o.fieldStringWithConstantValue)) return false;
        if (!Objects.equals(this.fieldFloat32, o.fieldFloat32)) return false;
        if (!Objects.equals(this.fieldFloat64, o.fieldFloat64)) return false;
        if (!Objects.equals(this.fieldUint8, o.fieldUint8)) return false;
        if (!Objects.equals(this.fieldUint16, o.fieldUint16)) return false;
        if (!Objects.equals(this.fieldUint32, o.fieldUint32)) return false;
        if (!Objects.equals(this.fieldUint64, o.fieldUint64)) return false;
        if (!Objects.equals(this.fieldInt8, o.fieldInt8)) return false;
        if (!Objects.equals(this.fieldInt16, o.fieldInt16)) return false;
        if (!Objects.equals(this.fieldInt32, o.fieldInt32)) return false;
        if (!Objects.equals(this.fieldInt64, o.fieldInt64)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.fieldAny, this.fieldBool, this.fieldBytes, this.fieldString, this.fieldStringWithConstantValue, this.fieldFloat32, this.fieldFloat64, this.fieldUint8, this.fieldUint16, this.fieldUint32, this.fieldUint64, this.fieldInt8, this.fieldInt16, this.fieldInt32, this.fieldInt64);
    }
}
