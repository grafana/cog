package dashboard;

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
}
