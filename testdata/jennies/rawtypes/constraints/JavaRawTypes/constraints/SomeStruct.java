package constraints;


public class SomeStruct {
    public Long id;
    public Long maybeId;
    public String title;
    public refStruct refStruct;
    public SomeStruct() {
    }
    
    public SomeStruct(Long id,Long maybeId,String title,refStruct refStruct) {
        this.id = id;
        this.maybeId = maybeId;
        this.title = title;
        this.refStruct = refStruct;
    }
}
