package dashboard;

import java.util.Objects;
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

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof Panel)) return false;
        Panel o = (Panel) other;
        if (!Objects.equals(this.title, o.title)) return false;
        if (!Objects.equals(this.type, o.type)) return false;
        if (!Objects.equals(this.datasource, o.datasource)) return false;
        if (!Objects.equals(this.options, o.options)) return false;
        if (!Objects.equals(this.targets, o.targets)) return false;
        if (!Objects.equals(this.fieldConfig, o.fieldConfig)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.title, this.type, this.datasource, this.options, this.targets, this.fieldConfig);
    }
}
