# Grafana Foundation SDK – Java

A set of tools, types and *builder libraries* for building and manipulating Grafana objects in Java.

> [!NOTE]
> This branch contains **types and builders generated for Grafana {{ .Extra.GrafanaVersion }}.**
> Other supported versions of Grafana can be found at [this repository's root](https://github.com/grafana/grafana-foundation-sdk/).

## Installing

### Gradle
```kotlin
implementation("com.grafana:grafana-foundation-sdk:{{ .Extra.GrafanaVersion|registryToSemver }}-{{ .Extra.BuildTimestamp }}")
```

### Maven
```xml
<dependency>
    <groupId>com.grafana</groupId>
    <artifactId>grafana-foundation-sdk</artifactId>
    <version>{{ .Extra.GrafanaVersion|registryToSemver }}-{{ .Extra.BuildTimestamp }}</version>
</dependency>
```

## Example usage

### Building a dashboard

```java
public class Main {

    public static void main(String[] args) {
        Dashboard dashboard = new Dashboard.Builder("Sample Dashboard").
                Uid("generated-from-java").
                Tags(List.of("generated", "from", "java")).
                Refresh("1m").Time(new DashboardDashboardTime.Builder().
                        From("now-30m").
                        To("now")
                ).
                Timezone(Constants.TimeZoneBrowser).
                WithRow(new RowPanel.Builder("Overview")).
                WithPanel(new PanelBuilder().
                        Title("Network Received").
                        Unit("bps").
                        Min(0.0).
                        WithTarget(new Dataquery.Builder().
                                Expr("rate(node_network_receive_bytes_total{job=\"integrations/raspberrypi-node\", device!=\"lo\"}[$__rate_interval]) * 8").
                                LegendFormat({{ `"{{ device }}"` }})
                        )
                ).build();

        try {
            System.out.println(dashboard.toJSON());
        } catch (JsonProcessingException e) {
            e.printStackTrace();
        }
    }
}
```

### Unmarshaling a dashboard

```java
public class Main {

    public static void main(String[] args) {
        ObjectMapper mapper = new ObjectMapper();
        try {
            InputStream json = Main.class.getResourceAsStream("/dashboard.json");
            Dashboard dashboard = mapper.readValue(json, Dashboard.class);
            System.out.println(dashboard.toJSON());
        } catch (IOException e) {
            e.printStackTrace();
        }
    }
}
```

### Defining a custom query type

While the SDK ships with support for all core datasources and their query types,
it can be extended for private/third-party plugins.

To do so, define a type and a builder for the custom query:

```java
public class CustomQuery implements Dataquery {
    
    @JsonProperty("refId")
    private String refId;
    @JsonProperty("hide")
    private Boolean hide;
    @JsonProperty("query")
    private String query;

    public static class Builder implements com.grafana.foundation.cog.Builder<Dataquery> {
        private final CustomQuery internal;

        public Builder(String query) {
            internal = new CustomQuery();
            this.internal.query = query;
        }

        public Builder refId(String refId) {
            this.internal.refId = refId;
            return this;
        }

        public Builder hide(Boolean hide) {
            this.internal.hide = hide;
            return this;
        }

        @Override
        public CustomQuery build() {
            return this.internal;
        }
    }
}
```

Register the type with cog, and use it as usual to build a dashboard:

```java
public class Main {

    public static void main(String[] args) {
        Registry.registerDataquery("custom", CustomQuery.class);

        Dashboard dashboard = new Dashboard.Builder("Custom query type").
                uid("test-custom-query-type").
                refresh("1m").
                time(new DashboardDashboardTime.Builder().from("now-30m").to("now")).
                withRow(new RowPanel.Builder("Overview")).
                withPanel(new PanelBuilder().
                        title("Sample panel").
                        withTarget(
                                new CustomQuery.Builder("query here")
                        )
                ).
                build();

        try {
            System.out.println(dashboard.toJSON());
        } catch (JsonProcessingException e) {
            throw new RuntimeException(e);
        }
    }
}
```

## Maturity

> [!WARNING]
> The code in this repository should be considered experimental. Documentation is only
available alongside the code. It comes with no support, but we are keen to receive
feedback and suggestions on how to improve it, though we cannot commit
to resolution of any particular issue.

Grafana Labs defines experimental features as follows:

> Projects and features in the Experimental stage are supported only by the Engineering
teams; on-call support is not available. Documentation is either limited or not provided
outside of code comments. No SLA is provided.
>
> Experimental projects or features are primarily intended for open source engineers who
want to participate in ensuring systems stability, and to gain consensus and approval
for open source governance projects.
>
> Projects and features in the Experimental phase are not meant to be used in production
environments, and the risks are unknown/high.

## License

[Apache 2.0 License](./LICENSE)
