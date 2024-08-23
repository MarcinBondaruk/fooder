<?php

namespace Calendar\Controller;

use Calendar\Model\LeapYear;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\Response;

final class LeapYearController {
    public function getAnswer(Request $request, int $year): Response {
        $leapYear = new LeapYear();
        $content = $year.' is not a leap year.';

        if ($leapYear->isLeapYear($year)) {
            $content = $year.' is a leap year.';
        }

        $response = new Response($content);
        $response->setStatusCode(200);
        $response->prepare($request);
        return $response;
    }
}