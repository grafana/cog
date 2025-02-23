package promql;

import java.util.List;
import java.util.LinkedList;

public class FuncCallExprBuilder implements cog.Builder<FuncCallExpr> {
    protected final FuncCallExpr internal;
    
    public FuncCallExprBuilder() {
        this.internal = new FuncCallExpr();
        this.internal.type = "funcCallExpr";
    }
    public FuncCallExprBuilder function(String function) {
        if (!(function.length() >= 1)) {
            throw new IllegalArgumentException("function.length() must be >= 1");
        }
        this.internal.function = function;
        return this;
    }
    
    public FuncCallExprBuilder args(cog.Builder<List<Expr>> args) {
        this.internal.args = args.build();
        return this;
    }
    
    public FuncCallExprBuilder arg(cog.Builder<Expr> arg) {
		if (this.internal.args == null) {
			this.internal.args = new LinkedList<>();
		}
        this.internal.args.add(arg.build());
        return this;
    }
    public FuncCallExpr build() {
        return this.internal;
    }
}
