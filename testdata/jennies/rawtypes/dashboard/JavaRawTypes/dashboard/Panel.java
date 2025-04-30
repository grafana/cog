package dashboard;

import java.util.List;
import cog.variants.Dataquery;

public class Panel {
    public String title;
    public String type;
    public DataSourceRef datasource;
    public Object options;
    public List<Dataquery> targets;
    public FieldConfigSource fieldConfig;
    public Panel() {
        this.title = "";
        this.type = "";
    }
    public Panel(String title,String type,DataSourceRef datasource,Object options,List<Dataquery> targets,FieldConfigSource fieldConfig) {
        this.title = title;
        this.type = type;
        this.datasource = datasource;
        this.options = options;
        this.targets = targets;
        this.fieldConfig = fieldConfig;
    }
}
