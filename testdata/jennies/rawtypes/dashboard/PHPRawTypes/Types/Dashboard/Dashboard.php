<?php

namespace Types\Dashboard;

class Dashboard
{
    public string $title;

    /**
     * @var array<\Types\Dashboard\Panel>
     */
    public ?array $panels;
}
