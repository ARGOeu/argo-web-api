#!/bin/bash



for i in {0..100}
do
for j in {0..4}
do
time curl "http://webserver01.grid.auth.gr:8080/api/v1/service_availability_in_profile?profile_name=ROC_CRITICAL&group_type=Site&start_time=2013-08-01T10:00:00Z&end_time=2013-08-30T10:00:00Z&type=HOURLY&output=XML" >> /dev/null &
done
wait
done


