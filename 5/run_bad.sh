for i in {1..20};
    do curl -s -o /dev/null -H "Authorization: Bearer wrong" http://localhost:8082/tasks
done
