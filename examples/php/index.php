<?php

require_once __DIR__.'/vendor/autoload.php';

echo 'lala'.PHP_EOL;

$dashboard = new \Grafana\Foundation\Types\Dashboard\Dashboard(uid: "foo");
$dashboard->title = "Some dashboard";
$dashboard->tags = ["foo", "bar"];

var_dump(json_encode($dashboard));
