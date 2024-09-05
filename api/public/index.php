<?php
require_once dirname(__DIR__).'/vendor/autoload.php';

use Symfony\Component\HttpFoundation\Request;

$routes = include __DIR__.'/../src/routes.php';
$container = include __DIR__.'/../src/container.php';

$request = Request::createFromGlobals();

$response = $container->get('framework')->handle($request);
$response->send();
