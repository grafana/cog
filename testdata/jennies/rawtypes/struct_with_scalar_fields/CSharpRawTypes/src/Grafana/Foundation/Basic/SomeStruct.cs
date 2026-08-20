namespace Grafana.Foundation.Basic;


// This
// is
// a
// comment
public class SomeStruct
{
    // Anything can go in there.
    // Really, anything.
    public object FieldAny;
    public bool FieldBool;
    public byte FieldBytes;
    public string FieldString;
    public string FieldStringWithConstantValue;
    public float FieldFloat32;
    public double FieldFloat64;
    public byte FieldUint8;
    public ushort FieldUint16;
    public uint FieldUint32;
    public ulong FieldUint64;
    public sbyte FieldInt8;
    public short FieldInt16;
    public int FieldInt32;
    public long FieldInt64;

    public SomeStruct()
    {
        this.FieldAny = new object();
        this.FieldBool = false;
        this.FieldBytes = (byte) 0;
        this.FieldString = "";
        this.FieldStringWithConstantValue = "";
        this.FieldFloat32 = 0f;
        this.FieldFloat64 = 0d;
        this.FieldUint8 = 0;
        this.FieldUint16 = 0;
        this.FieldUint32 = 0;
        this.FieldUint64 = 0UL;
        this.FieldInt8 = 0;
        this.FieldInt16 = 0;
        this.FieldInt32 = 0;
        this.FieldInt64 = 0L;
    }

    public SomeStruct(object fieldAny, bool fieldBool, byte fieldBytes, string fieldString, string fieldStringWithConstantValue, float fieldFloat32, double fieldFloat64, byte fieldUint8, ushort fieldUint16, uint fieldUint32, ulong fieldUint64, sbyte fieldInt8, short fieldInt16, int fieldInt32, long fieldInt64)
    {
        this.FieldAny = fieldAny;
        this.FieldBool = fieldBool;
        this.FieldBytes = fieldBytes;
        this.FieldString = fieldString;
        this.FieldStringWithConstantValue = fieldStringWithConstantValue;
        this.FieldFloat32 = fieldFloat32;
        this.FieldFloat64 = fieldFloat64;
        this.FieldUint8 = fieldUint8;
        this.FieldUint16 = fieldUint16;
        this.FieldUint32 = fieldUint32;
        this.FieldUint64 = fieldUint64;
        this.FieldInt8 = fieldInt8;
        this.FieldInt16 = fieldInt16;
        this.FieldInt32 = fieldInt32;
        this.FieldInt64 = fieldInt64;
    }
}
