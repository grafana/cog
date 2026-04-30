package dashboard;

import java.util.Objects;
import java.util.List;

public class Dashboard {
    public String title;
    public List<Panel> panels;
    public Dashboard() {
        this.title = "";
    }
    public Dashboard(String title,List<Panel> panels) {
        this.title = title;
        this.panels = panels;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof Dashboard)) return false;
        Dashboard o = (Dashboard) other;
        if (!Objects.equals(this.title, o.title)) return false;
        if (!Objects.equals(this.panels, o.panels)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.title, this.panels);
    }
}
