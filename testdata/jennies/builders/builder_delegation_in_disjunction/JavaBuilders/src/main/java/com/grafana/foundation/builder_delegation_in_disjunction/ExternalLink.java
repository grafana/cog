package com.grafana.foundation.builder_delegation_in_disjunction;


public class ExternalLink {
    public String url;
    
    public static class Builder {
        private ExternalLink internal;
        
        public Builder() {
            this.internal = new ExternalLink();
        }
    public Builder setUrl(String url) {
    this.internal.url = url;
        return this;
    }
    public ExternalLink build() {
            return this.internal;
        }
    }
}
