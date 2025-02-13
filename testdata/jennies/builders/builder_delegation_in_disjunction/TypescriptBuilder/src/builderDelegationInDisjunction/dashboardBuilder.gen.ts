import * as cog from '../cog';
import * as builderDelegationInDisjunction from '../builderDelegationInDisjunction';

export class DashboardBuilder implements cog.Builder<builderDelegationInDisjunction.Dashboard> {
    protected readonly internal: builderDelegationInDisjunction.Dashboard;

    constructor() {
        this.internal = builderDelegationInDisjunction.defaultDashboard();
    }

    /**
     * Builds the object.
     */
    build(): builderDelegationInDisjunction.Dashboard {
        return this.internal;
    }

    // will be expanded to cog.Builder<DashboardLink> | string
    singleLinkOrString(singleLinkOrString: cog.Builder<builderDelegationInDisjunction.DashboardLink> | string): this {
        const singleLinkOrStringResource = cog.isBuilder(singleLinkOrString) ? singleLinkOrString.build() : singleLinkOrString;
        this.internal.singleLinkOrString = singleLinkOrStringResource;
        return this;
    }

    // will be expanded to [](cog.Builder<DashboardLink> | string)
    linksOrStrings(linksOrStrings: (cog.Builder<builderDelegationInDisjunction.DashboardLink> | string)[]): this {
        const linksOrStringsResources = linksOrStrings.map(builder1 => cog.isBuilder(builder1) ? builder1.build() : builder1);
        this.internal.linksOrStrings = linksOrStringsResources;
        return this;
    }

    disjunctionOfBuilders(disjunctionOfBuilders: cog.Builder<builderDelegationInDisjunction.DashboardLink> | cog.Builder<builderDelegationInDisjunction.ExternalLink>): this {
        const disjunctionOfBuildersResource = disjunctionOfBuilders.build();
        this.internal.disjunctionOfBuilders = disjunctionOfBuildersResource;
        return this;
    }
}

