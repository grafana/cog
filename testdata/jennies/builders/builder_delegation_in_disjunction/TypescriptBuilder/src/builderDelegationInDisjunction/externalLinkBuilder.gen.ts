import * as cog from '../cog';
import * as builderDelegationInDisjunction from '../builderDelegationInDisjunction';

export class ExternalLinkBuilder implements cog.Builder<builderDelegationInDisjunction.ExternalLink> {
    protected readonly internal: builderDelegationInDisjunction.ExternalLink;

    constructor() {
        this.internal = builderDelegationInDisjunction.defaultExternalLink();
    }

    /**
     * Builds the object.
     */
    build(): builderDelegationInDisjunction.ExternalLink {
        return this.internal;
    }

    url(url: string): this {
        this.internal.url = url;
        return this;
    }
}

