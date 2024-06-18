<?php

namespace Grafana\Foundation\Runtime;

/**
 * @template T
 */
interface Builder
{
    /**
     * @return T
     */
    public function build();
}
