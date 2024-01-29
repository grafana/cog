import * as cog from '../cog';
import * as builder_delegation from '../builder_delegation';

export class DashboardLinkBuilder implements cog.Builder<builder_delegation.DashboardLink> {
    private readonly internal: builder_delegation.DashboardLink;

    constructor() {
        this.internal = builder_delegation.defaultDashboardLink();
    }

    build(): builder_delegation.DashboardLink {
        return this.internal;
    }

    title(title: string): this {
        this.internal.title = title;
        return this;
    }

    url(url: string): this {
        this.internal.url = url;
        return this;
    }
}
