import * as cog from '../cog';
import * as builder_delegation from '../builder_delegation';

export class DashboardBuilder implements cog.Builder<builder_delegation.Dashboard> {
    private readonly internal: builder_delegation.Dashboard;

    constructor() {
        this.internal = builder_delegation.defaultDashboard();
    }

    build(): builder_delegation.Dashboard {
        return this.internal;
    }

    id(id: number): this {
        this.internal.id = id;
        return this;
    }

    title(title: string): this {
        this.internal.title = title;
        return this;
    }

    // will be expanded to []cog.Builder<DashboardLink>
    links(links: cog.Builder<builder_delegation.DashboardLink>[]): this {
        const linksResources = links.map(builder1 => builder1.build());
        this.internal.links = linksResources;
        return this;
    }

    // will be expanded to [][]cog.Builder<DashboardLink>
    linksOfLinks(linksOfLinks: cog.Builder<builder_delegation.DashboardLink>[][]): this {
        const linksOfLinksResources = linksOfLinks.map(builder1 => builder1.map(builder2 => builder2.build()));
        this.internal.linksOfLinks = linksOfLinksResources;
        return this;
    }

    // will be expanded to cog.Builder<DashboardLink>
    singleLink(singleLink: cog.Builder<builder_delegation.DashboardLink>): this {
        const singleLinkResource = singleLink.build();
        this.internal.singleLink = singleLinkResource;
        return this;
    }
}
