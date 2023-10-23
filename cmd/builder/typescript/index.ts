import {DashboardCursorSync, DashboardLinkType} from "../../../generated/types/dashboard/types_gen";
import {GraphDrawStyle, TooltipDisplayMode} from "../../../generated/types/common/types_gen";
import {DashboardBuilder} from "../../../generated/dashboard/dashboard/builder_gen";
import {RowPanelBuilder} from "../../../generated/dashboard/rowpanel/builder_gen";
import {PanelBuilder} from "../../../generated/timeseries/panel/builder_gen";
import {VizTooltipOptionsBuilder} from "../../../generated/common/viztooltipoptions/builder_gen";
import {TimePickerBuilder} from "../../../generated/dashboard/timepicker/builder_gen";

const timeseriesPanel = new PanelBuilder()
    .title("Some timeseries panel")
    .transparent(true)
    .description("Let there be data")
    .decimals(2)
    .min(0)
    .max(200)
    .lineWidth(5)
    .drawStyle(GraphDrawStyle.GraphDrawStyleBars)
    .tooltip(new VizTooltipOptionsBuilder().mode(TooltipDisplayMode.TooltipDisplayModeSingle));

const overviewRow = new RowPanelBuilder("Overview");

const builder = new DashboardBuilder("Some title")
    .uid("test-dashboard-codegen")
    .description("Some description")

    .refresh("1m")
    .time({from: "now-3h", to: "now"})
    .timezone("utc")

    .timepicker(
        new TimePickerBuilder()
            .refresh_intervals(["30s", "1m", "5m"]),
    )

    .tooltip(DashboardCursorSync.DashboardCursorSyncCrosshair)
    .tags(["generated", "from", "cue"])
    .links([
        {
            // TODO: this is painful.
            title: "Some link",
            url: "http://google.com",
            type: DashboardLinkType.DashboardLinkTypeLink,
            tags: [],
            icon: "cloud",
            tooltip: "",
            asDropdown: false,
            targetBlank: false,
            includeVars: false,
            keepTime: false,
        },
    ])

    .panel(overviewRow.build())
    .panel(timeseriesPanel.build())
;

console.log(JSON.stringify(builder.build(), null, 2));

