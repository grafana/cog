package dashboard;

import java.util.Objects;

public class DataSourceRef {
    public String type;
    public String uid;
    public DataSourceRef() {
    }
    public DataSourceRef(String type,String uid) {
        this.type = type;
        this.uid = uid;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof DataSourceRef)) return false;
        DataSourceRef o = (DataSourceRef) other;
        if (!Objects.equals(this.type, o.type)) return false;
        if (!Objects.equals(this.uid, o.uid)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.type, this.uid);
    }
}
