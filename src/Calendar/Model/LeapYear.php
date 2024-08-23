<?php
namespace Calendar\Model;

final class LeapYear
{
    public function isLeapYear(int $year): bool
    {
        if (null === $year) {
            $year = date('Y');
        }

        return 0 == $year % 400 || (0 == $year % 4 && 0 != $year % 100);
    }
}