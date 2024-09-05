<?php

use CookBook\Controller as CookbookController;
use CookBook\Repository as CookbookRepository;
use CookBook\Service as CookbookService;
use Fooder\ErrorListener;
use Fooder\Framework;
use Groceries\Controller as GroceriesController;
use Symfony\Component\DependencyInjection\ContainerBuilder;
use Symfony\Component\DependencyInjection\Reference;
use Symfony\Component\EventDispatcher\EventDispatcher;
use Symfony\Component\HttpFoundation\RequestStack;
use Symfony\Component\HttpKernel\Controller\ArgumentResolver;
use Symfony\Component\HttpKernel\Controller\ContainerControllerResolver;
use Symfony\Component\HttpKernel\EventListener\RouterListener;
use Symfony\Component\Routing\Matcher\UrlMatcher;
use Symfony\Component\Routing\RequestContext;

/**
 * This part contains framework setup
 */
$container = new ContainerBuilder();
$container->register('context', RequestContext::class);
$container->register('matcher', UrlMatcher::class)
    ->setArguments([$routes, new Reference('context')]);

$container->register('request_stack', RequestStack::class);
$container->register('controller_resolver', ContainerControllerResolver::class)
    ->setArguments([$container]);
$container->register('argument_resolver', ArgumentResolver::class);

$container->register('listener.router_listener', RouterListener::class)
    ->setArguments([new Reference('matcher'), new Reference('request_stack')]);

$container->register('listener.error_listener', ErrorListener::class);

$container->register('event_dispatcher', EventDispatcher::class)
    ->addMethodCall('addSubscriber', [new Reference('listener.router_listener')])
    ->addMethodCall('addSubscriber', [new Reference('listener.error_listener')]);

$container->register('framework', Framework::class)
    ->setArguments([
        new Reference('event_dispatcher'),
        new Reference('controller_resolver'),
        new Reference('request_stack'),
        new Reference('argument_resolver'),
    ]);

/**
 * setup by Service
 */
$container->register('cookbook.repo', CookbookRepository::class);
$container->register('cookbook.svc', CookbookService::class)
    ->setArguments([new Reference('cookbook.repo')]);
$container->register('cookbook.ctrl', CookbookController::class)
    ->setArguments([new Reference('cookbook.svc')]);

$container->register('groceries.ctrl', GroceriesController::class);

return $container;