import {DashboardCursorSync, DashboardLinkType} from "../../../generated/types/dashboard/types_gen";
import {GraphDrawStyle, TooltipDisplayMode} from "../../../generated/types/common/types_gen";
import {dashboardBuilder} from "../../../generated/dashboard/dashboard/builder_gen";
import {RowPanelBuilder} from "../../../generated/dashboard/rowpanel/builder_gen";
import {PanelBuilder} from "../../../generated/timeseries/panel/builder_gen";
import {VizTooltipOptionsBuilder} from "../../../generated/common/viztooltipoptions/builder_gen";

const timeseriesPanel = new PanelBuilder()
    .title("Some timeseries panel")
    .transparent(true)
    .description("Let there be data")
    .decimals(2)
    .min(0)
    .max(200)
    .lineWidth(5)
    .drawStyle(GraphDrawStyle.Bars)
    .tooltip(new VizTooltipOptionsBuilder().mode(TooltipDisplayMode.Single));

const overviewRow = new RowPanelBuilder("Overview");

const builder = new dashboardBuilder("Some title")
    .uid("test-dashboard-codegen")
    .description("Some description")
    .time({from: "now-3h", to: "now"})
    .refresh("1m")
    .timezone("utc")
    .tooltip(DashboardCursorSync.Crosshair)
    .tags(["generated", "from", "cue"])
    .rows([
        overviewRow.build(),
        timeseriesPanel.build(),
    ])
    .links([
        {
            // TODO: this is painful.
            title: "Some link",
            url: "http://google.com",
            type: DashboardLinkType.Link,
            tags: [],
            icon: "cloud",
            tooltip: "",
            asDropdown: false,
            targetBlank: false,
            includeVars: false,
            keepTime: false,
        },
    ]);

console.log(builder.build());

