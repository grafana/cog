package constraints;


public class SomeStruct {
    public Long id;
    public Long maybeId;
    public String title;
    public RefStruct refStruct;
    public SomeStruct() {
        this.id = 0L;
        this.title = "";
    }
    public SomeStruct(Long id,Long maybeId,String title,RefStruct refStruct) {
        this.id = id;
        this.maybeId = maybeId;
        this.title = title;
        this.refStruct = refStruct;
    }
}
