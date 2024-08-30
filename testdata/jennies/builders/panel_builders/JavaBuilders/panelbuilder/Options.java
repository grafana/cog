package panelbuilder;

import java.util.List;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectWriter;

public class Options { 
    @JsonProperty("onlyFromThisDashboard")
    public Boolean onlyFromThisDashboard; 
    @JsonProperty("onlyInTimeRange")
    public Boolean onlyInTimeRange; 
    @JsonProperty("tags")
    public List<String> tags; 
    @JsonProperty("limit")
    public Integer limit; 
    @JsonProperty("showUser")
    public Boolean showUser; 
    @JsonProperty("showTime")
    public Boolean showTime; 
    @JsonProperty("showTags")
    public Boolean showTags; 
    @JsonProperty("navigateToPanel")
    public Boolean navigateToPanel; 
    @JsonProperty("navigateBefore")
    public String navigateBefore; 
    @JsonProperty("navigateAfter")
    public String navigateAfter;
    
    public String toJSON() throws JsonProcessingException {
        ObjectWriter ow = new ObjectMapper().writer().withDefaultPrettyPrinter();
        return ow.writeValueAsString(this);
    }

    public Options() {}

    public Options(Boolean onlyFromThisDashboard,Boolean onlyInTimeRange,List<String> tags,Integer limit,Boolean showUser,Boolean showTime,Boolean showTags,Boolean navigateToPanel,String navigateBefore,String navigateAfter) {
        this.onlyFromThisDashboard = onlyFromThisDashboard;
        this.onlyInTimeRange = onlyInTimeRange;
        this.tags = tags;
        this.limit = limit;
        this.showUser = showUser;
        this.showTime = showTime;
        this.showTags = showTags;
        this.navigateToPanel = navigateToPanel;
        this.navigateBefore = navigateBefore;
        this.navigateAfter = navigateAfter;
    }
    
}
