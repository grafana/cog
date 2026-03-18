package time_hint;


public class ObjWithTimeField {
    public String registeredAt;
    public String duration;
    public ObjWithTimeField() {
        this.registeredAt = "";
        this.duration = "";
    }
    public ObjWithTimeField(String registeredAt,String duration) {
        this.registeredAt = registeredAt;
        this.duration = duration;
    }
}
