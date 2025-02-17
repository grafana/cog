import * as cog from '../cog';
import * as builderDelegationInDisjunction from '../builderDelegationInDisjunction';

export class DashboardLinkBuilder implements cog.Builder<builderDelegationInDisjunction.DashboardLink> {
    protected readonly internal: builderDelegationInDisjunction.DashboardLink;

    constructor() {
        this.internal = builderDelegationInDisjunction.defaultDashboardLink();
    }

    /**
     * Builds the object.
     */
    build(): builderDelegationInDisjunction.DashboardLink {
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

