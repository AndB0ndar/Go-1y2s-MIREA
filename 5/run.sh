for i in {1..50}; do
    curl -s -o /dev/null -H "Authorization: Bearer demo-token" http://localhost:8082/tasks
done
