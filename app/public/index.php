<?php
require_once dirname(__DIR__).'/vendor/autoload.php';

use Fooder\ErrorListener;
use Fooder\Framework;
use Symfony\Component\EventDispatcher\EventDispatcher;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\RequestStack;
use Symfony\Component\HttpKernel\Controller\ArgumentResolver;
use Symfony\Component\HttpKernel\Controller\ControllerResolver;
use Symfony\Component\HttpKernel\EventListener\RouterListener;
use Symfony\Component\Routing\Matcher\UrlMatcher;
use Symfony\Component\Routing\RequestContext;

$request = Request::createFromGlobals();
$requestStack = new RequestStack();
$routes = include __DIR__.'/../src/routes.php';

$context = new RequestContext();
$matcher = new UrlMatcher($routes, $context);

$controllerResolver = new ControllerResolver();
$argumentResolver = new ArgumentResolver();

$eventDispatcher = new EventDispatcher();
$eventDispatcher->addSubscriber(new RouterListener($matcher, $requestStack));

$eventDispatcher->addSubscriber(new ErrorListener);

$framework = new Framework($eventDispatcher, $controllerResolver, $requestStack, $argumentResolver);

$response = $framework->handle($request);
$response->send();
