<?php

require_once __DIR__.'/vendor/autoload.php';

echo 'lala'.PHP_EOL;

$dashboard = new \Grafana\Foundation\Types\Dashboard\Dashboard();
$dashboard->schemaVersion = 1;
$dashboard->templating = new \Grafana\Foundation\Types\Dashboard\DashboardDashboardTemplating();
$dashboard->annotations = new \Grafana\Foundation\Types\Dashboard\AnnotationContainer();
$dashboard->title = "Some dashboard";
$dashboard->tags = ["foo", "bar"];

var_dump(json_encode($dashboard));
