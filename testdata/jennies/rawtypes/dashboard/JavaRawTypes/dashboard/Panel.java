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
}
