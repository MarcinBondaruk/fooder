<?php

use Recipe\Controller as RecipeController;
use Recipe\Repository\RecipeRepository as RecipeRepository;
use Recipe\CookbookService;
use Fooder\Listener\ErrorListener;
use Fooder\Framework;
use Grocery\Controller as GroceriesController;
use Recipe\Repository\RecipeQueryRepository;
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

$container->register('pdo', PDO::class)
    ->setArguments([
        // TODO: get arguments from envs
    ]);

/**
 * setup by Service
 */
$container->register('recipe.repo', RecipeRepository::class)
    ->setArguments([
        new Reference('pdo'),
    ]);
$container->register('recipe.query_repo', RecipeQueryRepository::class)
    ->setArguments([
        new Reference('pdo'),
    ]);
$container->register('recipe.svc', CookbookService::class)
    ->setArguments([
        new Reference('recipe.repo'),
        new Reference('recipe.query_repo'),
    ]);
$container->register('recipe.ctrl', RecipeController::class)
    ->setArguments([new Reference('recipe.svc')]);

$container->register('groceries.ctrl', GroceriesController::class);

return $container;