<?php

require_once __DIR__.'/vendor/autoload.php';

$dashBuilder = new \Grafana\Foundation\Dashboard\DashboardBuilder(title: "Awesome dashboard");
$dashBuilder->refresh("5m");
$dash = $dashBuilder->build();

var_dump(json_encode($dash));
