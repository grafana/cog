<?php

namespace Grafana\Foundation\Promql;

/**
 * @implements \Grafana\Foundation\Cog\Builder<\Grafana\Foundation\Promql\FuncCallExpr>
 */
class FuncCallExprBuilder implements \Grafana\Foundation\Cog\Builder
{
    protected \Grafana\Foundation\Promql\FuncCallExpr $internal;

    public function __construct()
    {
    	$this->internal = new \Grafana\Foundation\Promql\FuncCallExpr();
    $this->internal->type = "funcCallExpr";
    }

    /**
     * Builds the object.
     * @return \Grafana\Foundation\Promql\FuncCallExpr
     */
    public function build()
    {
        return $this->internal;
    }

    public function function(string $function): static
    {
        if (!(strlen($function) >= 1)) {
            throw new \ValueError('strlen($function) must be >= 1');
        }
        $this->internal->function = $function;
    
        return $this;
    }

    /**
     * @param array<\Grafana\Foundation\Cog\Builder<\Grafana\Foundation\Promql\Expr>> $args
     */
    public function args(array $args): static
    {
            $argsResources = [];
            foreach ($args as $r1) {
                    $argsResources[] = $r1->build();
            }
        $this->internal->args = $argsResources;
    
        return $this;
    }

    /**
     * Modified by veneer 'Duplicate[args]'
     * Modified by veneer 'ArrayToAppend'
     * @param \Grafana\Foundation\Cog\Builder<\Grafana\Foundation\Promql\Expr> $arg
     */
    public function arg(\Grafana\Foundation\Cog\Builder $arg): static
    {
        $argResource = $arg->build();
        $this->internal->args[] = $argResource;
    
        return $this;
    }

}
