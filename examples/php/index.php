<?php

require_once __DIR__.'/vendor/autoload.php';

echo 'lala'.PHP_EOL;

$dashboard = new \Grafana\Foundation\Dashboard\Dashboard(uid: "foo");
$dashboard->title = "Some dashboard";
$dashboard->tags = ["foo", "bar"];

var_dump(json_encode($dashboard));


$dashBuilder = new \Grafana\Foundation\Dashboard\DashboardBuilder(title: "Awesome dashboard");
$dashBuilder->refresh("5m");
$dash = $dashBuilder->build();

var_dump(json_encode($dash));
