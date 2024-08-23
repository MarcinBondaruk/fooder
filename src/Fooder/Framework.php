<?php
namespace Fooder;

use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\HttpKernel\Controller\ArgumentResolverInterface;
use Symfony\Component\HttpKernel\Controller\ControllerResolverInterface;
use Symfony\Component\Routing\Exception\ResourceNotFoundException;
use Symfony\Component\Routing\Matcher\UrlMatcherInterface;

final class Framework {
    public function __construct(
        private UrlMatcherInterface $urlMatcher,
        private ControllerResolverInterface $controllerResolver,
        private ArgumentResolverInterface $argumentResolver
    ) {
    }

    public function handle(Request $request): Response
    {
        try {
            $request->attributes->add($this->urlMatcher->match($request->getPathInfo()));

            $controller = $this->controllerResolver->getController($request);
            $arguments = $this->argumentResolver->getArguments($request, $controller);

            return call_user_func($controller, ...$arguments);
        } catch (ResourceNotFoundException $e) {
            return new Response('Page not found', 404);
        } catch (\Exception $e) {
            return new Response('Unexpected error occured: '.$e->getMessage(), 500);
        }
    }
}