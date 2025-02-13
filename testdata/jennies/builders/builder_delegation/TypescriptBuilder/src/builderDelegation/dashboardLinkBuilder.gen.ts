import * as cog from '../cog';
import * as builderDelegation from '../builderDelegation';

export class DashboardLinkBuilder implements cog.Builder<builderDelegation.DashboardLink> {
    protected readonly internal: builderDelegation.DashboardLink;

    constructor() {
        this.internal = builderDelegation.defaultDashboardLink();
    }

    /**
     * Builds the object.
     */
    build(): builderDelegation.DashboardLink {
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

