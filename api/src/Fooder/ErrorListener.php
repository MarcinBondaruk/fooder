<?php
namespace Fooder;

use Symfony\Component\EventDispatcher\EventSubscriberInterface;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\HttpKernel\Event\ExceptionEvent;
use Symfony\Component\HttpKernel\Exception\NotFoundHttpException;
use Symfony\Component\HttpKernel\KernelEvents;

final class ErrorListener implements EventSubscriberInterface
{
    public function onError(ExceptionEvent $event): void {
        $e = $event->getThrowable();

        if ($e instanceof NotFoundHttpException) {

            $event->setResponse(new Response('Page not found', 404));
        } else {
            var_dump($event->getThrowable());
            $event->setResponse(new Response('Something went wrong...', 500));
        }
    }

    public static function getSubscribedEvents(): array
    {
        return [
            KernelEvents::EXCEPTION => ['onError']
        ];
    }
}