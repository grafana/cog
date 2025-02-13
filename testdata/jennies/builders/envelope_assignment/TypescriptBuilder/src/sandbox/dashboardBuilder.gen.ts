import * as cog from '../cog';
import * as sandbox from '../sandbox';

export class DashboardBuilder implements cog.Builder<sandbox.Dashboard> {
    protected readonly internal: sandbox.Dashboard;

    constructor() {
        this.internal = sandbox.defaultDashboard();
    }

    /**
     * Builds the object.
     */
    build(): sandbox.Dashboard {
        return this.internal;
    }

    withVariable(name: string,value: string): this {
        if (!this.internal.variables) {
            this.internal.variables = [];
        }
        this.internal.variables.push({
        name: name,
        value: value,
    });
        return this;
    }
}

