<?php
require_once dirname(__DIR__).'/vendor/autoload.php';

use Fooder\Framework;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpKernel\Controller\ArgumentResolver;
use Symfony\Component\HttpKernel\Controller\ControllerResolver;
use Symfony\Component\Routing\Matcher\UrlMatcher;
use Symfony\Component\Routing\RequestContext;

function renderTemplate(Request $request): string
{
    extract($request->attributes->all());
    ob_start();
    include(sprintf(__DIR__."/../templates/%s.php", $_route));

    return ob_get_clean();
}

$request = Request::createFromGlobals();
$routes = include __DIR__.'/../src/routes.php';

$context = new RequestContext();
$context->fromRequest($request);
$matcher = new UrlMatcher($routes, $context);

$controllerResolver = new ControllerResolver();
$argumentResolver = new ArgumentResolver();

$framework = new Framework($matcher, $controllerResolver, $argumentResolver);

$response = $framework->handle($request);

$response->send();
