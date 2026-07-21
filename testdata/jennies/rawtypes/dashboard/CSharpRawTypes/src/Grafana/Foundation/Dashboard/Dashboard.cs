namespace Grafana.Foundation.Dashboard;

using System.Collections.Generic;

public class Dashboard
{
    public string Title;
    public List<Panel> Panels;

    public Dashboard()
    {
        this.Title = "";
    }

    public Dashboard(string title, List<Panel> panels)
    {
        this.Title = title;
        this.Panels = panels;
    }
}
